package arenaclient

import (
	"fmt"
	"time"

	"github.com/kubeflow/arena/pkg/apis/config"
	apistraining "github.com/kubeflow/arena/pkg/apis/training"
	"github.com/kubeflow/arena/pkg/apis/types"
	"github.com/kubeflow/arena/pkg/apis/utils"
	"github.com/kubeflow/arena/pkg/training"
)

var (
	errJobNotFoundMessage = "Not found training job %s in namespace %s,please use 'arena submit' to create it."
)

// TrainingJobClient provides some operators for managing training jobs.
type TrainingJobClient struct {
	// namespace store the namespace
	namespace            string
	arenaSystemNamespace string
	configer             *config.ArenaConfiger
}

// NewTrainingJobClient creates a TrainingJobClient
func NewTrainingJobClient(namespace, arenaSystemNamespace string, configer *config.ArenaConfiger) *TrainingJobClient {
	return &TrainingJobClient{
		namespace:            namespace,
		arenaSystemNamespace: arenaSystemNamespace,
		configer:             configer,
	}
}

// Namespace sets the namespace
func (t *TrainingJobClient) Namespace(namespace string) *TrainingJobClient {
	copyTrainingJobClient := &TrainingJobClient{
		namespace:            namespace,
		arenaSystemNamespace: t.arenaSystemNamespace,
		configer:             t.configer,
	}
	return copyTrainingJobClient
}

// Submit submits a training job
func (t *TrainingJobClient) Submit(job *apistraining.Job) error {
	switch job.Type() {
	case types.TFTrainingJob:
		args := job.Args().(*types.SubmitTFJobArgs)
		return training.SubmitTFJob(t.namespace, args)
	case types.PytorchTrainingJob:
		args := job.Args().(*types.SubmitPyTorchJobArgs)
		return training.SubmitPytorchJob(t.namespace, args)
	case types.MPITrainingJob:
		args := job.Args().(*types.SubmitMPIJobArgs)
		return training.SubmitMPIJob(t.namespace, args)
	case types.HorovodTrainingJob:
		args := job.Args().(*types.SubmitHorovodJobArgs)
		return training.SubmitHorovodJob(t.namespace, args)
	case types.VolcanoTrainingJob:
		args := job.Args().(*types.SubmitVolcanoJobArgs)
		return training.SubmitVolcanoJob(t.namespace, args)
	case types.ETTrainingJob:
		args := job.Args().(*types.SubmitETJobArgs)
		return training.SubmitETJob(t.namespace, args)
	case types.SparkTrainingJob:
		args := job.Args().(*types.SubmitSparkJobArgs)
		return training.SubmitSparkJob(t.namespace, args)
	}
	return nil
}

// ScaleIn scales in job
func (t *TrainingJobClient) ScaleIn(job *apistraining.Job) error {
	switch job.Type() {
	case types.ETTrainingJob:
		args := job.Args().(*types.ScaleInETJobArgs)
		return training.SubmitScaleInETJob(t.namespace, args)
	}
	return nil
}

// ScaleOut scales out job
func (t *TrainingJobClient) ScaleOut(job *apistraining.Job) error {
	switch job.Type() {
	case types.ETTrainingJob:
		args := job.Args().(*types.ScaleOutETJobArgs)
		return training.SubmitScaleOutETJob(t.namespace, args)
	}
	return nil
}

// Get returns a training job information
func (t *TrainingJobClient) Get(jobName string, jobType types.TrainingJobType) (*types.TrainingJobInfo, error) {
	job, err := training.SearchTrainingJob(jobName, t.namespace, jobType)
	if err != nil {
		return nil, err
	}
	jobInfo := training.BuildJobInfo(job, true)
	return jobInfo, nil
}

// GetAndPrint print training job information
func (t *TrainingJobClient) GetAndPrint(jobName string, jobType types.TrainingJobType, format string, showEvent bool, showGPU bool) error {
	if utils.TransferPrintFormat(format) == types.UnknownFormat {
		return fmt.Errorf("Unknown output format,only support:[wide|json|yaml]")
	}
	job, err := training.SearchTrainingJob(jobName, t.namespace, jobType)
	if err != nil {
		if err == types.ErrTrainingJobNotFound {
			return fmt.Errorf(errJobNotFoundMessage, jobName, t.namespace)
		}
		return err
	}
	training.PrintTrainingJob(job, format, showEvent, showGPU)
	return nil
}

// List returns all training jobs
func (t *TrainingJobClient) List(allNamespaces bool, trainingType types.TrainingJobType) ([]*types.TrainingJobInfo, error) {
	jobs, err := training.ListTrainingJobs(t.namespace, allNamespaces, trainingType)
	if err != nil {
		return nil, err
	}
	jobInfos := []*types.TrainingJobInfo{}
	for _, job := range jobs {
		jobInfos = append(jobInfos, training.BuildJobInfo(job, true))
	}
	return jobInfos, nil
}

// ListAndPrint lists and prints the job informations
func (t *TrainingJobClient) ListAndPrint(allNamespaces bool, format string, trainingType types.TrainingJobType) error {
	if utils.TransferPrintFormat(format) == types.UnknownFormat {
		return fmt.Errorf("Unknown output format,only support:[wide|json|yaml]")
	}
	jobs, err := training.ListTrainingJobs(t.namespace, allNamespaces, trainingType)
	if err != nil {
		return err
	}
	training.DisplayTrainingJobList(jobs, format, allNamespaces)
	return nil
}

// Logs returns the training job log
func (t *TrainingJobClient) Logs(jobName string, jobType types.TrainingJobType, args *types.LogArgs) error {
	args.Namespace = t.namespace
	args.JobName = jobName
	return training.AcceptJobLog(jobName, jobType, args)
}

// Delete deletes the target training job
func (t *TrainingJobClient) Delete(jobType types.TrainingJobType, jobNames ...string) error {
	for _, jobName := range jobNames {
		err := training.DeleteTrainingJob(jobName, t.namespace, jobType)
		if err != nil {
			if err == types.ErrTrainingJobNotFound {
				return nil
			}
			return err
		}
	}
	return nil
}

// LogViewer returns the log viewer
func (t *TrainingJobClient) LogViewer(jobName string, jobType types.TrainingJobType) ([]string, error) {
	job, err := training.SearchTrainingJob(jobName, t.namespace, jobType)
	if err != nil {
		return nil, err
	}
	return job.GetJobDashboards(t.configer.GetClientSet(), t.namespace, t.arenaSystemNamespace)
}

// Prune cleans the not running training jobs
func (t *TrainingJobClient) Prune(allNamespaces bool, since time.Duration) error {
	return training.PruneTrainingJobs(t.namespace, allNamespaces, since)
}

func (t *TrainingJobClient) Top(args []string, namespace string, allNamespaces bool, jobType types.TrainingJobType, instanceName string, notStop bool, format types.FormatStyle) error {
	return training.TopTrainingJobs(args, namespace, allNamespaces, jobType, instanceName, notStop, format)
}
