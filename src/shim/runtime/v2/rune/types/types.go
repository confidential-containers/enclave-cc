package types

const (
	// PodSandbox identifies an infra container that will be used to create the pod
	PodSandbox = "sandbox"
	// PodContainer identifies a container that should be associated with an existing pod
	PodContainer = "container"

	// Default identifies the default namespace used by objects
	Default = "default"
	// KubeSystem identifies the namespace used by the Kubernetes system objects
	KubeSystem = "kube-system"
)
