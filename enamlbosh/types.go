package enamlbosh

type Client struct {
	user string
	pass string
	host string
	port int
}

func NewClient(user, pass, host string, port int) *Client {
	return &Client{
		user: user,
		pass: pass,
		host: host,
		port: port,
	}
}
