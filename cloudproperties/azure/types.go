package azurecloudproperties

type ResourcePool struct {
	InstanceType string `yaml:"instance_type,omitempty"` // [String, required]: Type of the instance. Example: m3.medium.
}

type Network struct {
	VnetName   string `yaml:"virtual_network_name,omitempty"` //: VNET-NAME # <--- Replace
	SubnetName string `yaml:"subnet_name,omitempty"`          //: SUBNET-NAME # <--- Replace
}
