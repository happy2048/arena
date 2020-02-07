// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"os"

	"fmt"

	"github.com/kubeflow/arena/pkg/util"
	"github.com/kubeflow/arena/pkg/util/helm"
	"github.com/kubeflow/arena/pkg/workflow"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

var (
	servingVersion string
	servingType    string
)

// NewDeleteCommand
func NewServingDeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "delete a serving job",
		Short: "delete a serving job and its associated pods",
		Run: func(cmd *cobra.Command, args []string) {
			util.SetLogLevel(logLevel)
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			setupKubeconfig()
			client, err := initKubeClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = updateNamespace(cmd)
			if err != nil {
				log.Debugf("Failed due to %v", err)
				fmt.Println(err)
				os.Exit(1)
			}
			exitCode := 0
			for _, jobName := range args {
				err = deleteServingJob(client, jobName)
				if err != nil {
					log.Errorf("Failed due to %v", err)
					exitCode = 2
				}
			}
			os.Exit(exitCode)
		},
	}
	command.Flags().StringVar(&servingVersion, "version", "", "The serving version to delete.")
	command.Flags().StringVar(&servingType, "type", "", "The serving type to delete, the possible option is custom, tf-serving, trt-servign. (optional)")

	// command.MarkFlagRequired("version")

	return command
}

func deleteServingJob(client *kubernetes.Clientset, servingJob string) error {
	var servingTypes []string
	err := helm.DeleteRelease(servingJob)
	if err == nil {
		log.Infof("Delete the job %s successfully.", servingJob)
		return nil
	}

	log.Debugf("%s wasn't deleted by helm due to %v", servingJob, err)

	// 1. If serving Version is not set, detect it automatically
	if servingVersion == "" {
		servings, err := ListServingsByName(client, servingJob)
		if err != nil {
			return err
		}

		if len(servings) == 0 {
			return fmt.Errorf("There is no serving job found with the name %s, please check it with `arena serve list | grep %s`",
				servingJob,
				servingJob)
		} else if len(servings) > 1 {
			return fmt.Errorf("There are more than one serving job found with the name %s, please delete it with `arena serve delete %s --version {version}`",
				servingJob,
				servingJob)
		}

		servingVersion = servings[0].Version
	}

	servingJobWithVersion := servingJob + "-" + servingVersion

	// 2. If serving Type is not set, detect it automatically
	if servingType == "" {
		servingTypes = getServingTypes(servingJobWithVersion, namespace)
		if len(servingTypes) == 0 {
			return fmt.Errorf("There is no serving job found with the name %s with version %s, please check it with `arena serve list | grep %s`",
				servingJob,
				servingType,
				servingJob)
		} else if len(servingTypes) > 1 {
			return fmt.Errorf("There are more than one serving job found with the name %s, please check it with `arena serve list | grep %s`",
				servingJob,
				servingJob)
		}

		servingType = servingTypes[0]
	}

	// 4. Delete serving
	err = workflow.DeleteJob(servingJobWithVersion, namespace, servingType)
	if err != nil {
		return err
	}
	log.Infof("The Serving job %s with version %s has been deleted successfully", servingJob, servingVersion)
	return nil
}

func deleteServingJobByHelm(servingJob string) error {
	return helm.DeleteRelease(servingJob)
}
