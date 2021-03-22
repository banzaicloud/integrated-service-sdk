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

package backup

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/mitchellh/mapstructure"
)

const (
	// supported DNS provider names
	dnsRoute53 = "route53"
	dnsAzure   = "azure"
	dnsGoogle  = "google"
	dnsBanzai  = "banzaicloud-dns"
)

// +kubebuilder:object:generate=true

type ServiceSpec struct {
	Cluster       ClusterConfig `json:"cluster" mapstructure:"cluster"`
	Bucket        BucketConfig `json:"bucket" mapstructure:"bucket"`
	SecretName    string `json:"secretName"`

	UseClusterSecret      bool `json:"UseClusterSecret"`
	ServiceAccountRoleARN string `json:"secretName"`
	UseProviderSecret     bool `json:"useProviderSecret"`
}

// +kubebuilder:object:generate=true

type ClusterConfig struct {
	Name         string `json:"name"`
	Provider     string `json:"provider"`
	Distribution string `json:"distribution"`
	Location     string `json:"location"`
	RBACEnabled  bool `json:"rbacEnabled"`
	ResourceGroup string `json:"resourceGroup"`
}


// +kubebuilder:object:generate=true

type BucketConfig struct {
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Provider string `json:"provider"`
	Location string `json:"location"`
	StorageAccount string `json:"storageAccount"`
	ResourceGroup  string `json:"resourceGroup"`
}


func (s ServiceSpec) Validate() error {
    var secretNameErr error
	if s.SecretName == "" {
		secretNameErr = requiredStringFieldError{fieldName: "secretName"}
	}
	return errors.Combine(s.Cluster.Validate(), s.Bucket.Validate(), secretNameErr)
}


func (s ClusterConfig) Validate() error {
	// TODO check all required fields.
	// resourceGroup required only for Azure
	if s.Name == "" {
		return requiredStringFieldError{fieldName: "cluster.name"}
	}
	return nil
}

func (s BucketConfig) Validate() error {
	// TODO check all required fields.
	// storageAccount, resourceGroup required only for Azure
	if s.Name == "" {
		return requiredStringFieldError{fieldName: "bucket.name"}
	}
	return nil
}

// +kubebuilder:object:generate=true

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
