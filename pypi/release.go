package pypi

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gosuri/uiprogress"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

// ReleasesSince queries PyPi for all the releases sine `time`, where
// `time` is seconds from epocj
func (c *Client) ReleasesSince(t time.Time) ([]db.Release, error) {
	changes := [][]interface{}{}

	err := c.client.Call("changelog", []interface{}{t.Unix(), true}, &changes)
	if err != nil {
		return nil, err
	}
	releases := []db.Release{}
	for _, change := range changes {
		if change[3] == "new release" {
			releases = append(releases, db.Release{
				Name:    change[0].(string),
				Version: change[1].(string),
				Time:    time.Unix(change[2].(int64), 0),
			})
		}

	}
	return releases, nil
}

// AllReleases queries PyPi for all the release ever
func (c *Client) AllReleases() ([]db.Release, error) {
	names, err := c.names()
	if err != nil {
		return nil, err
	}
	uiprogress.Start()
	bar := uiprogress.AddBar(len(names))
	bar.AppendCompleted()
	bar.PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf(
			"%v; %d/s",
			b.Current(),
			int(float64(b.Current())/b.TimeElapsed().Seconds()),
		)
	})
	releases := make(chan db.Release)
	c.addReleases(names, releases, bar)
	close(releases)
	return releaseChanToSlice(releases), nil
}

func releaseChanToSlice(releasesChan <-chan db.Release) []db.Release {
	releases := []db.Release{}
	for rel := range releasesChan {
		releases = append(releases, rel)
	}
	return releases
}

func (c *Client) addReleases(names []string, releases chan<- db.Release, bar *uiprogress.Bar) {
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		name := name
		go func() {
			c.addVersions(name, releases, bar)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (c *Client) addReleasesLimitedWorkers(names []string, releases chan<- db.Release, bar *uiprogress.Bar) {
	numWorkers := 1000
	namesChan := make(chan string, len(names))
	for _, name := range names {
		namesChan <- name
	}
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		func() {
			defer wg.Done()
			select {
			case name := <-namesChan:
				c.addVersions(name, releases, bar)
			default:
				return
			}
		}()
	}
	wg.Wait()
}

func (c *Client) addVersions(name string, releases chan<- db.Release, bar *uiprogress.Bar) {
	_log := logrus.WithField("name", name)
	_log.Debug("Getting versions")
	versions, err := c.versions(name)
	if err != nil {
		_log.WithError(err).Warning("Can't get versions")
	}
	bar.Incr()
	_log.WithField("count", len(versions)).Debug("Got versions")
	for _, version := range versions {
		time, err := c.releaseTime(name, version)
		if err != nil {
			_log.WithError(err).Warning("Can't get release time")
		}
		releases <- db.Release{
			Name:    name,
			Version: version,
			Time:    time,
		}
	}
}

func (c *Client) names() (names []string, err error) {
	logrus.Debug("Getting packages")
	err = c.client.Call("list_packages", []interface{}{}, &names)
	logrus.WithField("count", len(names)).Info("Got all packages")
	return
}

func (c *Client) versions(name string) (versions []string, err error) {
	err = c.client.Call("package_releases", []interface{}{name}, &versions)
	return
}

type ReleaseURL struct {
	UploadTime time.Time `xmlrpc:"upload_time"`
}

func (c *Client) releaseTime(name string, version string) (t time.Time, err error) {
	var rus []ReleaseURL
	err = c.client.Call("release_urls", []interface{}{name, version}, &rus)
	if err == nil {
		t = rus[0].UploadTime
	}
	return
}
