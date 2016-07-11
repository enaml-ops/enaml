package enamlbosh

import "net/http"

//Client a bosh client object
type Client struct {
	user string
	pass string
	host string
	port int
}

//HttpClientDoer - interface for a http.Client.Doer
type HttpClientDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

//NewClient - constrcutor for a bosh client
func NewClient(user, pass, host string, port int) *Client {
	return &Client{
		user: user,
		pass: pass,
		host: host,
		port: port,
	}
}

//BoshInfo - info object for bosh
type BoshInfo struct {
	Name               string
	UUID               string
	Version            string
	User               string
	CPI                string
	UserAuthentication map[string]interface{}
	Features           map[string]interface{}
}

//CloudConfigResponseBody - response body struct for get cloud config calls
type CloudConfigResponseBody struct {
	Properties string
}

//BoshTask - an object representing a bosh task response
type BoshTask struct {
	ID          int
	State       string
	Description string
	Timestamp   int
	Result      string
	User        string
}

//DeployedStemcell - response of stemcells already deployed
type DeployedStemcell struct {
	Name    string
	Version string
}
