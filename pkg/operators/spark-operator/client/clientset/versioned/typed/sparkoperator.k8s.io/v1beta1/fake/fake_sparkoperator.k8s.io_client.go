// Code generated by k8s code-generator DO NOT EDIT.

/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/kubeflow/arena/pkg/operators/spark-operator/client/clientset/versioned/typed/sparkoperator.k8s.io/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSparkoperatorV1beta1 struct {
	*testing.Fake
}

func (c *FakeSparkoperatorV1beta1) ScheduledSparkApplications(namespace string) v1beta1.ScheduledSparkApplicationInterface {
	return &FakeScheduledSparkApplications{c, namespace}
}

func (c *FakeSparkoperatorV1beta1) SparkApplications(namespace string) v1beta1.SparkApplicationInterface {
	return &FakeSparkApplications{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSparkoperatorV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
