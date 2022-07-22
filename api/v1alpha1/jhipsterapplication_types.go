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
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type AppType string

const (
	Monolith     AppType = "monolith"
	Microservice AppType = "microservice"
	Gateway      AppType = "gateway"
)

type DbType string

const (
	PostgreSQL DbType = "postgresql"
	MySQL      DbType = "mysql"
	NoDb       DbType = "none"
)

type DatabaseConfiguration struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
}

/*
   	name String
       image String
   	ingressDomain String
       ingressClassName String
       tls Boolean
       appType AppType
       env String
       useAntiAffinity Boolean
       profilesActive String
       useExternalDb Boolean
       databaseType DatabaseType
       serviceType String
*/

// JHipsterApplicationSpec defines the desired state of JHipsterApplication
type JHipsterApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// the name of the application to deploy
	Name string `json:"name"`
	// the image of the container
	Image string `json:"image"`
	// for gateways and monoliths, the ingress domain
	IngressDomain *string `json:"ingressDomain,omitempty"`

	// for gateways and monoliths, the class name to be used. Default is 'nginx'

	//+kubebuilder:default:="nginx"
	IngressClassName *string `json:"ingressClassName,omitempty"`

	// if true, the nginx is generated to use TLS as-well
	//+kubebuilder:default:=false
	Tls bool `json:"tls"`

	// the type of the application, one of 'microservice', 'gateway' or 'monolith'
	//+kubebuilder:default:="microservice"
	AppType AppType `json:"appType"`

	// if true, deployments whill have PodAntiAffinity configured
	//+kubebuilder:default:=true
	UsePodAntiAffinity bool `json:"usePodAntiAffinity"`

	// the active Spring profiles for this application
	//+kubebuilder:default:="prod"
	ProfilesActive string `json:"profilesActive"`

	// if true, the operator doesn't generate a database workload, but assumes user provides the connection
	//+kubebuilder:default:=false
	UseExternalDb bool `json:"useExternalDb"`

	// if useExternalDb is true, this credentials are used to connect to database. Otherwise, if not empty, these
	// credentials are used for the managed database
	DatabaseConfiguration DatabaseConfiguration `json:"databaseConfiguration,omitempty"`

	// the type of the applications database, default is 'postgresql'
	//+kubebuilder:default:="postgresql"
	DatabaseType DbType `json:"databaseType"`

	// the service type used for the applications service
	ServiceType *v1.ServiceType `json:"serviceType,omitempty"`

	// overrides for the deployment of the application
	//+kubebuilder:skipversion
	ApplicationOverrides v12.DeploymentSpec `json:"applicationOverrides,omitempty"`

	// overrides for the deployment of the managed database
	//+kubebuilder:skipversion
	DatabaseOverrides v12.DeploymentSpec `json:"databaseOverrides,omitempty"`

	// overrides for the generated ingress resource
	//+kubebuilder:skipversion
	IngressOverrides v13.IngressSpec `json:"ingressOverrides,omitempty"`

	// overrides for the generated service resource
	//+kubebuilder:skipversion
	ServiceOverrides v1.ServiceSpec `json:"serviceOverrides,omitempty"`
}

// JHipsterApplicationStatus defines the observed state of JHipsterApplication
type JHipsterApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=jhiapp;jhiapps

// JHipsterApplication is the Schema for the jhipsterapplications API
type JHipsterApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JHipsterApplicationSpec   `json:"spec,omitempty"`
	Status JHipsterApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JHipsterApplicationList contains a list of JHipsterApplication
type JHipsterApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JHipsterApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JHipsterApplication{}, &JHipsterApplicationList{})
}
