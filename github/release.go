package github

import (
	"fmt"
	"strings"

	log "github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/saulshanabrook/pypi-dockerhub/release"
)

func (c *Client) AddRelease(rel *release.Release) error {
	rel.Log().Debug("Github: Checking tag exists")
	tagExists, err := c.tagExists(rel)
	if err != nil {
		return wrapError(err, "checking if tag exists")
	}
	if tagExists {
		rel.Log().Debug("Github: Already exists; skipping")
		return nil
	}

	rel.Log().Debug("Github: Getting current dockerfile")
	prevFileContent, err := c.currentDockerfile(rel)
	var rcr *github.RepositoryContentResponse
	if err == nil {
		rel.Log().Debug("Github: Found; updating it with new version")
		rcr, err = c.updateDockerfile(rel, *prevFileContent.SHA)
	} else {
		rel.Log().Debug("Github: Didn't find; adding it")
		rcr, err = c.addDockerfile(rel)
	}
	if err != nil {
		return wrapError(err, "updating/adding dockerfile")
	}
	rel.Log().Debug("Github: Adding tag")
	return wrapError(c.addTag(rel, rcr), "adding tag")
}

func (c *Client) currentDockerfile(rel *release.Release) (content *github.RepositoryContent, err error) {
	content, _, _, err = c.client.Repositories.GetContents(c.owner, c.repo, rel.DockerfilePath(), &github.RepositoryContentGetOptions{})
	return
}

func (c *Client) updateDockerfile(rel *release.Release, sha string) (rcr *github.RepositoryContentResponse, err error) {
	message := rel.GithubCommitMessage()
	newFileContent := github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(rel.DockerfileContents()),
		SHA:     &sha,
	}
	rcr, _, err = c.client.Repositories.UpdateFile(c.owner, c.repo, rel.DockerfilePath(), &newFileContent)
	return
}

func (c *Client) addDockerfile(rel *release.Release) (rcr *github.RepositoryContentResponse, err error) {
	message := rel.GithubCommitMessage()
	newFileContent := github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(rel.DockerfileContents()),
	}
	rcr, _, err = c.client.Repositories.CreateFile(c.owner, c.repo, rel.DockerfilePath(), &newFileContent)
	return
}

func (c *Client) addTagRef(tag *github.Tag) (err error) {
	refName := fmt.Sprintf("tags/%v", *tag.Tag)
	ref := github.Reference{
		Ref:    &refName,
		Object: tag.Object,
	}
	_, _, err = c.client.Git.CreateRef(c.owner, c.repo, &ref)
	if err != nil && strings.Contains(err.Error(), "Reference already exists") {
		log.Debug("Github: Ref already existed; deleting and trying again")
		_, err = c.client.Git.DeleteRef(c.owner, c.repo, refName)
		if err != nil {
			return wrapError(err, "deleting ref")
		}
		c.addTagRef(tag)
	}
	return wrapError(err, "creating ref")
}
func (c *Client) addTag(rel *release.Release, rcr *github.RepositoryContentResponse) (err error) {
	tagName := rel.GithubTagName()
	tagMessage := rel.GithubTagMessage(rcr)
	objectType := "commit"
	tag := &github.Tag{
		Tag:     &tagName,
		SHA:     rcr.SHA,
		Message: &tagMessage,
		Object: &github.GitObject{
			Type: &objectType,
			URL:  rcr.URL,
			SHA:  rcr.SHA,
		},
	}
	tag, _, err = c.client.Git.CreateTag(c.owner, c.repo, tag)
	if err != nil {
		return err
	}
	return c.addTagRef(tag)
}

func (c *Client) tagExists(rel *release.Release) (bool, error) {
	t1, t2, err := c.client.Git.GetRef(
		c.owner, c.repo, fmt.Sprintf("tags/%v", rel.GithubTagName()))
	fmt.Printf("%v\n\n%v\n\n%v", t1, t2, err)
	if err != nil && (strings.Contains(err.Error(), "404 Not Found") || strings.Contains(err.Error(), "409 Git Repository is empty")) {
		return false, nil
	}
	return true, err
}
