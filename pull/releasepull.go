package pull

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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
func (r *Release) Pull(releaseLocation string) (filename string, err error) {
	u, uerr := url.Parse(releaseLocation)
	if uerr != nil || !(u.Scheme == "http" || u.Scheme == "https") {
		// assume a local file, ensure it exists
		if _, ferr := os.Stat(releaseLocation); os.IsNotExist(ferr) {
			err = fmt.Errorf("Could not pull %s. The file doesn't exist or isn't a valid http(s) URL", releaseLocation)
			return
		}
		filename = releaseLocation
	} else {
		// remote file, ensure its in the local cache
		filename = r.CacheDir + "/" + path.Base(releaseLocation)
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			fmt.Println("Could not find release in local cache. Downloading now.")
			err = r.download(releaseLocation, filename)
		}
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
