package topnode

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/kubeflow/arena/pkg/apis/types"
	v1 "k8s.io/api/core/v1"
)

var NormalNodeDescription = `
  1.This node has none gpu devices
`

var normalNodeTemplate = `
Name:    %v
Status:  %v
Role:    %v
Type:    %v
Address: %v
Description:
%v
-----------------------------------------------------------------------------------------
`

type normalNode struct {
	node *v1.Node
	pods []*v1.Pod
	baseNode
}

func NewNormalNode(node *v1.Node, pods []*v1.Pod, index int, args buildNodeArgs) (Node, error) {
	return &normalNode{
		node: node,
		pods: pods,
		baseNode: baseNode{
			index:    index,
			node:     node,
			pods:     pods,
			nodeType: types.NormalNode,
		},
	}, nil
}

func (n *normalNode) AllDevicesAreHealthy() bool {
	return true
}

func (n *normalNode) convert2NodeInfo() types.NormalNodeInfo {
	role := strings.Join(n.Role(), ",")
	return types.NormalNodeInfo{
		CommonNodeInfo: types.CommonNodeInfo{
			Name:        n.Name(),
			IP:          n.IP(),
			Status:      n.Status(),
			Role:        role,
			Type:        types.NormalNode,
			Description: NormalNodeDescription,
		},
	}
}

func (n *normalNode) Convert2NodeInfo() interface{} {
	return n.convert2NodeInfo()
}

func (n *normalNode) WideFormat() string {
	role := strings.Join(n.Role(), ",")
	if role == "" {
		role = "<none>"
	}
	nodeInfo := n.convert2NodeInfo()
	return fmt.Sprintf(strings.TrimRight(normalNodeTemplate, "\n"),
		n.Name(),
		n.Status(),
		role,
		n.Type(),
		n.IP(),
		strings.Trim(nodeInfo.Description, "\n"),
	)
}

/*
display like:

===================================== NormalNode ========================================

Name:    cn-shanghai.192.168.7.186
Status:  Ready
Role:    <none>
Type:    Normal
Address: 192.168.7.186
-----------------------------------------------------------------------------------------

===================================== End ===============================================
*/
func displayNormalNodeDetails(w *tabwriter.Writer, nodes []Node) {
	if len(nodes) == 0 {
		return
	}
	PrintLine(w, "===================================== GPU MODE: None =======================================")
	for _, node := range nodes {
		PrintLine(w, node.WideFormat())
	}
	PrintLine(w, "")
}

func displayNormalNodeSummary(w *tabwriter.Writer, nodes []Node, isUnhealthy, showMode bool) (int, int, int) {
	for _, node := range nodes {
		nodeInfo := node.Convert2NodeInfo().(types.NormalNodeInfo)
		items := []string{}
		items = append(items, node.Name())
		items = append(items, node.IP())
		role := nodeInfo.Role
		if role == "" {
			role = "<none>"
		}
		items = append(items, role)
		items = append(items, node.Status())
		items = append(items, "0")
		items = append(items, "0")
		if showMode {
			for _, typeInfo := range types.NodeTypeSlice {
				if typeInfo.Name == types.NormalNode {
					items = append(items, typeInfo.Alias)
				}
			}
		}
		if isUnhealthy {
			items = append(items, fmt.Sprintf("0"))
		}
		PrintLine(w, items...)
	}
	return 0, 0, 0
}

func NewNormalNodeProcesser() NodeProcesser {
	return &nodeProcesser{
		nodeType:            types.NormalNode,
		key:                 "normalNodes",
		builder:             NewNormalNode,
		canBuildNode:        func(node *v1.Node) bool { return true },
		displayNodesDetails: displayNormalNodeDetails,
		displayNodesSummary: displayNormalNodeSummary,
	}
}
