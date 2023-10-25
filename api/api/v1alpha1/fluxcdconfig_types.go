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

package v1alpha1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FluxCDConfigSpec defines the desired state of FluxCDConfig
type FluxCDConfigSpec struct {
	InstallCRDs               bool                       `json:"installCRDs"`
	Multitenancy              Multitenancy               `json:"multitenancy"`
	ClusterDomain             string                     `json:"clusterDomain"`
	Cli                       CliSpec                    `json:"cli"`
	HelmController            ControllerSpec             `json:"helmController"`
	ImageAutomationController ControllerSpec             `json:"imageAutomationController"`
	ImageReflectionController ControllerSpec             `json:"imageReflectionController"`
	KustomizeController       KustomizeControllerSpec    `json:"kustomizeController"`
	NotificationController    NotificationControllerSpec `json:"notificationController"`
	SourceController          SourceControllerSpec       `json:"sourceController"`
	Policies                  Policies                   `json:"policies"`
	Rbac                      Rbac                       `json:"rbac"`
	LogLevel                  string                     `json:"logLevel"`
	WatchAllNamespaces        bool                       `json:"watchAllNamespaces"`
	//+optional
	ImagePullSecrets []core.LocalObjectReference `json:"imagePullSecrets"`
	//+optional
	ExtraObjects []runtime.RawExtension `json:"extraObjects"`
	Prometheus   PrometheusSpec         `json:"prometheus"`
}

type Multitenancy struct {
	Enabled bool `json:"enabled"`
	// +kubebuilder:default=default
	DefaultServiceAccount string `json:"defaultServiceAccount"`
	// +kubebuilder:default=true
	Privileged bool `json:"privileged"`
}

type CliSpec struct {
	//+kubebuilder:default=ghcr.io/fluxcd/flux-cli
	Image string `json:"image"`
	//+kubebuilder:default=v2.1.1
	Tag            string                `json:"tag"`
	NodeSelector   map[string]string     `json:"nodeSelector"`
	Affinity       core.Affinity         `json:"affinity"`
	Tolerations    []core.Toleration     `json:"tolerations"`
	Annotations    map[string]string     `json:"annotations"`
	ServiceAccount CliServiceAccountSpec `json:"serviceAccount"`
}

type CliServiceAccountSpec struct {
	//+kubebuilder:default=true
	Automount bool `json:"automount"`
}

type ControllerSpec struct {
	Create    bool                 `json:"create"`
	Image     string               `json:"image"`
	Tag       string               `json:"tag"`
	Resources ResourceRequirements `json:"resources"`

	//+Optional
	PriorityClassName string            `json:"priorityClassName"`
	Annotations       map[string]string `json:"annotations"`
	//+Optional
	Labels    map[string]string `json:"labels"`
	Container ContainerSpec     `json:"container"`
	//+Optional
	ExtraEnv       []core.EnvVar      `json:"extraEnv"`
	ServiceAccount ServiceAccountSpec `json:"serviceAccount"`

	//+Optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy"`
	//+Optional
	NodeSelector map[string]string `json:"nodeSelector"`
	//+Optional
	Affinity core.Affinity `json:"affinity"`
	//+Optional
	Tolerations []core.Toleration `json:"tolerations"`
}

type KustomizeControllerSpec struct {
	Create    bool                 `json:"create"`
	Image     string               `json:"image"`
	Tag       string               `json:"tag"`
	Resources ResourceRequirements `json:"resources"`
	//+Optional
	PriorityClassName string            `json:"priorityClassName"`
	Annotations       map[string]string `json:"annotations"`
	//+Optional
	Labels    map[string]string `json:"labels"`
	Container ContainerSpec     `json:"container"`
	EnvFrom   EnvFromSource     `json:"envFrom"`
	//+Optional
	ExtraEnv []core.EnvVar `json:"extraEnv"`
	//+Optional
	ExtraSecretMounts []core.VolumeMount `json:"extraSecretMounts"`
	ServiceAccount    ServiceAccountSpec `json:"serviceAccount"`

	//+Optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy"`
	Secret          SecretSpec      `json:"secret"`
	//+Optional
	NodeSelector map[string]string `json:"nodeSelector"`
	//+Optional
	Affinity core.Affinity `json:"affinity"`
	//+Optional
	Tolerations []core.Toleration `json:"tolerations"`
}

type SecretSpec struct {
	Create bool              `json:"create"`
	Name   string            `json:"name"`
	Data   map[string]string `json:"data"`
}

type NotificationControllerSpec struct {
	Create    bool                 `json:"create"`
	Image     string               `json:"image"`
	Tag       string               `json:"tag"`
	Resources ResourceRequirements `json:"resources"`
	//+Optional
	PriorityClassName string            `json:"priorityClassName"`
	Annotations       map[string]string `json:"annotations"`
	//+Optional
	Labels    map[string]string `json:"labels"`
	Container ContainerSpec     `json:"container"`
	//+Optional
	ExtraEnv       []core.EnvVar      `json:"extraEnv"`
	ServiceAccount ServiceAccountSpec `json:"serviceAccount"`

	//+Optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy     `json:"imagePullPolicy"`
	Service         ServiceSpec         `json:"service"`
	WebhookReceiver WebhookReceiverSpec `json:"webhookReceiver"`
	//+Optional
	NodeSelector map[string]string `json:"nodeSelector"`
	//+Optional
	Affinity core.Affinity `json:"affinity"`
	//+Optional
	Tolerations []core.Toleration `json:"tolerations"`
}

type WebhookReceiverSpec struct {
	Service ServiceSpec `json:"service"`
}

type SourceControllerSpec struct {
	Create    bool                 `json:"create"`
	Image     string               `json:"image"`
	Tag       string               `json:"tag"`
	Resources ResourceRequirements `json:"resources"`
	//+Optional
	PriorityClassName string `json:"priorityClassName"`
	//+Optional
	Annotations map[string]string `json:"annotations"`
	//+Optional
	Labels    map[string]string `json:"labels"`
	Container ContainerSpec     `json:"container"`
	//+Optional
	ExtraEnv       []core.EnvVar      `json:"extraEnv"`
	ServiceAccount ServiceAccountSpec `json:"serviceAccount"`

	//+Optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy"`
	Service         ServiceSpec     `json:"service"`
	//+Optional
	NodeSelector map[string]string `json:"nodeSelector"`
	//+Optional
	Affinity core.Affinity `json:"affinity"`
	//+Optional
	Tolerations []core.Toleration `json:"tolerations"`
}

type ServiceSpec struct {
	//+Optional
	Labels map[string]string `json:"labels"`
	//+Optional
	Annotations map[string]string `json:"annotations"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits core.ResourceList `json:"limits"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests core.ResourceList `json:"requests"`
}

type EnvFromSource struct {
	Map    LocalObjectReference `json:"map"`
	Secret LocalObjectReference `json:"secret"`
}

type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// TODO: Add other useful fields. apiVersion, kind, uid?
	// +optional
	Name string `json:"name"`
}

type ContainerSpec struct {
	//+Optional
	AdditionalArgs []string `json:"additionalArgs"`
}

type ServiceAccountSpec struct {
	Create      bool              `json:"create"`
	Automount   bool              `json:"automount"`
	Annotations map[string]string `json:"annotations"`
}

type Policies struct {
	Create bool `json:"create"`
}

type Rbac struct {
	Create            bool `json:"create"`
	CreateAggregation bool `json:"createAggregation"`
}

type PrometheusSpec struct {
	PodMonitor PodMonitorSpec `json:"podMonitor"`
}

type PodMonitorSpec struct {
	Create bool `json:"create"`

	PodMetricsEndpoints []MetricsEndpoints `json:"podMetricsEndpoints"`
}

type MetricsEndpoints struct {
	Port string `json:"port"`

	Relabelings []Relabeling `json:"relabelings"`
}

type Relabeling struct {
	SourceLabels []string `json:"sourceLabels"`

	Action string `json:"action"`

	Regex string `json:"regex"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluxCDConfig is the Schema for the fluxcdconfigs API
type FluxCDConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec FluxCDConfigSpec `json:"spec"`
}

//+kubebuilder:object:root=true

// FluxCDConfigList contains a list of FluxCDConfig
type FluxCDConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FluxCDConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FluxCDConfig{}, &FluxCDConfigList{})
}
