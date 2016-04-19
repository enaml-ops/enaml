package enaml

func NewInstanceGroup() *InstanceGroup {
	return &InstanceGroup{}
}

func (s *InstanceGroup) AddJob(job InstanceJob) (err error) {
	s.Jobs = append(s.Jobs, job)
	return
}
