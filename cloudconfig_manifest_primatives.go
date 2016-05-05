package enaml

type CloudConfigManifest struct {
	AZs         []AZ                `yaml:"azs,omitempty"`
	VMTypes     []VMType            `yaml:"vm_types,omitempty"`
	DiskTypes   []DiskType          `yaml:"disk_types,omitempty"`
	Networks    []DeploymentNetwork `yaml:"networks,omitempty"`
	Compilation *Compilation        `yaml:"compilation,omitempty"`
}

func (s *CloudConfigManifest) GetManifest() CloudConfigManifest {
	return *s
}

//ContainsAZName -
func (s *CloudConfigManifest) ContainsAZName(azName string) (result bool) {
	result = false
	for _, az := range s.AZs {
		if az.Name == azName {
			result = true
			return
		}
	}
	return
}

//ContainsVMType -
func (s *CloudConfigManifest) ContainsVMType(vmTypeName string) (result bool) {
	result = false
	for _, vmType := range s.VMTypes {
		if vmType.Name == vmTypeName {
			result = true
			return
		}
	}
	return
}

//ContainsDiskType -
func (s *CloudConfigManifest) ContainsDiskType(diskTypeName string) (result bool) {
	result = false
	for _, diskType := range s.DiskTypes {
		if diskType.Name == diskTypeName {
			result = true
			return
		}
	}
	return
}

type VMType struct {
	Name            string      `yaml:"name,omitempty"`
	CloudProperties interface{} `yaml:"cloud_properties,omitempty"`
}

type AZ struct {
	Name            string      `yaml:"name,omitempty"`
	CloudProperties interface{} `yaml:"cloud_properties,omitempty"`
}

func (s *CloudConfigManifest) AddAZ(az AZ) (err error) {
	s.AZs = append(s.AZs, az)
	return
}

func (s *CloudConfigManifest) AddVMType(vmt VMType) (err error) {
	s.VMTypes = append(s.VMTypes, vmt)
	return
}

func (s *CloudConfigManifest) AddDiskType(dskt DiskType) (err error) {
	s.DiskTypes = append(s.DiskTypes, dskt)
	return
}

func (s *CloudConfigManifest) AddNetwork(ntw DeploymentNetwork) (err error) {
	s.Networks = append(s.Networks, ntw)
	return
}

func (s *CloudConfigManifest) SetCompilation(cpl *Compilation) (err error) {
	s.Compilation = cpl
	return
}
