package pypi

import "github.com/kolo/xmlrpc"

type Client struct {
	client *xmlrpc.Client
}

func NewClient() (*Client, error) {
	xmlClient, err := xmlrpc.NewClient("https://pypi.python.org/pypi", nil)
	return &Client{client: xmlClient}, err
}
