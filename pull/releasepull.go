package pull

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/mitchellh/ioprogress"
)

func NewRelease(cache string) *Release {
	return &Release{CacheDir: cache}
}

type Release struct {
	CacheDir string
}

func (s *Release) Pull(url string) (filename string) {

	name := path.Base(url)
	filename = s.CacheDir + "/" + name

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
