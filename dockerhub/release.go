package dockerhub

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Release interface {
	DockerHubTag() string
	DockerHubRepoShortDescription() string
	DockerHubRepoFullDescription() string
	DockerfilePath() string
	GitTagName() string
}

func (c *Client) AddRelease(rel Release) (err error) {
	log.Debug("DockerHub: checking if repo exists")
	repoExists, err := c.checkRepoExists(rel)
	if err != nil {
		return wrapError(err, "checking repo exists")
	}
	if repoExists {
		log.Debug("DockerHub: Exists; checking if automated builds exist")
		automatedExists, err := c.checkAutomatedExists(rel)
		if err != nil {
			return wrapError(err, "checking automated builds exist")
		}
		if automatedExists {
			log.Debug("DockerHub: Exists; checking if build exists")
			buildExists, err := c.checkBuildExists(rel)
			if err != nil {
				return wrapError(err, "checking build exists")
			}
			if !buildExists {
				log.Debug("DockerHub: creating build")
				if err = c.createBuild(rel); err != nil {
					return wrapError(err, "creating build")
				}
			} else {
				log.Debug("DockerHub: build already exists")
			}
		} else {
			log.Debug("DockerHub: doesn't exist; deleting whole repo")
			c.DeleteRepo(rel)
			repoExists = !repoExists
		}

	}

	if !repoExists {
		log.Debug("DockerHub: Doesn't exist; creating repo and build")
		if err = c.createRepoAndBuild(rel); err != nil {
			return wrapError(err, "creating repo and build")
		}
	} else {
		log.Debug("DockerHub: Exists; checking if build exists")
		buildExists, err := c.checkBuildExists(rel)
		if err != nil {
			return wrapError(err, "checking build exists")
		}
		if !buildExists {
			log.Debug("DockerHub: creating build")
			if err = c.createBuild(rel); err != nil {
				return wrapError(err, "creating build")
			}
		} else {
			log.Debug("DockerHub: build already exists")
		}
	}
	log.Debug("DockerHub: setting full description")
	err = c.setFullDescription(rel)
	return wrapError(err, "setting full description")
}

func (c *Client) TriggerRelease(rel Release) error {
	log.Debug("DockerHub: triggering build")
	return wrapError(c.triggerBuild(rel), "triggering build")
}

func (c *Client) checkRepoExists(rel Release) (bool, error) {
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

func (c *Client) checkAutomatedExists(rel Release) (bool, error) {
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

func (c *Client) createRepoAndBuild(rel Release) error {
	body := struct {
		Active            bool       `json:"active"`
		BuildTags         []buildTag `json:"build_tags"`
		Description       string     `json:"description"`
		DockerHubRepoName string     `json:"dockerHub_repo_name"`
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
			Name:               rel.DockerHubTag(),
			SourceType:         "Tag",
			SourceName:         rel.GitTagName(),
			DockerfileLocation: rel.DockerfilePath(),
		}},
		Description:       rel.DockerHubRepoShortDescription(),
		DockerHubRepoName: fmt.Sprintf("%v/%v", c.Owner, c.Name),
		IsPrivate:         false,
		Name:              c.Name,
		Namespace:         c.Owner,
		Provider:          "github",
		VCSRepoName:       fmt.Sprintf("%v/%v", c.github.Owner, c.github.Name),
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

func (c *Client) checkBuildExists(rel Release) (bool, error) {
	var resJSON autobuildReponse
	if _, err := c.callRepo(rel, "autobuild/", "GET", "", 200, &resJSON); err != nil {
		return false, err
	}
	for _, bt := range resJSON.BuildTags {
		if bt.Name == rel.DockerHubTag() {
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

func (c *Client) createBuild(rel Release) error {
	body := completeBuildTag{
		buildTag: buildTag{
			Name:               rel.DockerHubTag(),
			SourceType:         "Tag",
			SourceName:         rel.GitTagName(),
			DockerfileLocation: rel.DockerfilePath(),
		},
		IsNew:     true,
		Namespace: c.Owner,
		RepoName:  c.Name,
	}
	res, err := c.callRepo(rel, "autobuild/tags/", "POST", body, 201, nil)
	if err != nil {
		return wrapError(err, "adding build")
	}
	return wrapError(res.Body.Close(), "closing body on autobuild/tags")
}

func (c *Client) triggerBuild(rel Release) (err error) {
	_, err = c.callRepo(rel, "autobuild/trigger-build/", "POST", "", 201, nil)
	return
}

type fullDescriptionBody struct {
	FullDescription string `json:"full_description"`
}

func (c *Client) setFullDescription(rel Release) (err error) {
	_, err = c.callRepo(rel, "", "PATCH", fullDescriptionBody{rel.DockerHubRepoFullDescription()}, 200, nil)
	return
}
