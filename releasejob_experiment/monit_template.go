package releasejob_experiment

import (
	"io"

	"github.com/alecthomas/template"
)

type monit struct {
	ProcessName string
	PID         string
}

func writeMonitFile(f io.Writer, processname string, pidpath string) error {
	monitTemplate := `check process {{.ProcessName}}
	with pidfile {{.PID}} 
	start program "/var/vcap/jobs/{{.ProcessName}}/bin/{{.ProcessName}}_ctl start" with timeout 120 seconds
	stop program "/var/vcap/jobs/{{.ProcessName}}/bin/{{.ProcessName}}_ctl stop" with timeout 30 seconds
	group vcap
`
	m := monit{
		ProcessName: processname,
		PID:         pidpath,
	}
	tmpl, err := template.New("monit-file-create").Parse(monitTemplate)

	if err != nil {
		return err
	}
	err = tmpl.Execute(f, m)

	if err != nil {
		return err
	}
	return nil
}
