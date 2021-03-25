// Copyright © 2021 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +kubebuilder:object:generate=true

package backup

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/mitchellh/mapstructure"
	v1 "k8s.io/api/core/v1"
)

type ServiceSpec struct {
	ChartValues ValueOverrides `json:"chartValues"`
}

// +kubebuilder:object:generate=true
type ValueOverrides struct {
	Configuration   Configuration          `json:"Configuration"`
	Credentials     Credentials            `json:"Credentials"`
	Image           Image                  `json:"Image"`
	RBAC            Rbac                   `json:"Rbac"`
	InitContainers  []v1.Container         `json:"initContainers"`
	CleanUpCRDs     bool                   `json:"cleanUpCRDs"`
	ServiceAccount  ServiceAccount         `json:"ServiceAccount"`
	SecurityContext SecurityContext        `json:"SecurityContext"`
	Affinity        map[string]interface{} `json:"affinity"`
}

// +kubebuilder:object:generate=true
type SecurityContext struct {
	FsGroup int `json:"fsGroup"`
}

// +kubebuilder:object:generate=true
type ServiceAccount struct {
	Server Server `json:"Server"`
}

// +kubebuilder:object:generate=true
type Server struct {
	Create      bool              `json:"create"`
	Name        string            `json:"name"`
	Annotations map[string]string `json:"annotations"`
}

// +kubebuilder:object:generate=true
type Rbac struct {
	Create bool `json:"create"`
}

// +kubebuilder:object:generate=true
type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}

// +kubebuilder:object:generate=true
type Configuration struct {
	Provider               string                 `json:"provider"`
	VolumeSnapshotLocation VolumeSnapshotLocation `json:"VolumeSnapshotLocation"`
	BackupStorageLocation  BackupStorageLocation  `json:"BackupStorageLocation"`
	RestoreOnlyMode        bool                   `json:"restoreOnlyMode"`
	LogLevel               string                 `json:"logLevel"`
}

// +kubebuilder:object:generate=true
type Credentials struct {
	ExistingSecret string `json:"existingSecret"`
}

// +kubebuilder:object:generate=true
type VolumeSnapshotLocation struct {
	Name     string                       `json:"name"`
	Provider string                       `json:"provider"`
	Config   VolumeSnapshotLocationConfig `json:"config,omitempty"`
}

type VolumeSnapshotLocationConfig struct {
	Region        string `json:"region,omitempty"`
	Profile       string `json:"profile,omitempty"`
	ApiTimeout    string `json:"apiTimeout,omitempty"`
	ResourceGroup string `json:"resourceGroup,omitempty"`
}

// +kubebuilder:object:generate=true
type BackupStorageLocation struct {
	Name     string                      `json:"name"`
	Provider string                      `json:"provider"`
	Bucket   string                      `json:"bucket"`
	Prefix   string                      `json:"prefix"`
	Config   BackupStorageLocationConfig `json:"config,omitempty"`
}

type BackupStorageLocationConfig struct {
	Region                  string `json:"region,omitempty"`
	Profile                 string `json:"profile,omitempty"`
	S3ForcePathStyle        string `json:"s3ForcePathStyle,omitempty"`
	S3Url                   string `json:"s3Url,omitempty"`
	KMSKeyId                string `json:"kmsKeyId,omitempty"`
	ResourceGroup           string `json:"resourceGroup,omitempty"`
	StorageAccount          string `json:"storageAccount,omitempty"`
	StorageAccountKeyEnvVar string `json:"storageAccountKeyEnvVar,omitempty"`
}

func BindIntegratedServiceSpec(spec map[string]interface{}) (ServiceSpec, error) {
	var boundSpec ServiceSpec
	if err := mapstructure.Decode(spec, &boundSpec); err != nil {
		return boundSpec, errors.WrapIf(err, "failed to bind integrated service spec")
	}
	return boundSpec, nil
}

type requiredStringFieldError struct {
	fieldName string
}

func (e requiredStringFieldError) Error() string {
	return fmt.Sprintf("%s must be specified and cannot be empty", e.fieldName)
}
