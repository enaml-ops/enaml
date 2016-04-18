package generators

type JobManifest struct {
	Properties map[string]JobManifestProperty `yaml:"properties"`
}

type JobManifestProperty struct {
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
}

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
	CacheDir  string
	OutputDir string
}
