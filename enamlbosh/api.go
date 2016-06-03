package enamlbosh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/enaml-ops/enaml"
	"github.com/xchapter7x/lo"
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

func (s *Client) GetTask(taskID int, httpClient HttpClientDoer) (bt BoshTask, err error) {
	var req *http.Request
	var res *http.Response

	if req, err = http.NewRequest("GET", s.buildBoshURL("/tasks/"+strconv.Itoa(taskID)), nil); err == nil {
		req.SetBasicAuth(s.user, s.pass)
		req.Header.Set("content-type", "text/yaml")

		if res, err = httpClient.Do(req); err == nil {
			lo.G.Debug("task request complete")
			var b []byte
			b, err = ioutil.ReadAll(res.Body)
			lo.G.Debug("rest resp: ", string(b))
			json.Unmarshal(b, &bt)
		}
	}

	if bt.ID != taskID && err == nil {
		err = fmt.Errorf("could not find the given task: %v", taskID)
	}
	return
}

func (s *Client) PostRemoteRelease(rls enaml.Release, httpClient HttpClientDoer) (bt []BoshTask, err error) {

	if rls.URL == "" || rls.SHA1 == "" {
		err = fmt.Errorf("url or sha not set. these are required for remote stemcells URL: %s , SHA: %s", rls.URL, rls.SHA1)

	} else {
		var req *http.Request
		var res *http.Response
		var reqMap = map[string]string{
			"location": rls.URL,
			"sha1":     rls.SHA1,
		}
		var reqBytes, _ = json.Marshal(reqMap)
		var reqBody = bytes.NewReader(reqBytes)

		if req, err = http.NewRequest("POST", s.buildBoshURL("/releases"), reqBody); err == nil {
			req.SetBasicAuth(s.user, s.pass)
			req.Header.Set("content-type", "application/json")

			if res, err = httpClient.Do(req); err == nil {
				lo.G.Debug("release request complete")
				var b []byte

				if b, err = ioutil.ReadAll(res.Body); err == nil {
					lo.G.Debug("rest resp: ", string(b))
					err = json.Unmarshal(b, &bt)
				}
			}
		}
	}
	return
}

func (s *Client) PostRemoteStemcell(sc enaml.Stemcell, httpClient HttpClientDoer) (bt []BoshTask, err error) {

	if sc.URL == "" || sc.SHA1 == "" {
		err = fmt.Errorf("url or sha not set. these are required for remote stemcells URL: %s , SHA: %s", sc.URL, sc.SHA1)

	} else {
		var req *http.Request
		var res *http.Response
		var reqMap = map[string]string{
			"location": sc.URL,
			"sha1":     sc.SHA1,
		}
		var reqBytes, _ = json.Marshal(reqMap)
		var reqBody = bytes.NewReader(reqBytes)

		if req, err = http.NewRequest("POST", s.buildBoshURL("/stemcells"), reqBody); err == nil {
			req.SetBasicAuth(s.user, s.pass)
			req.Header.Set("content-type", "application/json")

			if res, err = httpClient.Do(req); err == nil {
				lo.G.Debug("stemcell request complete")
				var b []byte

				if b, err = ioutil.ReadAll(res.Body); err == nil {
					lo.G.Debug("rest resp: ", string(b))
					err = json.Unmarshal(b, &bt)
				}
			}
		}
	}
	return
}

func (s *Client) PostDeployment(deploymentManifest enaml.DeploymentManifest, httpClient HttpClientDoer) (boshTask []BoshTask, err error) {
	var req *http.Request
	var res *http.Response
	var reqBody = bytes.NewReader(deploymentManifest.Bytes())

	if req, err = http.NewRequest("POST", s.buildBoshURL("/deployments"), reqBody); err == nil {
		req.SetBasicAuth(s.user, s.pass)
		req.Header.Set("content-type", "text/yaml")

		if res, err = httpClient.Do(req); err == nil {
			lo.G.Debug("deployment request complete")
			var b []byte

			if b, err = ioutil.ReadAll(res.Body); err == nil {
				lo.G.Debug("rest resp: ", string(b))
				err = json.Unmarshal(b, &boshTask)
			}
		}
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
			lo.G.Debug("rest resp: ", string(b))
			json.Unmarshal(b, &resBody)
			cloudconfig = enaml.NewCloudConfigManifest([]byte(resBody[0].Properties))
		}
	}
	return
}

func (s *Client) GetInfo(httpClient HttpClientDoer) (bi *BoshInfo, err error) {
	var req *http.Request
	var res *http.Response
	bi = new(BoshInfo)

	if req, err = http.NewRequest("GET", s.buildBoshURL("/info"), nil); err == nil {
		req.SetBasicAuth(s.user, s.pass)
		req.Header.Set("content-type", "text/yaml")

		if res, err = httpClient.Do(req); err == nil {
			var b []byte
			b, err = ioutil.ReadAll(res.Body)
			lo.G.Debug("rest resp: ", string(b))
			json.Unmarshal(b, bi)
		}
	}
	return

}

func (s *Client) buildBoshURL(urlpath string) (boshurl string) {
	boshurl = s.host + ":" + strconv.Itoa(s.port) + urlpath
	return
}
