package diff

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/kr/pretty"
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/pull"
)

func NewDiff(cacheDir string) *Diff {
	return &Diff{CacheDir: cacheDir}
}

type Diff struct {
	CacheDir string
}

func (s *Diff) ReleaseDiff(releaseURLA, releaseURLB string) (diffset []string, err error) {
	release := pull.NewRelease(s.CacheDir)
	filenameA, _ := release.Pull(releaseURLA)
	filenameB, _ := release.Pull(releaseURLB)
	GetReleaseManifest(filenameA)
	GetReleaseManifest(filenameB)
	return
}

func (s *Diff) JobDiffBetweenReleases(jobname, releaseURLA, releaseURLB string) (diffset []string, err error) {
	var (
		jobA      *tar.Reader
		jobB      *tar.Reader
		filenameA string
		filenameB string
		ok        bool
	)
	release := pull.NewRelease(s.CacheDir)
	filenameA, err = release.Pull(releaseURLA)
	if err != nil {
		err = fmt.Errorf("An error occurred downloading %s. %s", releaseURLA, err.Error())
		return
	}
	filenameB, err = release.Pull(releaseURLB)
	if err != nil {
		err = fmt.Errorf("An error occurred downloading %s. %s", releaseURLB, err.Error())
		return
	}
	jobA, ok = ProcessReleaseArchive(filenameA)[jobname]

	if !ok {
		err = errors.New("could not find jobname in release A")
	}
	jobB, ok = ProcessReleaseArchive(filenameB)[jobname]

	if !ok {
		err = errors.New("could not find jobname in release B")
	}
	bufA := new(bytes.Buffer)
	bufA.ReadFrom(jobA)
	bufB := new(bytes.Buffer)
	bufB.ReadFrom(jobB)
	diffset = JobPropertiesDiff(bufA.Bytes(), bufB.Bytes())
	return
}

func GetReleaseManifest(srcFile string) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	tarReader := getTarballReader(f)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if path.Base(name) == "release.MF" {
				fmt.Println("found the release manifest")
			}
		}
	}
}

func ProcessReleaseArchive(srcFile string) (jobs map[string]*tar.Reader) {
	jobs = make(map[string]*tar.Reader)
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	tarReader := getTarballReader(f)

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
				jobName := strings.Split(path.Base(name), ".")[0]
				jobs[jobName] = jobManifest
			}
		}
	}
	return
}

func getTarballReader(reader io.Reader) *tar.Reader {
	gzf, err := gzip.NewReader(reader)

	if err != nil {
		fmt.Println(err)
	}
	return tar.NewReader(gzf)
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

func JobPropertiesDiff(a, b []byte) []string {
	var objA enaml.JobManifest
	var objB enaml.JobManifest
	yaml.Unmarshal(a, &objA)
	yaml.Unmarshal(b, &objB)
	mp := pretty.Diff(objA, objB)
	return mp
}
