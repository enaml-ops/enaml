package enamlbosh

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
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

func transport(insecureSkipVerify bool) *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
}

// NewClientBasic creates a bosh client configured for basic auth.
// It can be configured to ignore SSL warnings by setting sslIgnore to true.
func NewClientBasic(user, pass, host string, port int, sslIgnore bool) *Client {
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
	return c
}

// NewClientUAA creates a bosh client configured for UAA and obtains
// a token from the UAA server.
//
// The user and pass arguments are the ops manager credentials.
// The id and secret arguments are the client ID and client secret for UAA.
// The host and port arguments are for the bosh director.
func NewClientUAA(user, pass, id, secret, host string, port int, uaaURL string, sslIgnore bool) (*Client, error) {
	c := NewClientBasic(user, pass, host, port, sslIgnore)
	cfg := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", uaaURL),
			TokenURL: fmt.Sprintf("%s/oauth/token", uaaURL),
		},
	}
	tok, err := cfg.PasswordCredentialsToken(context.Background(), user, pass)
	c.token = tok
	return c, err
}

// BoshInfo contains data about a bosh.
type BoshInfo struct {
	Name               string
	UUID               string
	Version            string
	User               string
	CPI                string
	UserAuthentication map[string]interface{}
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
