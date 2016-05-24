package enamlbosh

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

func (s *Client) GetCloudConfig(httpClient HttpClientDoer) (cloudconfig *enaml.CloudConfigManifest, err error) {
	var req *http.Request
	var res *http.Response
	var resBody = make([]CloudConfigResponseBody, 1)

	if req, err = http.NewRequest("GET", s.buildBoshURL("/cloud_configs?limit=1"), nil); err == nil {
		req.SetBasicAuth(s.user, s.pass)
		req.Header.Set("content-type", "text/yaml")

		if res, err = httpClient.Do(req); err == nil {
			var b []byte
			b, err = ioutil.ReadAll(res.Body)
			json.Unmarshal(b, &resBody)
			cloudconfig = enaml.NewCloudConfigManifest([]byte(resBody[0].Properties))
		}
	}
	return
}

func (s *Client) buildBoshURL(urlpath string) (boshurl string) {
	boshurl = s.host + ":" + strconv.Itoa(s.port) + urlpath
	return
}
