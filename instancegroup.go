package enaml

func NewInstanceGroup(name string, instanceCount int, vmName string, stemCellAlias string) *InstanceGroup {
	return &InstanceGroup{
		Name:      name,
		Instances: instanceCount,
		VMType:    vmName,
		Stemcell:  stemCellAlias,
	}
}

func (s *InstanceGroup) AddAZ(az string) (err error) {
	s.AZs = append(s.AZs, az)
	return
}

func (s *InstanceGroup) AddNetwork(network Network) (err error) {
	s.Networks = append(s.Networks, network)
	return
}

func (s *InstanceGroup) GetJobByName(name string) (j *InstanceJob) {

	for _, job := range s.Jobs {
		if job.Name == name {
			j = &job
			break
		}
	}
	return
}

func (s *InstanceGroup) GetNetworkByName(name string) (n *Network) {

	for _, network := range s.Networks {
		if network.Name == name {
			n = &network
			break
		}
	}
	return
}

func (s *InstanceGroup) AddJob(job *InstanceJob) (err error) {
	s.Jobs = append(s.Jobs, *job)
	return
}
