package generators

type elementStruct struct {
	ElementName     string
	ElementType     string
	ElementYamlName string
}

type jobStructTemplate struct {
	JobName  string
	Elements []elementStruct
}

type ReleaseJobsGenerator struct {
	OutputDir string
}
