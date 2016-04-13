package vsphere

type CloudProviderProperties struct {
	VCenter   *VCenter
	Agent     Agent     `yaml:",flow"`
	Blobstore Blobstore `yaml:",flow"`
	NTP       NTP       `yaml:",flow"`
}

type Agent map[string]string
type Blobstore map[string]string
type NTP []string

type VCenter struct {
	Address     string
	User        string
	Password    string
	DataCenters []DataCenter
}

type DataCenter struct {
	Name                       string
	VMFolder                   string   `yaml:"vm_folder"`
	TemplateFolder             string   `yaml:"template_folder"`
	DatastorePattern           string   `yaml:"datastore_pattern"`
	PersistentDatastorePattern string   `yaml:"persistent_datastore_pattern"`
	DiskPath                   string   `yaml:"disk_path"`
	Clusters                   []string `yaml:",flow"`
}
