package serving

import "github.com/spf13/cobra"

var (
	serveLong = `serve a job.

Available Commands:
  tensorflow,tf  Submit a TensorFlow Serving Job.
  tensorrt,trt   Submit a TensorRT Serving Job
  custom         Submit a Custom Serving Job  
  kfserving,kfs  Submit a kubeflow Serving Job`
)

func NewServeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "serve",
		Short: "Serve a job.",
		Long:  serveLong,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	command.AddCommand(NewSubmitTFServingJobCommand())
	command.AddCommand(NewSubmitTRTServingJobCommand())
	command.AddCommand(NewSubmitCustomServingJobCommand())
	command.AddCommand(NewSubmitKFServingJobCommand())
	command.AddCommand(NewListCommand())
	command.AddCommand(NewDeleteCommand())
	command.AddCommand(NewGetCommand())
	command.AddCommand(NewLogsCommand())
	command.AddCommand(NewTrafficRouterSplitCommand())
	return command
}
