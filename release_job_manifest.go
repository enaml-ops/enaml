package enaml

type JobManifest struct {
	Properties map[string]JobManifestProperty `yaml:"properties"`
}

type JobManifestProperty struct {
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
}
