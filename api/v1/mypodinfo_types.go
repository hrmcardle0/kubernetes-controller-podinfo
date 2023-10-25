/*
Copyright 2023.

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

package v1

import (
	//v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MyPodinfoRedis specifies whether to use redis or not
type MyPodinfoRedis struct {

	// specficy whether redis is to be enabled
	Enabled string `json:"enabled,omitempty"`
}

// MyPodinfoImage defines the image to use
type MyPodinfoImage struct {

	// Image is the image to use
	Image string `json:"image,omitempty"`
	// Name is the name to give the container
	Name string `json:"name,omitempty"`
}

// MyPodinfoResource defines resource the pod will use
type MyPodinfoResource struct {

	// memoryLimit specifies the max memory the pod will use
	MemoryLimit string `json:"memoryLimit,omitempty"`

	// cpuRequest is the amount of CPU the pod is requesting
	CpuRequest string `json:"cpuRequest,omitempty"`
}

// MyPodinfoSpec defines the desired state of MyPodinfo
type MyPodinfoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// replicaCount is the number of pods to generate
	ReplicaCount int `json:"replicaCount,omitempty"`
	// resources holds information about the pods
	Resources MyPodinfoResource `json:"resources,omitempty"`
	// image is the contanier image to use
	Image MyPodinfoImage `json:"image,omitempty"`
	// redis is whether to enable redis or not
	Redis MyPodinfoRedis `json:"redis,omitempty"`
}

// MyPodinfoStatus defines the observed state of MyPodinfo
type MyPodinfoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyPodinfo is the Schema for the mypodinfoes API
type MyPodinfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyPodinfoSpec   `json:"spec,omitempty"`
	Status MyPodinfoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyPodinfoList contains a list of MyPodinfo
type MyPodinfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyPodinfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyPodinfo{}, &MyPodinfoList{})
}
