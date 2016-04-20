package pull

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/mitchellh/ioprogress"
)

// NewRelease creates a new Release instance
func NewRelease(cache string) *Release {
	return &Release{CacheDir: cache}
}

// Release is a BOSH release with a configurable cache dir
type Release struct {
	CacheDir string
}

// Pull downloads the specified Release to the local cache dir
func (r *Release) Pull(url string) (filename string, err error) {
	name := path.Base(url)
	filename = r.CacheDir + "/" + name

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Could not find release in local cache. Downloading now.")
		err = r.download(url, filename)
	}
	return
}

func (r *Release) download(url, local string) (err error) {
	var out *os.File
	out, err = os.Create(local)
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	progressR := &ioprogress.Reader{
		Reader: resp.Body,
		Size:   resp.ContentLength,
	}
	_, err = io.Copy(out, progressR)
	return
}
