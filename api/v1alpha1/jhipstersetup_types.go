/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type ServiceDiscoveryType string

const (
	Consul ServiceDiscoveryType = "consul"
	Eureka ServiceDiscoveryType = "eureka"
	None   ServiceDiscoveryType = "none"
)

// JHipsterSetupSpec defines the desired state of JHipsterSetup
type JHipsterSetupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:default:="consul"
	// sets the service discovery solution for
	ServiceDiscoveryType ServiceDiscoveryType `json:"serviceDiscoveryType"`
	//+kubebuilder:default:=true
	UseMonitoring bool `json:"useMonitoring"`
	//+kubebuilder:default:=true
	UseDynamicStorage bool `json:"useDynamicStorage"`

	// the storage class name used for managed stuff
	StorageClassName *string `json:"storageClassName,omitempty"`

	//+kubebuilder:default:=false
	// if true, the operator installs istio for this setup
	// todo: implement istio. Currently it doesn't work
	UseIstio bool `json:"useIstio,omitempty"`
}

// JHipsterSetupStatus defines the observed state of JHipsterSetup
type JHipsterSetupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=jhisetup;jhisetups

// JHipsterSetup is the Schema for the jhipstersetups API
type JHipsterSetup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JHipsterSetupSpec   `json:"spec,omitempty"`
	Status JHipsterSetupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JHipsterSetupList contains a list of JHipsterSetup
type JHipsterSetupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JHipsterSetup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JHipsterSetup{}, &JHipsterSetupList{})
}
