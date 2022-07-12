package types

const (
	// PodSandbox identifies an infra container that will be used to create the pod
	PodSandbox = "sandbox"

	// DefaultNS identifies the default namespace used by objects
	DefaultNS = "default"
	// KubeSystemNS identifies the namespace used by the Kubernetes system objects
	KubeSystemNS = "kube-system"
)
