package generators

import "github.com/enaml-ops/enaml"

type ReleaseJobsGenerator struct {
	OutputDir string
}

type jobStructTemplate struct {
	PackageName string
	JobName     string
	Elements    []elementStruct
}

type elementStruct struct {
	ElementName     string
	ElementType     string
	ElementYamlName string
	ElementComments string
}

type ObjectField struct {
	ElementName       string
	ElementType       string
	ElementAnnotation string
	Meta              enaml.JobManifestProperty
}

type record struct {
	Length int
	Slice  []string
	Orig   string
	Yaml   enaml.JobManifestProperty
}

type processing struct {
	max  int
	recs []record
}
