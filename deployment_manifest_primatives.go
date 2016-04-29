package enaml

type Deployment interface {
	GetDeployment() DeploymentManifest
}

type CloudConfig interface {
	GetManifest() CloudConfigManifest
}

type DeploymentManifest struct {
	Name           string              `yaml:"name"`
	DirectorUUID   string              `yaml:"director_uuid,omitempty"`
	Releases       []Release           `yaml:"releases,omitempty"`
	Stemcells      []Stemcell          `yaml:"stemcells,omitempty"`
	InstanceGroups []*InstanceGroup    `yaml:"instance_groups,omitempty"`
	Networks       []DeploymentNetwork `yaml:"networks,omitempty"`
	ResourcePools  []ResourcePool      `yaml:"resource_pools,omitempty"`
	DiskPools      []DiskPool          `yaml:"disk_pools,omitempty"`
	Compilation    *Compilation        `yaml:"compilation,omitempty"`
	Update         Update              `yaml:"update,omitempty"`
	Jobs           []Job               `yaml:"jobs,omitempty"`
	Properties     Properties          `yaml:"properties,omitempty"`
	CloudProvider  *CloudProvider      `yaml:"cloud_provider,omitempty"`
}

type InstanceGroup struct {
	Name               string        `yaml:"name"`
	ResourcePool       string        `yaml:"resource_pool,omitempty"`
	PersistentDisk     int           `yaml:"persistent_disk,omitempty"`
	PersistentDiskType string        `yaml:"persistent_disk_type,omitempty"`
	Instances          int           `yaml:"instances"`
	VMType             string        `yaml:"vm_type,omitempty"`
	Stemcell           string        `yaml:"stemcell,omitempty"`
	AZs                []string      `yaml:"azs,flow,omitempty"`
	Networks           []Network     `yaml:"networks,flow"`
	Jobs               []InstanceJob `yaml:"jobs"`
}

type InstanceJob struct {
	Name       string      `yaml:"name"`
	Release    string      `yaml:"release"`
	Properties interface{} `yaml:"properties"`
}

type DeploymentNetwork interface{}

type Release struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version,omitempty"`
	URL     string `yaml:"url,omitempty"`
	SHA1    string `yaml:"sha1,omitempty"`
}

func NewVIPNetwork(name string) VIPNetwork {
	return VIPNetwork{
		Name: name,
		Type: "vip",
	}
}

type VIPNetwork struct {
	Name            string          `yaml:"name"`
	Type            string          `yaml:"type"`
	CloudProperties CloudProperties `yaml:"cloud_properties,omitempty"`
}

type DynamicNetwork struct {
	Name            string          `yaml:"name"`
	Type            string          `yaml:"type"`
	DNS             []string        `yaml:"dns"`
	CloudProperties CloudProperties `yaml:"cloud_properties"`
}

func NewManualNetwork(name string) ManualNetwork {
	return ManualNetwork{
		Name: name,
		Type: "manual",
	}
}

type ManualNetwork struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Subnets []Subnet `yaml:"subnets"`
}

func (s *ManualNetwork) AddSubnet(subnet Subnet) (err error) {
	s.Subnets = append(s.Subnets, subnet)
	return
}

func NewSubnet(octet string, azName string) Subnet {
	return Subnet{
		Range:   octet + ".0/24",
		Gateway: octet + ".1",
		AZ:      azName,
	}
}

type Subnet struct {
	Range           string          `yaml:"range,omitempty"`
	Gateway         string          `yaml:"gateway,omitempty"`
	DNS             []string        `yaml:"dns,omitempty"`
	Reserved        []string        `yaml:"reserved,omitempty"`
	Static          []string        `yaml:"static,omitempty"`
	AZ              string          `yaml:"az,omitempty"`
	AZs             []string        `yaml:"azs,omitempty"`
	CloudProperties CloudProperties `yaml:"cloud_properties"`
}

func (s *Subnet) AddDNS(dns string) (err error) {
	s.DNS = append(s.DNS, dns)
	return
}

func (s *Subnet) AddReserved(rsv string) (err error) {
	s.Reserved = append(s.Reserved, rsv)
	return
}

type ResourcePool struct {
	Name            string                 `yaml:"name"`
	Network         string                 `yaml:"network"`
	Size            int                    `yaml:"size,omitempty"`
	Stemcell        Stemcell               `yaml:"stemcell"`
	CloudProperties CloudProperties        `yaml:"cloud_properties"`
	Env             map[string]interface{} `yaml:"env,omitempty"`
}

type Stemcell struct {
	Alias   string `yaml:"alias,omitempty"`
	OS      string `yaml:"os,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
	URL     string `yaml:"url,omitempty"`
	SHA1    string `yaml:"sha1,omitempty"`
}

type DiskType DiskPool
type DiskPool struct {
	Name            string          `yaml:"name"`
	DiskSize        int             `yaml:"disk_size"`
	CloudProperties CloudProperties `yaml:"cloud_properties"`
}

type Compilation struct {
	Workers             int             `yaml:"workers"`
	Network             string          `yaml:"network"`
	ReuseCompilationVMs bool            `yaml:"reuse_compilation_vms"`
	CloudProperties     CloudProperties `yaml:"cloud_properties,omitempty"`
	VMType              string          `yaml:"vm_type,omitempty"`
	AZ                  string          `yaml:"az,omitempty"`
}

type Update struct {
	Canaries        int    `yaml:"canaries"`
	MaxInFlight     int    `yaml:"max_in_flight"`
	CanaryWatchTime string `yaml:"canary_watch_time"`
	UpdateWatchTime string `yaml:"update_watch_time"`
	Serial          bool   `yaml:"serial,omitempty"`
}

type Job struct {
	Name           string                 `yaml:"name"`
	Templates      []Template             `yaml:"templates,flow"`
	Lifecycle      string                 `yaml:"lifeycle,omitempty"`
	PersistentDisk string                 `yaml:"persistent_disk,omitempty"`
	Properties     Properties             `yaml:"properties,omitempty"`
	ResourcePool   string                 `yaml:"resource_pool"`
	Update         map[string]interface{} `yaml:"update,omitempty"`
	Instances      int                    `yaml:"instances"`
	Networks       []Network              `yaml:"networks"`
}

type Network struct {
	Name      string        `yaml:"name"`
	StaticIPs []string      `yaml:"static_ips,omitempty"`
	Default   []interface{} `yaml:"default,omitempty"`
}

type Template struct {
	Name    string `yaml:"name"`
	Release string `yaml:"release"`
}

type CloudProvider struct {
	Template   Template                `yaml:"template,flow"`
	MBus       string                  `yaml:"mbus"`
	Properties CloudProviderProperties `yaml:"properties"`
}
type CloudProviderProperties interface{}
type Properties map[string]interface{}
type CloudProperties interface{}
