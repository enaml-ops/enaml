package enamlbosh

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/enaml-ops/enaml"
)

func (s *Client) NewCloudConfigRequest(cloudconfig enaml.CloudConfigManifest) (req *http.Request, err error) {
	var b []byte
	var body io.Reader
	if b, err = cloudconfig.Bytes(); err == nil {
		body = bytes.NewReader(b)
		req, err = http.NewRequest("POST", s.buildBoshURL("/cloud_configs"), body)
		req.SetBasicAuth(s.user, s.pass)
		req.Header.Set("content-type", "text/yaml")
	}
	return
}

func (s *Client) buildBoshURL(urlpath string) (boshurl string) {
	boshurl = s.host + ":" + strconv.Itoa(s.port) + urlpath
	return
}
