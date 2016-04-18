package generators

const (
	structTemplate = `package releasejobs
	type {{.JobName}} struct {
		{{ range $key, $value := .Elements }}
		{{ $value.ElementName }} {{ $value.ElementType }} ` + "`" + `yaml:"{{$value.ElementYamlName}},omitempty"` + "`" + `
		{{ end }}
	}
	`
)
