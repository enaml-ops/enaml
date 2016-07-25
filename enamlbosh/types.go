package enamlbosh

import (
	"net/http"

	"golang.org/x/oauth2"
)

// Client provides an interface to the bosh director.
type Client struct {
	user  string
	pass  string
	host  string
	port  int
	http  *http.Client
	token *oauth2.Token
}

// BoshInfo contains data about a bosh.
type BoshInfo struct {
	Name               string
	UUID               string
	Version            string
	User               string
	CPI                string
	UserAuthentication map[string]interface{} `json:"user_authentication"`
	Features           map[string]interface{}
}

// CloudConfigResponseBody is the response for cloud config GET calls.
type CloudConfigResponseBody struct {
	Properties string
}

// BoshTask represents a bosh task response.
type BoshTask struct {
	ID          int
	State       string
	Description string
	Timestamp   int
	Result      string
	User        string
}

// DeployedStemcell is the response for stemcells that have already been deployed.
type DeployedStemcell struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	OS      string `json:"operating_system"`
}
