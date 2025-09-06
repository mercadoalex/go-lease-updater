package types

import "time"

type OwnerReference struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Name       string `yaml:"name"`
	UID        string `yaml:"uid"`
}

type LeaseSpec struct {
	HolderIdentity       string    `yaml:"holderIdentity"`
	LeaseDurationSeconds int       `yaml:"leaseDurationSeconds"`
	RenewTime            time.Time `yaml:"renewTime"`
}

type Lease struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       LeaseSpec `yaml:"spec"`
}

type Metadata struct {
	CreationTimestamp string           `yaml:"creationTimestamp"`
	Name              string           `yaml:"name"`
	Namespace         string           `yaml:"namespace"`
	OwnerReferences   []OwnerReference `yaml:"ownerReferences"`
	ResourceVersion   string           `yaml:"resourceVersion"`
	UID               string           `yaml:"uid"`
}
