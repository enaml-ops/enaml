package generators

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/xchapter7x/enaml/pull"
)

func GenerateReleaseJobsPackage(releaseURL string, cacheDir string, outputDir string) (err error) {
	gen := &ReleaseJobsGenerator{
		OutputDir: outputDir,
	}
	release := pull.NewRelease(cacheDir)
	var filename string
	filename, err = release.Pull(releaseURL)
	if err != nil {
		err = fmt.Errorf("An error occurred downloading %s. %s", releaseURL, err.Error())
		return
	}
	gen.ProcessFile(filename)
	return
}

func (s *ReleaseJobsGenerator) ProcessFile(srcFile string) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	tarReader := s.getTarballReader(f)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if strings.HasPrefix(name, "./jobs/") {
				jobTarball := s.getTarballReader(tarReader)
				jobManifest := s.getJobManifestFromTarball(jobTarball)
				s.processJobManifest(jobManifest, name)
			}
		}
	}
}

func (s *ReleaseJobsGenerator) getJobManifestFromTarball(jobTarball *tar.Reader) (res *tar.Reader) {
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

func (s *ReleaseJobsGenerator) processJobManifest(jobTarball io.Reader, tarballFilename string) {
	jobname := strings.Split(path.Base(tarballFilename), ".")[0]
	buf := new(bytes.Buffer)
	buf.ReadFrom(jobTarball)
	generate(jobname, buf.Bytes(), path.Join(s.OutputDir, jobname))
}

func (s *ReleaseJobsGenerator) getTarballReader(reader io.Reader) *tar.Reader {
	gzf, err := gzip.NewReader(reader)

	if err != nil {
		fmt.Println(err)
	}
	return tar.NewReader(gzf)
}
