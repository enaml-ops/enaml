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

	"github.com/alecthomas/template"
)

type spec struct {
	ProcessName string
}

func writeSpecFile(f io.Writer, processname string) error {
	specTemplate := `---
name: {{.ProcessName}} 

templates:
  {{.ProcessName}}_ctl.erb: bin/{{.ProcessName}}_ctl

packages:
- gemfire
- jdk8
- gemfire-utils
- jq

properties:
  external_dependencies.router.system_domain:
    description: "System domain"
  gemfire.locator.addresses:
    description: "List of GemFire Locator addresses of the form X.X.X.X"
  gemfire.locator.port:
    description: "Port the Locator will listen on"
    default: "55221"
  gemfire.locator.rest_port:
    description: "Port the Locator will listen on for REST API"
    default: "8080"
  gemfire.locator.vm_memory:
    description: "RAM allocated to the locator VM in MB"
  gemfire.cluster_topology.number_of_locators:
    description: "Current topology"
    default: "2"
  gemfire.cluster_topology.min_number_of_locators:
    description: "min number of locators which should be present"
    default: "2"
`
	s := spec{
		ProcessName: processname,
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
