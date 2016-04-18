package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/ioprogress"
)

func init() {
	if c := os.Getenv("CACHE_DIR"); c != "" {
		cacheDir = c
	}
	os.MkdirAll(cacheDir, 0755)
}

const (
	CacheDir = ".cache"
)

var (
	cacheDir = CacheDir
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "generate-jobs",
			Aliases: []string{"gj"},
			Usage:   "generate golang structs for the jobs in a given release",
			Action: func(c *cli.Context) {
				GenerateReleaseJobsPackage(c.Args().First())
				println("completed generating release job structs for ", c.Args().First())
			},
		},
		{
			Name:    "diff-release",
			Aliases: []string{"dr"},
			Usage:   "show a diff between 2 releases given",
			Action: func(c *cli.Context) {
				println("unimplemented")
				println("release job properties diff", c.Args().First())
			},
		},
		{
			Name:    "diff-job",
			Aliases: []string{"dj"},
			Usage:   "show diff between jobs across 2 releases",
			Action: func(c *cli.Context) {
				println("unimplemented")
				println("release job properties diff", c.Args().First())
			},
		},
	}
	app.Run(os.Args)
}

func GenerateReleaseJobsPackage(releaseURL string) (err error) {
	filename := downloadFromURL(releaseURL)
	processFile(filename)
	return
}

func downloadFromURL(url string) (filename string) {
	name := path.Base(url)
	filename = cacheDir + "/" + name

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Could not find release in local cache. Downloading now.")
		out, _ := os.Create(filename)
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		progressR := &ioprogress.Reader{
			Reader: resp.Body,
			Size:   resp.ContentLength,
		}
		io.Copy(out, progressR)
	}
	return
}

func processFile(srcFile string) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	tarReader := getTarballReader(f)
	i := 0

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if strings.HasPrefix(name, "./jobs/") {
				jobTarball := getTarballReader(tarReader)
				jobManifest := getJobManifestFromTarball(jobTarball)
				processJobManifest(jobManifest, name)
			}
		}
		i++
	}
}

func getJobManifestFromTarball(jobTarball *tar.Reader) (res *tar.Reader) {
	var jobManifestFilename = "./job.MF"

	for {
		header, _ := jobTarball.Next()
		if header.Name == jobManifestFilename {
			res = jobTarball
			break
		}
	}
	return
}

func processJobManifest(jobTarball io.Reader, tarballFilename string) {
	var elements []elementStruct
	var defaultElementType = "interface{}"
	buf := new(bytes.Buffer)
	buf.ReadFrom(jobTarball)
	manifestYaml := JobManifest{}
	yaml.Unmarshal(buf.Bytes(), &manifestYaml)

	for k, v := range manifestYaml.Properties {
		myType := defaultElementType

		if v.Default != nil {
			myType = fmt.Sprint(reflect.ValueOf(v.Default).Type())
		}
		elements = append(elements, elementStruct{
			ElementName:     parseElementName(k),
			ElementType:     myType,
			ElementYamlName: k,
		})
	}
	jobName := strings.Split(path.Base(tarballFilename), ".")[0]
	jobName = strings.ToUpper(jobName[:1]) + jobName[1:]
	job := jobStructTemplate{
		JobName:  jobName,
		Elements: elements,
	}
	tmpl, err := template.New("job").Parse(structTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, job)
	if err != nil {
		panic(err)
	}
}

func parseElementName(name string) string {
	f := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '.'
	})
	for i := range f {
		f[i] = strings.ToUpper(f[i][:1]) + f[i][1:]
	}
	return strings.Join(f, "")
}

func getTarballReader(reader io.Reader) *tar.Reader {
	gzf, err := gzip.NewReader(reader)

	if err != nil {
		fmt.Println(err)
	}
	return tar.NewReader(gzf)
}

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

const (
	structTemplate = `package releasejobs
	type {{.JobName}} struct {
		{{ range $key, $value := .Elements }}
		{{ $value.ElementName }} {{ $value.ElementType }} ` + "`" + `yaml:"{{$value.ElementYamlName}},omitempty"` + "`" + `
		{{ end }}
	}
	`
)
