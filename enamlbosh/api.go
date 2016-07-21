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

// newRequest is like http.NewRequest, with the exception that it will add
// basic auth headers if the client is configured for basic auth.
func (s *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err == nil {
		req.SetBasicAuth(s.user, s.pass)
	}
	return req, err
}

func (s *Client) NewCloudConfigRequest(cloudconfig enaml.CloudConfigManifest) (*http.Request, error) {
	b, err := cloudconfig.Bytes()
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(b)
	req, err := s.newRequest("POST", s.buildBoshURL("/cloud_configs"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "text/yaml")
	return req, nil
}

func (s *Client) GetTask(taskID int, httpClient HttpClientDoer) (BoshTask, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/tasks/"+strconv.Itoa(taskID)), nil)
	if err != nil {
		return BoshTask{}, err
	}
	req.Header.Set("content-type", "text/yaml")
	res, err := httpClient.Do(req)
	if err != nil {
		return BoshTask{}, err
	}
	defer res.Body.Close()
	lo.G.Debug("task request complete")
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return BoshTask{}, err
	}
	lo.G.Debug("rest resp: ", string(b))
	var bt BoshTask
	json.Unmarshal(b, &bt)

	if bt.ID != taskID {
		return bt, fmt.Errorf("could not find the given task: %v", taskID)
	}
	return bt, nil
}

func (s *Client) PostRemoteRelease(rls enaml.Release, httpClient HttpClientDoer) (BoshTask, error) {
	if rls.URL == "" || rls.SHA1 == "" {
		return BoshTask{}, fmt.Errorf("url or sha not set. these are required for remote stemcells URL: %s , SHA: %s", rls.URL, rls.SHA1)
	}
	reqMap := map[string]string{
		"location": rls.URL,
		"sha1":     rls.SHA1,
	}
	reqBytes, _ := json.Marshal(reqMap)
	reqBody := bytes.NewReader(reqBytes)

	req, err := s.newRequest("POST", s.buildBoshURL("/releases"), reqBody)
	if err != nil {
		return BoshTask{}, err
	}

	req.Header.Set("content-type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return BoshTask{}, err
	}

	defer res.Body.Close()
	lo.G.Debug("release request complete")

	var bt BoshTask
	err = json.NewDecoder(res.Body).Decode(&bt)
	return bt, err
}

func (s *Client) GetStemcells(httpClient HttpClientDoer) (stemcells []DeployedStemcell, err error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/stemcells"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var stemcells []DeployedStemcell
	err = json.NewDecoder(res.Body).Decode(&stemcells)
	return stemcells, err
}

func (s *Client) CheckRemoteStemcell(sc enaml.Stemcell, httpClient HttpClientDoer) (exists bool, err error) {
	if (sc.Name == "" && sc.OS == "") || sc.Version == "" {
		return false, fmt.Errorf("name or version not set. these are required to check for remote stemcells Name: %s , Version: %s", sc.Name, sc.Version)
	}

	stemcells, err := s.GetStemcells(httpClient)
	if err != nil {
		return false, err
	}

	for _, stemcell := range stemcells {
		if (stemcell.Name == sc.Name || stemcell.OS == sc.OS) && stemcell.Version == sc.Version {
			return true, nil
		}
	}
	return false, nil
}

func (s *Client) PostRemoteStemcell(sc enaml.Stemcell, httpClient HttpClientDoer) (BoshTask, error) {
	if sc.URL == "" || sc.SHA1 == "" {
		return BoshTask{}, fmt.Errorf("url or sha not set. these are required for remote stemcells URL: %s , SHA: %s", sc.URL, sc.SHA1)
	}
	reqMap := map[string]string{
		"location": sc.URL,
		"sha1":     sc.SHA1,
	}
	reqBytes, _ := json.Marshal(reqMap)
	reqBody := bytes.NewReader(reqBytes)

	req, err := s.newRequest("POST", s.buildBoshURL("/stemcells"), reqBody)
	if err != nil {
		return BoshTask{}, err
	}

	req.Header.Set("content-type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return BoshTask{}, err
	}

	defer res.Body.Close()
	lo.G.Debug("stemcell request complete")

	var bt BoshTask
	err = json.NewDecoder(res.Body).Decode(&bt)
	return bt, err
}

func (s *Client) PostDeployment(deploymentManifest enaml.DeploymentManifest, httpClient HttpClientDoer) (BoshTask, error) {
	reqBody := bytes.NewReader(deploymentManifest.Bytes())
	req, err := s.newRequest("POST", s.buildBoshURL("/deployments"), reqBody)
	if err != nil {
		return BoshTask{}, err
	}

	req.Header.Set("content-type", "text/yaml")
	res, err := httpClient.Do(req)
	if err != nil {
		return BoshTask{}, err
	}

	defer res.Body.Close()
	lo.G.Debug("deployment request complete")
	var bt BoshTask
	err = json.NewDecoder(res.Body).Decode(&bt)
	return bt, err
}

func (s *Client) GetCloudConfig(httpClient HttpClientDoer) (*enaml.CloudConfigManifest, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/cloud_configs?limit=1"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "text/yaml")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var cc []CloudConfigResponseBody
	err = json.NewDecoder(res.Body).Decode(&cc)
	if err != nil {
		return nil, err
	}
	return enaml.NewCloudConfigManifest([]byte(cc[0].Properties)), nil
}

func (s *Client) GetInfo(httpClient HttpClientDoer) (*BoshInfo, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/info"), nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var bi BoshInfo
	err = json.NewDecoder(res.Body).Decode(&bi)
	if err != nil {
		return nil, err
	}
	return &bi, nil
}

func (s *Client) buildBoshURL(urlpath string) string {
	return s.host + ":" + strconv.Itoa(s.port) + urlpath
}
