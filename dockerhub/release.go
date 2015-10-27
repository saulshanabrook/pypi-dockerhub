package dockerhub

import (
	"fmt"

	"github.com/saulshanabrook/pypi-dockerhub/release"
)

func (c *Client) AddRelease(rel *release.Release) (err error) {
	rel.Log().Debug("Dockerhub: checking if repo exists")
	repoExists, err := c.checkRepoExists(rel)
	if err != nil {
		return wrapError(err, "checking repo exists")
	}
	if repoExists {
		rel.Log().Debug("Dockerhub: Exists; checking if automated builds exist")
		automatedExists, err := c.checkAutomatedExists(rel)
		if err != nil {
			return wrapError(err, "checking automated builds exist")
		}
		if automatedExists {
			rel.Log().Debug("Dockerhub: Exists; checking if build exists")
			buildExists, err := c.checkBuildExists(rel)
			if err != nil {
				return wrapError(err, "checking build exists")
			}
			if !buildExists {
				rel.Log().Debug("Dockerhub: creating build")
				if err = c.createBuild(rel); err != nil {
					return wrapError(err, "creating build")
				}
			} else {
				rel.Log().Debug("Dockerhub: build already exists")
			}
		} else {
			rel.Log().Debug("Dockerhub: doesn't exist; deleting whole repo")
			c.DeleteRepo(rel)
			repoExists = !repoExists
		}

	}

	if !repoExists {
		rel.Log().Debug("Dockerhub: Doesn't exist; creating repo and build")
		if err = c.createRepoAndBuild(rel); err != nil {
			return wrapError(err, "creating repo and build")
		}
	} else {
		rel.Log().Debug("Dockerhub: Exists; checking if build exists")
		buildExists, err := c.checkBuildExists(rel)
		if err != nil {
			return wrapError(err, "checking build exists")
		}
		if !buildExists {
			rel.Log().Debug("Dockerhub: creating build")
			if err = c.createBuild(rel); err != nil {
				return wrapError(err, "creating build")
			}
		} else {
			rel.Log().Debug("Dockerhub: build already exists")
		}
	}
	rel.Log().Debug("Dockerhub: triggering build")
	if err = c.triggerBuild(rel); err != nil {
		return wrapError(err, "triggering build")
	}

	rel.Log().Debug("Dockerhub: setting full description")
	err = c.setFullDescription(rel)
	return wrapError(err, "setting full description")
}

func (c *Client) checkRepoExists(rel *release.Release) (bool, error) {
	res, err := c.callRepo(rel, "", "GET", nil, 0, nil)
	if err != nil {
		return false, err
	}
	if res.StatusCode == 404 {
		return false, nil
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, wrongResponseError(res, "repo should have either been a 404 or a 200")
}

func (c *Client) checkAutomatedExists(rel *release.Release) (bool, error) {
	res, err := c.callRepo(rel, "autobuild/", "GET", nil, 0, nil)
	if err != nil {
		return false, err
	}
	if res.StatusCode == 404 {
		return false, nil
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, wrongResponseError(res, "autobuild should have either been a 404 or a 200")
}

type buildTag struct {
	Name               string `json:"name"`
	SourceType         string `json:"source_type"`
	SourceName         string `json:"source_name"`
	DockerfileLocation string `json:"dockerfile_location"`
}

func (c *Client) createRepoAndBuild(rel *release.Release) error {
	body := struct {
		Active            bool       `json:"active"`
		BuildTags         []buildTag `json:"build_tags"`
		Description       string     `json:"description"`
		DockerhubRepoName string     `json:"dockerhub_repo_name"`
		IsPrivate         bool       `json:"is_private"`
		Name              string     `json:"name"`
		Namespace         string     `json:"namespace"`
		Provider          string     `json:"provider"`
		VCSRepoName       string     `json:"vcs_repo_name"`
	}{
		Active: false,
		BuildTags: []buildTag{{
			Name:               "latest",
			SourceType:         "Branch",
			SourceName:         "master",
			DockerfileLocation: rel.DockerfilePath(),
		}, {
			Name:               rel.DockerhubTag(),
			SourceType:         "Tag",
			SourceName:         rel.GithubTagName(),
			DockerfileLocation: rel.DockerfilePath(),
		}},
		Description:       rel.DockerhubRepoShortDescription(),
		DockerhubRepoName: fmt.Sprintf("%v/%v", c.dockerhubOwner, rel.DockerhubName()),
		IsPrivate:         false,
		Name:              rel.DockerhubName(),
		Namespace:         c.dockerhubOwner,
		Provider:          "github",
		VCSRepoName:       fmt.Sprintf("%v/%v", c.githubOwner, c.githubRepo),
	}
	res, err := c.callRepo(rel, "autobuild/", "POST", body, 201, nil)
	if err != nil {
		return err
	}
	if err = res.Body.Close(); err != nil {
		return wrapError(err, "closing body on autobuild/")
	}

	var resJSON struct{ Active bool }
	res, err = c.callRepo(rel, "autobuild/", "PATCH", struct {
		Active bool `json:"active"`
	}{false}, 200, &resJSON)
	if err != nil {
		return wrapError(err, "turning off autobuild")
	}
	if resJSON.Active != false {
		return fmt.Errorf("Couldnt turn off autobuilding")
	}
	err = res.Body.Close()
	return wrapError(err, "closing body on PATCH autobuild/")
}

type autobuildReponse struct {
	BuildTags []buildTag `json:"build_tags"`
}

func (c *Client) checkBuildExists(rel *release.Release) (bool, error) {
	var resJSON autobuildReponse
	if _, err := c.callRepo(rel, "autobuild/", "GET", "", 200, &resJSON); err != nil {
		return false, err
	}
	for _, bt := range resJSON.BuildTags {
		if bt.Name == rel.DockerhubTag() {
			return true, nil
		}
	}
	return false, nil
}

type completeBuildTag struct {
	buildTag
	IsNew     bool   `json:"is_new"`
	Namespace string `json:"namespace"`
	RepoName  string `json:"repo_name"`
}

func (c *Client) createBuild(rel *release.Release) error {
	body := completeBuildTag{
		buildTag: buildTag{
			Name:               rel.DockerhubTag(),
			SourceType:         "Tag",
			SourceName:         rel.GithubTagName(),
			DockerfileLocation: rel.DockerfilePath(),
		},
		IsNew:     true,
		Namespace: c.dockerhubOwner,
		RepoName:  rel.DockerhubName(),
	}
	res, err := c.callRepo(rel, "autobuild/tags/", "POST", body, 201, nil)
	if err != nil {
		return wrapError(err, "adding build")
	}
	return wrapError(res.Body.Close(), "closing body on autobuild/tags")
}

func (c *Client) triggerBuild(rel *release.Release) (err error) {
	_, err = c.callRepo(rel, "autobuild/trigger-build/", "POST", "", 201, nil)
	return
}

type fullDescriptionBody struct {
	FullDescription string `json:"full_description"`
}

func (c *Client) setFullDescription(rel *release.Release) (err error) {
	_, err = c.callRepo(rel, "", "PATCH", fullDescriptionBody{rel.DockerhubRepoFullDescription()}, 200, nil)
	return
}
