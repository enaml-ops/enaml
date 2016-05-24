package enamlbosh

import "net/http"

type Client struct {
	user string
	pass string
	host string
	port int
}

type HttpClientDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewClient(user, pass, host string, port int) *Client {
	return &Client{
		user: user,
		pass: pass,
		host: host,
		port: port,
	}
}

type CloudConfigResponseBody struct {
	Properties string
}
