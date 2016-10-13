/*
*
*
*type Inventory struct {
	Material string
	Count    uint
}
sweaters := Inventory{"wool", 17}
tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
if err != nil { panic(err) }
err = tmpl.Execute(os.Stdout, sweaters)
if err != nil { panic(err) }
*
*
*
*/

package releasejob_experiment

import (
	"io"

	"text/template"
)

type spec struct {
	ProcessName string
	Properties  []JobProperty
	Packages    []string
}

func writeSpecFile(f io.Writer, processname string, properties []JobProperty, packages []string) error {
	specTemplate := `---
name: {{.ProcessName}} 

templates:
  {{.ProcessName}}_ctl.erb: bin/{{.ProcessName}}_ctl

packages:{{ range $key, $value := .Packages }}
- {{ $value }}{{ end }}

properties:{{ range $key, $value := .Properties }}
	{{ $value.Name }}:
		description:{{ $value.Description }}
		{{if $value.Default}}default: {{$value.Default}}{{ end }}{{end}}
`
	s := spec{
		ProcessName: processname,
		Properties:  properties,
		Packages:    packages,
	}
	tmpl, err := template.New("monit-file-create").Parse(specTemplate)

	if err != nil {
		return err
	}
	err = tmpl.Execute(f, s)

	if err != nil {
		return err
	}
	return nil
}
