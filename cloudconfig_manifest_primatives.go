package enaml

type CloudConfigManifest struct {
	AZs         []AZ                `yaml:"azs,omitempty"`
	VMTypes     []VMType            `yaml:"vm_types,omitempty"`
	DiskTypes   []DiskPool          `yaml:"disk_types,omitempty"`
	Networks    []DeploymentNetwork `yaml:"networks,omitempty"`
	Compilation *Compilation        `yaml:"compilation,omitempty"`
}

type VMType struct {
	Name            string      `yaml:"name:omitempty"`
	CloudProperties interface{} `yaml:"cloud_properties,omitempty"`
}

type AZ struct {
	Name            string      `yaml:"name:omitempty"`
	CloudProperties interface{} `yaml:"cloud_properties,omitempty"`
}
