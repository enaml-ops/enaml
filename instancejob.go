package enaml

func NewInstanceJob(name string, release string, properties interface{}) *InstanceJob {
	return &InstanceJob{
		Name:       name,
		Release:    release,
		Properties: properties,
	}
}
