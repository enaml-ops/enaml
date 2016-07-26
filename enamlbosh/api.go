package enamlbosh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/enaml-ops/enaml"
	"github.com/op/go-logging"
	"github.com/xchapter7x/lo"
)

// NewClient creates a new bosh client.  It queries bosh's /info endpoint
// to determine if user authentication should be done via UAA or basic auth.
func NewClient(user, pass, host string, port int, sslIgnore bool) (*Client, error) {
	c := &Client{
		user: user,
		pass: pass,
		host: host,
		port: port,
		http: &http.Client{Transport: transport(sslIgnore)},
	}
	c.http.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req.URL, _ = url.Parse(req.URL.Scheme + "://" + via[0].URL.Host + req.URL.Path)
		setAuth(c, req)
		return nil
	}

	info, err := c.GetInfo()
	if err != nil {
		return nil, err
	}

	authType, ok := info.UserAuthentication["type"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected user auth in response: %#v", info)
	}

	switch authType {
	case "basic":
		return c, nil
	case "uaa":
		opt, ok := info.UserAuthentication["options"].(map[string]interface{})
		if !ok {
			return nil, errors.New("missing UAA options in response")
		}
		uaaURL, ok := opt["url"].(string)
		if !ok {
			return nil, errors.New("couln't get UAA URL")
		}
		err = c.getToken(uaaURL + "/oauth/token")
		if err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, errors.New("unknown user auth type: " + authType)
	}
}

func (c *Client) getToken(tokURL string) error {
	cfg := oauth2.Config{
		ClientID:     "bosh_cli",
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			TokenURL: tokURL,
		},
	}
	// make sure we use our HTTP client for getting the token,
	// not http.DefausltClient
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c.http)
	tok, err := cfg.PasswordCredentialsToken(ctx, c.user, c.pass)
	if err != nil {
		return err
	}
	c.token = tok
	return nil
}

func transport(insecureSkipVerify bool) *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
}

// newRequest is like http.NewRequest, with the exception that it will add
// basic auth headers if the client is configured for basic auth.
func (s *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	setAuth(s, req)
	return req, nil
}

func setAuth(c *Client, r *http.Request) {
	if c.token == nil {
		r.SetBasicAuth(c.user, c.pass)
	} else {
		c.token.SetAuthHeader(r)
	}
}

func (s *Client) newCloudConfigRequest(cloudconfig enaml.CloudConfigManifest) (*http.Request, error) {
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

// PushCloudConfig uploads a cloud config to bosh.
func (s *Client) PushCloudConfig(manifest []byte) error {
	ccm := enaml.NewCloudConfigManifest(manifest)
	req, err := s.newCloudConfigRequest(*ccm)
	if err != nil {
		return err
	}
	res, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s error pushing cloud config to BOSH: %s", res.Status, string(body))
	}
	return nil
}

func (s *Client) GetTask(taskID int) (BoshTask, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/tasks/"+strconv.Itoa(taskID)), nil)
	if err != nil {
		return BoshTask{}, err
	}
	req.Header.Set("content-type", "text/yaml")
	res, err := s.http.Do(req)
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

func (s *Client) PostRemoteRelease(rls enaml.Release) (BoshTask, error) {
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
	res, err := s.http.Do(req)
	if err != nil {
		return BoshTask{}, err
	}

	defer res.Body.Close()
	lo.G.Debug("release request complete")

	var bt BoshTask
	err = json.NewDecoder(res.Body).Decode(&bt)
	return bt, err
}

func (s *Client) GetStemcells() ([]DeployedStemcell, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/stemcells"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var stemcells []DeployedStemcell
	err = json.NewDecoder(res.Body).Decode(&stemcells)
	return stemcells, err
}

func (s *Client) CheckRemoteStemcell(sc enaml.Stemcell) (exists bool, err error) {
	if (sc.Name == "" && sc.OS == "") || sc.Version == "" {
		return false, fmt.Errorf("name or version not set. these are required to check for remote stemcells Name: %s , Version: %s", sc.Name, sc.Version)
	}

	stemcells, err := s.GetStemcells()
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

func (s *Client) PostRemoteStemcell(sc enaml.Stemcell) (BoshTask, error) {
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
	res, err := s.http.Do(req)
	if err != nil {
		return BoshTask{}, err
	}

	defer res.Body.Close()
	lo.G.Debug("stemcell request complete")

	var bt BoshTask
	err = json.NewDecoder(res.Body).Decode(&bt)
	return bt, err
}

func (s *Client) PostDeployment(deploymentManifest enaml.DeploymentManifest) (BoshTask, error) {
	reqBody := bytes.NewReader(deploymentManifest.Bytes())
	req, err := s.newRequest("POST", s.buildBoshURL("/deployments"), reqBody)
	if err != nil {
		return BoshTask{}, err
	}

	req.Header.Set("content-type", "text/yaml")
	res, err := s.http.Do(req)
	if err != nil {
		return BoshTask{}, err
	}

	defer res.Body.Close()
	lo.G.Debug("deployment request complete")
	var bt BoshTask
	err = json.NewDecoder(res.Body).Decode(&bt)
	return bt, err
}

func (s *Client) GetCloudConfig() (*enaml.CloudConfigManifest, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/cloud_configs?limit=1"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "text/yaml")

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var cc []CloudConfigResponseBody
	var buf bytes.Buffer
	err = json.NewDecoder(io.TeeReader(res.Body, &buf)).Decode(&cc)
	if err != nil {
		if lo.G.IsEnabledFor(logging.DEBUG) {
			lo.G.Debug(string(buf.Bytes()))
		}
		return nil, err
	}
	if len(cc) > 0 {
		return enaml.NewCloudConfigManifest([]byte(cc[0].Properties)), nil
	} else {
		return &enaml.CloudConfigManifest{}, nil
	}
}

func (s *Client) GetInfo() (*BoshInfo, error) {
	req, err := s.newRequest("GET", s.buildBoshURL("/info"), nil)
	if err != nil {
		return nil, err
	}

	res, err := s.http.Do(req)
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
