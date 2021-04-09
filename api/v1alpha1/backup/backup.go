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

type Provider string

const AWSProvider Provider = "aws"
const AzureProvider Provider = "azure"
const GCPProvider Provider = "gcp"

type ServiceSpec struct {
	ChartValues ValueOverrides `json:"chartValues"`
}

type ValueOverrides struct {
	Configuration   Configuration   `json:"configuration"`
	Credentials     Credentials     `json:"credentials"`
	RBAC            Rbac            `json:"rbac"`
	CleanUpCRDs     bool            `json:"cleanUpCRDs"`
	ServiceAccount  ServiceAccount  `json:"serviceAccount,omitempty"`
	SecurityContext SecurityContext `json:"securityContext,omitempty"`
	Affinity        v1.Affinity     `json:"affinity,omitempty"`
}

type SecurityContext struct {
	FsGroup int `json:"fsGroup"`
}

type ServiceAccount struct {
	Server Server `json:"server"`
}

type Annotations map[string]string

type Server struct {
	Create bool   `json:"create"`
	Name   string `json:"name,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Rbac struct {
	Create bool `json:"create"`
}

type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}

type Configuration struct {
	Provider               Provider               `json:"provider"`
	VolumeSnapshotLocation VolumeSnapshotLocation `json:"volumeSnapshotLocation"`
	BackupStorageLocation  BackupStorageLocation  `json:"backupStorageLocation"`
	RestoreOnlyMode        bool                   `json:"restoreOnlyMode"`
	LogLevel               string                 `json:"logLevel"`
}

type Credentials struct {
	ExistingSecret string `json:"existingSecret"`
}

type VolumeSnapshotLocation struct {
	Name     string                       `json:"name"`
	Provider Provider                     `json:"provider"`
	Config   VolumeSnapshotLocationConfig `json:"config,omitempty"`
}

type VolumeSnapshotLocationConfig struct {
	Region        string `json:"region,omitempty"`
	Profile       string `json:"profile,omitempty"`
	ApiTimeout    string `json:"apiTimeout,omitempty"`
	ResourceGroup string `json:"resourceGroup,omitempty"`
}

type BackupStorageLocation struct {
	Name     string                      `json:"name"`
	Provider Provider                    `json:"provider"`
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

func (s ServiceSpec) Validate() error {
	var errs error

	if s.ChartValues.Credentials.ExistingSecret == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.credentials.existingSecret"})
	}

	if s.ChartValues.Configuration.Provider == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.provider"})
	}

	if err := validateProvider(s.ChartValues.Configuration.Provider); err != nil {
		errs = errors.Append(errs, err)
	}

	if s.ChartValues.Configuration.BackupStorageLocation.Provider == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.backupStorageLocation.provider"})
	}

	if err := validateProvider(s.ChartValues.Configuration.BackupStorageLocation.Provider); err != nil {
		errs = errors.Append(errs, err)
	}

	if s.ChartValues.Configuration.BackupStorageLocation.Name == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.backupStorageLocation.name"})
	}

	if s.ChartValues.Configuration.BackupStorageLocation.Bucket == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.backupStorageLocation.bucket"})
	}

	if s.ChartValues.Configuration.BackupStorageLocation.Prefix == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.backupStorageLocation.prefix"})
	}

	if s.ChartValues.Configuration.VolumeSnapshotLocation.Name == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.volumeSnapshotLocation.name"})
	}

	if s.ChartValues.Configuration.VolumeSnapshotLocation.Provider == "" {
		errs = errors.Append(errs, requiredStringFieldError{fieldName: "chartValues.configuration.volumeSnapshotLocation.provider"})
	}

	if err := validateProvider(s.ChartValues.Configuration.VolumeSnapshotLocation.Provider); err != nil {
		errs = errors.Append(errs, err)
	}

	return errors.Combine(errs)
}

type requiredStringFieldError struct {
	fieldName string
}

func (e requiredStringFieldError) Error() string {
	return fmt.Sprintf("%s must be specified and cannot be empty", e.fieldName)
}

type invalidProviderFieldError struct {
	fieldName string
}

func (e invalidProviderFieldError) Error() string {
	return fmt.Sprintf("%s is not a valid provider", e.fieldName)
}

func validateProvider(provider Provider) error {
	switch provider {
	case AWSProvider:
	case GCPProvider:
	case AzureProvider:
	default:
		return invalidProviderFieldError{fieldName: "chartValues.configuration.volumeSnapshotLocation.provider"}
	}
	return nil
}
