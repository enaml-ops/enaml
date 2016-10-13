package releasejob_experiment

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

//BoshJobMeta - release job meta data which defines
//the job its dependencies and functionality
type BoshJobMeta struct {
	Version       string
	Name          string
	Packages      []string
	JobProperties []JobProperty
	PIDFile       string
}

//JobProperty - an object defining a job property
type JobProperty struct {
	Description string
	Name        string
	EnvVar      string
	Default     string
}

//BoshJob - interface defining what a bosh release job should do
type BoshJob interface {
	Meta() BoshJobMeta
	Start() error
	Stop() error
}

//JobRunner - an object which can run a given bosh job
type JobRunner struct {
	Args []string
}

//BuildJob - function to create a compiled version of the current job
func BuildJob(jobMeta BoshJobMeta, dest string) error {
	b, err := json.Marshal(jobMeta)

	if err != nil {
		return err
	}
	fmt.Println("building job: ", string(b))
	monitFile, specFile, err := createJobFiles(dest, jobMeta.Name)

	if err != nil {
		return err
	}
	defer monitFile.Close()
	defer specFile.Close()
	err = writeMonitFile(monitFile, jobMeta.Name, jobMeta.PIDFile)

	if err != nil {
		return err
	}
	err = writeSpecFile(specFile, jobMeta.Name)
	return err
}

func createJobFiles(dest, name string) (monitFile *os.File, specFile *os.File, err error) {
	base := path.Join(dest, name, "templates")
	fmt.Println("creating: ", base)
	err = os.MkdirAll(base, 0777)

	if err != nil {
		return nil, nil, err
	}
	monitFile, err = os.Create(path.Join(dest, name, "monit"))

	if err != nil {
		return nil, nil, err
	}
	specFile, err = os.Create(path.Join(dest, name, "spec"))

	if err != nil {
		return nil, nil, err
	}
	err = copyFile(os.Args[0], path.Join(base, name+"_ctl"))
	return monitFile, specFile, err
}

//NewJobRunner - constructs and returns a JobRunner type object
func NewJobRunner(args []string) *JobRunner {
	return &JobRunner{
		Args: args,
	}
}

func (s *JobRunner) Run(job BoshJob) error {
	_, err := os.Stat(job.Meta().PIDFile)

	if s.Args[1] == "build" && len(s.Args) == 3 {
		return BuildJob(job.Meta(), s.Args[2])

	} else if s.Args[1] == "build" && len(s.Args) != 3 {
		return fmt.Errorf("build command uses format: `./job build <output-dir>`")
	}

	if s.cleanStart(err) {
		return job.Start()
	}

	if s.cleanStop(err) {
		return job.Stop()
	}

	if !s.cleanStart(err) {
		return fmt.Errorf("sorry PID file %v already exists.", job.Meta().PIDFile)
	}

	if !s.cleanStop(err) {
		return fmt.Errorf("PID file %v did not exist for running process.", job.Meta().PIDFile)
	}
	return fmt.Errorf("sorry we only take 'start|stop|build' as args")
}

func (s *JobRunner) cleanStart(err error) bool {
	return s.Args[1] == "start" && os.IsNotExist(err)
}

func (s *JobRunner) cleanStop(err error) bool {
	return s.Args[1] == "stop" && os.IsExist(err)
}

func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
