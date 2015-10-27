package pypi

import (
	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/kolo/xmlrpc"
	"github.com/saulshanabrook/pypi-dockerhub/release"
)

type Client struct {
	client *xmlrpc.Client
}

func NewClient() (*Client, error) {
	xmlClient, err := xmlrpc.NewClient("https://pypi.python.org/pypi", nil)
	return &Client{client: xmlClient}, err
}

// ReleasesSince returns a list of Releases since that time
// with the most recent last
func (c *Client) ReleasesSince(time int64) ([]*release.Release, error) {
	changes := [][]interface{}{}

	err := c.client.Call("changelog", []interface{}{time, true}, &changes)
	if err != nil {
		return nil, err
	}
	releases := []*release.Release{}
	for _, change := range changes {
		if change[3] == "new release" {
			releases = append(releases, &release.Release{
				Name:    change[0].(string),
				Version: change[1].(string),
				Time:    change[2].(int64),
			})
		}
	}
	return releases, nil
}
