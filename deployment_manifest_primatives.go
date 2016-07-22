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
	Update             Update        `yaml:"update,omitempty"`
	Lifecycle          string        `yaml:"lifecycle,omitempty"`
}

type InstanceJob struct {
	Name       string      `yaml:"name"`
	Release    string      `yaml:"release"`
	Properties interface{} `yaml:"properties"`
}

type DeploymentNetwork interface {
}

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

func (s DynamicNetwork) GetName() (name string) {
	return s.Name
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

func NewSubnet(cidr, gateway, azName string) Subnet {
	return Subnet{
		Range:   cidr,
		Gateway: gateway,
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

func (s *ResourcePool) SetStemcell(sc Stemcell) (err error) {
	s.Stemcell = sc
	return
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
	Canaries        int    `yaml:"canaries,omitempty"`
	MaxInFlight     int    `yaml:"max_in_flight,omitempty"`
	CanaryWatchTime string `yaml:"canary_watch_time,omitempty"`
	UpdateWatchTime string `yaml:"update_watch_time,omitempty"`
	Serial          bool   `yaml:"serial,omitempty"`
}

type Job struct {
	Name               string     `yaml:"name"`
	Templates          []Template `yaml:"templates,flow"`
	Lifecycle          string     `yaml:"lifeycle,omitempty"`
	PersistentDisk     string     `yaml:"persistent_disk,omitempty"`
	PersistentDiskPool string     `yaml:"persistent_disk_pool,omitempty"`
	Properties         Properties `yaml:"properties,omitempty"`
	ResourcePool       string     `yaml:"resource_pool"`
	Update             Update     `yaml:"update,omitempty"`
	Instances          int        `yaml:"instances"`
	Networks           []Network  `yaml:"networks"`
}

func (s *Job) AddProperty(propName string, prop interface{}) {
	s.Properties[propName] = prop
}

func (s *Job) AddTemplate(t Template) (err error) {
	s.Templates = append(s.Templates, t)
	return
}

func (s *Job) AddNetwork(n Network) (err error) {
	s.Networks = append(s.Networks, n)
	return
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
	SSHTunnel  SSHTunnel               `yaml:"ssh_tunnel"`
}

type SSHTunnel struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	PrivateKeyPath string `yaml:"private_key"`
}
type CloudProviderProperties interface{}
type Properties map[string]interface{}
type CloudProperties interface{}
