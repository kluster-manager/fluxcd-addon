/*
Copyright AppsCode Inc. and Contributors.

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
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FluxCDConfigSpec defines the desired state of FluxCDConfig
type FluxCDConfigSpec struct {
	InstallCRDs bool `json:"installCRDs"`
	// +optional
	CRDs CRDsSpec `json:"crds,omitempty"`
	// +optional
	Multitenancy Multitenancy `json:"multitenancy,omitempty"`
	// +optional
	ClusterDomain string `json:"clusterDomain,omitempty"`
	// +optional
	Cli CliSpec `json:"cli,omitempty"`
	// +optional
	HelmController ControllerSpec `json:"helmController,omitempty"`
	// +optional
	ImageAutomationController ControllerSpec `json:"imageAutomationController,omitempty"`
	// +optional
	ImageReflectionController ControllerSpec `json:"imageReflectionController,omitempty"`
	// +optional
	KustomizeController KustomizeControllerSpec `json:"kustomizeController,omitempty"`
	// +optional
	NotificationController NotificationControllerSpec `json:"notificationController,omitempty"`
	// +optional
	SourceController SourceControllerSpec `json:"sourceController,omitempty"`
	// +optional
	Policies Policies `json:"policies,omitempty"`
	// +optional
	Rbac Rbac `json:"rbac,omitempty"`
	// +optional
	LogLevel string `json:"logLevel,omitempty"`
	// +optional
	WatchAllNamespaces bool `json:"watchAllNamespaces,omitempty"`
	// +optional
	ImagePullSecrets []core.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// +optional
	ExtraObjects []runtime.RawExtension `json:"extraObjects,omitempty"`
	// +optional
	Prometheus PrometheusSpec `json:"prometheus,omitempty"`
}

type CRDsSpec struct {
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Multitenancy struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	DefaultServiceAccount string `json:"defaultServiceAccount,omitempty"`
	// +optional
	Privileged bool `json:"privileged,omitempty"`
}

type CliSpec struct {
	// +optional
	Image string `json:"image,omitempty"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// +optional
	Affinity core.Affinity `json:"affinity,omitempty"`
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	ServiceAccount CliServiceAccountSpec `json:"serviceAccount,omitempty"`
}

type CliServiceAccountSpec struct {
	// +optional
	Automount bool `json:"automount,omitempty"`
}

type ControllerSpec struct {
	Create bool `json:"create"`
	// +optional
	Image string `json:"image,omitempty"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Container ContainerSpec `json:"container,omitempty"`
	// +optional
	ExtraEnv []core.EnvVar `json:"extraEnv,omitempty"`
	// +optional
	ServiceAccount ServiceAccountSpec `json:"serviceAccount,omitempty"`
	// +optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy,omitempty"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// +optional
	Affinity core.Affinity `json:"affinity,omitempty"`
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

type KustomizeControllerSpec struct {
	Create bool `json:"create"`
	// +optional
	Image string `json:"image,omitempty"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Container ContainerSpec `json:"container,omitempty"`
	// +optional
	EnvFrom EnvFromSource `json:"envFrom,omitempty"`
	// +optional
	ExtraEnv []core.EnvVar `json:"extraEnv,omitempty"`
	// +optional
	ExtraSecretMounts []core.VolumeMount `json:"extraSecretMounts,omitempty"`
	// +optional
	ServiceAccount ServiceAccountSpec `json:"serviceAccount,omitempty"`
	// +optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy,omitempty"`
	// +optional
	Secret SecretSpec `json:"secret,omitempty"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// +optional
	Affinity core.Affinity `json:"affinity,omitempty"`
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

type SecretSpec struct {
	Create bool `json:"create,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	Data map[string]string `json:"data,omitempty"`
}

type NotificationControllerSpec struct {
	Create bool `json:"create"`
	// +optional
	Image string `json:"image,omitempty"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Container ContainerSpec `json:"container,omitempty"`
	// +optional
	ExtraEnv []core.EnvVar `json:"extraEnv,omitempty"`
	// +optional
	ServiceAccount ServiceAccountSpec `json:"serviceAccount,omitempty"`
	// +optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy,omitempty"`
	// +optional
	Service ServiceSpec `json:"service,omitempty"`
	// +optional
	WebhookReceiver WebhookReceiverSpec `json:"webhookReceiver,omitempty"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// +optional
	Affinity core.Affinity `json:"affinity,omitempty"`
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

type WebhookReceiverSpec struct {
	// +optional
	Service ServiceSpec `json:"service,omitempty"`
	// +optional
	Ingress IngressSpec `json:"ingress,omitempty"`
}

type IngressSpec struct {
	Create      bool                    `json:"create,omitempty"`
	Annotations map[string]string       `json:"annotations,omitempty"`
	Labels      map[string]string       `json:"labels,omitempty"`
	Hosts       []IngressRule           `json:"hosts,omitempty"`
	TLS         []networking.IngressTLS `json:"tls,omitempty"`
}

type IngressRule struct {
	Host  string            `json:"host,omitempty"`
	Paths []HTTPIngressPath `json:"paths,omitempty"`
}

type HTTPIngressPath struct {
	Path     string `json:"path,omitempty"`
	PathType string `json:"pathType,omitempty"`
}

type SourceControllerSpec struct {
	Create bool `json:"create"`
	// +optional
	Image string `json:"image,omitempty"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Container ContainerSpec `json:"container,omitempty"`
	// +optional
	ExtraEnv []core.EnvVar `json:"extraEnv,omitempty"`
	// +optional
	ServiceAccount ServiceAccountSpec `json:"serviceAccount,omitempty"`
	// +optional
	//+kubebuilder:validation:Enum=Always;Never;IfNotPresent;""
	ImagePullPolicy core.PullPolicy `json:"imagePullPolicy,omitempty"`
	// +optional
	Service ServiceSpec `json:"service,omitempty"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// +optional
	Affinity core.Affinity `json:"affinity,omitempty"`
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

type ServiceSpec struct {
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits core.ResourceList `json:"limits,omitempty"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests core.ResourceList `json:"requests,omitempty"`
}

type EnvFromSource struct {
	// +optional
	Map LocalObjectReference `json:"map,omitempty"`
	// +optional
	Secret LocalObjectReference `json:"secret,omitempty"`
}

type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// TODO: Add other useful fields. apiVersion, kind, uid?
	// +optional
	Name string `json:"name,omitempty"`
}

type ContainerSpec struct {
	// +optional
	AdditionalArgs []string `json:"additionalArgs,omitempty"`
}

type ServiceAccountSpec struct {
	Create bool `json:"create,omitempty"`
	// +optional
	Automount bool `json:"automount,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Policies struct {
	Create bool `json:"create,omitempty"`
}

type Rbac struct {
	Create bool `json:"create,omitempty"`
	// +optional
	CreateAggregation bool `json:"createAggregation,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

type PrometheusSpec struct {
	// +optional
	PodMonitor PodMonitorSpec `json:"podMonitor,omitempty"`
}

type PodMonitorSpec struct {
	Create bool `json:"create,omitempty"`
	// +optional
	PodMetricsEndpoints []MetricsEndpoints `json:"podMetricsEndpoints,omitempty"`
}

type MetricsEndpoints struct {
	// +optional
	Port string `json:"port,omitempty"`
	// +optional
	Relabelings []Relabeling `json:"relabelings,omitempty"`
}

type Relabeling struct {
	// +optional
	SourceLabels []string `json:"sourceLabels,omitempty"`
	// +optional
	Action string `json:"action,omitempty"`
	// +optional
	Regex string `json:"regex,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluxCDConfig is the Schema for the fluxcdconfigs API
type FluxCDConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec FluxCDConfigSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
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
