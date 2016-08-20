package awscloudproperties

type AZ struct {
	AvailabilityZoneName string   `yaml:"availability_zone,omitempty"`
	SecurityGroups       []string `yaml:"security_groups,omitempty"`
}

type Network struct {
	Subnet         string   `yaml:"subnet,omitempty"`
	SecurityGroups []string `yaml:"security_groups,omitempty"`
}

type VMType ResourcePool
type ResourcePool struct {
	InstanceType          string        `yaml:"instance_type,omitempty"`          // [String, required]: Type of the instance. Example: m3.medium.
	AvailabilityZone      string        `yaml:"availability_zone,omitempty"`      // [String, required]: Availability zone to use for creating instances. Example: us-east-1a.
	SecurityGroups        []string      `yaml:"security_groups,omitempty"`        // [Array, optional]: Array of security groups to apply for all VMs that are in this resource pool. Defaults to security groups specified by default_security_groups in the global CPI settings unless security groups are specified on one of the VM networks. Security groups can be specified either on a resource pool or on a network. Available in v46+.
	KeyName               string        `yaml:"key_name,omitempty"`               // [String, optional]: Key pair name. Defaults to key pair name specified by default_key_name in global CPI settings. Example: bosh.
	SpotBidPrice          *float64      `yaml:"spot_bid_price,omitempty"`         // [Float, optional]: Bid price in dollars for AWS spot instance. Using this option will slow down VM creation. Example: 0.03.
	SpotOnDemandFallback bool          `yaml:"spot_ondemand_fallback,omitempty"` // [Boolean, optional]: Set to true to use an on demand instance if a spot instance is not available during VM creation. Defaults to false. Available in v36.
	ELBs                  []string      `yaml:"elbs,omitempty"`                   // [Array, optional]: Array of ELB names that should be attached to created VMs. Example: [prod-elb]. Default is [].
	IamInstanceProfile    string        `yaml:"iam_instance_profile,omitempty"`   // [String, optional]: Name of an IAM instance profile. Example: director.
	PlacementGroup        string        `yaml:"placement_group,omitempty"`        // [String, optional]: Name of a placement group. Example: my-group.
	Tenancy               string        `yaml:"tenancy,omitempty"`                // [String, optional]: VM tenancy configuration. Example: dedicated. Default is default.
	RawInstanceStorage    bool          `yaml:"raw_instance_storage,omitempty"`   // [Boolean, optional]: Exposes all available instance storage via labeled disks. Defaults to false.
	EphemeralDisk         EphemeralDisk `yaml:"ephemeral_disk,omitempty,flow"`    //EBS backed ephemeral disk of custom size for when instance storage is not large enough or not available for selected instance type.
	RootDisk              RootDisk      `yaml:"root_disk,omitempty"`              // [Hash, optional]: EBS backed root disk of custom size.
}

type RootDisk struct {
	DiskSize int    `yaml:"size,omitempty"` //[Integer, required]: Specifies the disk size in megabytes.
	DiskType string `yaml:"type,omitempty"` //[String, optional]: Type of the disk: standard, gp2. Defaults to standard.
}

type EphemeralDisk struct {
	Size     int    `yaml:"size,omitempty"` //size [Integer, required]: Specifies the disk size in megabytes.
	DiskType string `yaml:"type,omitempty"` // [String, optional]: Type of the disk: standard, gp2. Defaults to standard.
	/*
		standard stands for EBS magnetic drives
		gp2 stands for EBS general purpose drives (SSD)
		io1 stands for EBS provisioned IOPS drives (SSD)
	*/
	IOPs int `yaml:"iops,omitempty"` // [Integer, optional]: Specifies the number of I/O operations per second to provision for the drive.
	//Only valid for io1 type drive.
	//Required when io1 type drive is specified.
}
