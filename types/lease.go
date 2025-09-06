package types

type LeaseList struct {
	Items []Lease `yaml:"items"`
}

type Lease struct {
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		RenewTime            string `yaml:"renewTime"`
		LeaseDurationSeconds int    `yaml:"leaseDurationSeconds"`
	} `yaml:"spec"`
}
