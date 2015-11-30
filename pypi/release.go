package pypi

import (
	"github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

// ReleasesSince returns a list of Releases since that time
// with the most recent last
func (c *Client) ReleasesSince(time int64) ([]db.Release, error) {
	changes := [][]interface{}{}

	err := c.client.Call("changelog", []interface{}{time, true}, &changes)
	if err != nil {
		return nil, err
	}
	releases := []db.Release{}
	for _, change := range changes {
		if change[3] == "new release" {
			releases = append(releases, db.Release{
				Name:    change[0].(string),
				Version: change[1].(string),
				Time:    change[2].(int64),
			})
		}
	}
	return releases, nil
}

// AllReleases returns all the releases
func (c *Client) AllReleases() ([]db.Release, error) {
	logrus.Debug("Getting packages")
	names, err := c.names()
	if err != nil {
		return nil, err
	}
	logrus.WithField("count", len(names)).Debug("Got all packages")
	releases := []db.Release{}
	for _, name := range names {
		logrus.WithField("name", name).Debug("Getting versions")
		versions, err := c.versions(name)
		if err != nil {
			return nil, err
		}
		logrus.WithField("count", len(versions)).Debug("Got versions")
		for _, version := range versions {
			releases = append(releases, db.Release{
				Name:    name,
				Version: version,
			})
		}
	}
	return releases, nil
}

func (c *Client) names() (names []string, err error) {
	err = c.client.Call("list_packages", []interface{}{}, &names)
	return
}

func (c *Client) versions(name string) (versions []string, err error) {
	err = c.client.Call("package_releases", []interface{}{name}, &versions)
	return
}
