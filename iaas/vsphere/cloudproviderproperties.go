package vsphere

func NewCloudProviderProperties(vcenter *VCenter) *CloudProviderProperties {
	return &CloudProviderProperties{
		VCenter: vcenter,
	}
}

func (s *CloudProviderProperties) AddAgent(key string, value string) (err error) {
	if s.Agent == nil {
		s.Agent = make(Agent)
	}
	s.Agent[key] = value
	return
}

func (s *CloudProviderProperties) AddBlobstore(key string, value string) (err error) {
	if s.Blobstore == nil {
		s.Blobstore = make(Blobstore)
	}
	s.Blobstore[key] = value
	return
}

func (s *CloudProviderProperties) AddNTP(ntp string) (err error) {
	s.NTP = append(s.NTP, ntp)
	return
}
