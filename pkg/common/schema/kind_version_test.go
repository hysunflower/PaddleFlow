/*
Copyright (c) 2023 PaddlePaddle Authors. All Rights Reserve.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToKindGroupVersion(t *testing.T) {
	testCases := []struct {
		name        string
		clusterType string
		framework   Framework
		annotations map[string]string
		wantRes     KindGroupVersion
	}{
		{
			name:        "default paddle job kind group version",
			clusterType: KubernetesType,
			framework:   FrameworkPaddle,
			wantRes:     PaddleKindGroupVersion,
		},
		{
			name:        "custom paddle job kind group version",
			clusterType: KubernetesType,
			framework:   FrameworkPaddle,
			annotations: map[string]string{
				JobKindAnnotation:         "PaddleJob",
				JobGroupVersionAnnotation: "kubeflow.org/v1",
			},
			wantRes: KindGroupVersion{Kind: "PaddleJob", Group: "kubeflow.org", APIVersion: "v1"},
		},
		{
			name:        "custom paddle job kind group version failed",
			clusterType: KubernetesType,
			framework:   FrameworkPaddle,
			annotations: map[string]string{
				JobKindAnnotation:         "PaddleJob",
				JobGroupVersionAnnotation: "xx/kubeflow.org/v1",
			},
			wantRes: PaddleKindGroupVersion,
		},
		{
			name:        "custom tf job kind group version",
			clusterType: KubernetesType,
			framework:   FrameworkTF,
			annotations: map[string]string{
				JobKindAnnotation:         "TFJob",
				JobGroupVersionAnnotation: "v1",
			},
			wantRes: KindGroupVersion{Kind: "TFJob", Group: "", APIVersion: "v1"},
		},
		{
			name:        "unexcepted framework",
			clusterType: KubernetesType,
			framework:   Framework("xxx"),
			wantRes:     KindGroupVersion{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("%s\n", tc.wantRes)
			kgv := ToKindGroupVersion(tc.clusterType, tc.framework, tc.annotations)
			assert.Equal(t, tc.wantRes, kgv)
		})
	}
}

func TestGetJobType(t *testing.T) {
	testCases := []struct {
		name             string
		kindGroupVersion KindGroupVersion
		jobType          JobType
	}{
		{
			name:             "single job",
			kindGroupVersion: StandaloneKindGroupVersion,
			jobType:          TypeSingle,
		},
		{
			name:             "distributed job",
			kindGroupVersion: PaddleKindGroupVersion,
			jobType:          TypeDistributed,
		},
		{
			name:             "workflow job",
			kindGroupVersion: WorkflowKindGroupVersion,
			jobType:          TypeWorkflow,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jobType := GetJobType(tc.kindGroupVersion)
			assert.Equal(t, tc.jobType, jobType)
		})
	}
}

func TestGetJobFramework(t *testing.T) {
	testCases := []struct {
		name        string
		kindVersion KindGroupVersion
		framework   Framework
	}{
		{
			name:        "single job framework",
			kindVersion: StandaloneKindGroupVersion,
			framework:   FrameworkStandalone,
		},
		{
			name:        "workflow job framework",
			kindVersion: WorkflowKindGroupVersion,
			framework:   "",
		},
		{
			name:        "spark app framework",
			kindVersion: SparkKindGroupVersion,
			framework:   FrameworkSpark,
		},
		{
			name:        "paddle job framework",
			kindVersion: PaddleKindGroupVersion,
			framework:   FrameworkPaddle,
		},
		{
			name:        "pytorch job framework",
			kindVersion: PyTorchKindGroupVersion,
			framework:   FrameworkPytorch,
		},
		{
			name:        "tf job framework",
			kindVersion: TFKindGroupVersion,
			framework:   FrameworkTF,
		},
		{
			name:        "mxnet job framework",
			kindVersion: MXNetKindGroupVersion,
			framework:   FrameworkMXNet,
		},
		{
			name:        "mpi job framework",
			kindVersion: MPIKindGroupVersion,
			framework:   FrameworkMPI,
		},
		{
			name:        "ray job framework",
			kindVersion: RayKindGroupVersion,
			framework:   FrameworkRay,
		},
		{
			name:        "other job framework",
			kindVersion: KindGroupVersion{},
			framework:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fw := GetJobFramework(tc.kindVersion)
			assert.Equal(t, tc.framework, fw)
		})
	}
}