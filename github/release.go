package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/github"

	log "github.com/Sirupsen/logrus"
)

type Release interface {
	DockerfileContents() string
	DockerfilePath() string
	GitCommitMessage() string
	GitTagName() string
	GitTagMessage() string
}

// AddRelease takes Release r and adds it to Github.
// 1. Commit contents r.FileContents at r.FilePath with message r.CommitMessage
// 2. Add tag r.TagName with message r.TagMessage
func (c *Client) AddRelease(rel Release) error {
	// log.Debug("Github: Checking tag exists")
	// tagExists, err := c.tagExists(rel)
	// if err != nil {
	// 	return wrapError(err, "checking if tag exists")
	// }
	// if tagExists {
	// 	log.Debug("Github: Already exists; skipping")
	// 	return nil
	// }

	log.Debug("Github: Getting current dockerfile")
	prevFileContent, err := c.currentDockerfile(rel)
	var rcr *github.RepositoryContentResponse
	if err == nil {
		log.Debug("Github: Found; updating it with new version")
		rcr, err = c.updateDockerfile(rel, *prevFileContent.SHA)
	} else {
		log.Debug("Github: Didn't find; adding it")
		rcr, err = c.addDockerfile(rel)
	}
	if err != nil {
		return wrapError(err, "updating/adding dockerfile")
	}
	log.Debug("Github: Adding tag")
	return wrapError(c.addTag(rel, rcr), "adding tag")
}

func (c *Client) currentDockerfile(rel Release) (content *github.RepositoryContent, err error) {
	content, _, _, err = c.client.Repositories.GetContents(c.Owner, c.Name, rel.DockerfilePath(), &github.RepositoryContentGetOptions{})
	return
}

func (c *Client) updateDockerfile(rel Release, sha string) (rcr *github.RepositoryContentResponse, err error) {
	message := rel.GitCommitMessage()
	newFileContent := github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(rel.DockerfileContents()),
		SHA:     &sha,
	}
	rcr, _, err = c.client.Repositories.UpdateFile(c.Owner, c.Name, rel.DockerfilePath(), &newFileContent)
	return
}

func (c *Client) addDockerfile(rel Release) (rcr *github.RepositoryContentResponse, err error) {
	message := rel.GitCommitMessage()
	newFileContent := github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(rel.DockerfileContents()),
	}
	rcr, _, err = c.client.Repositories.CreateFile(c.Owner, c.Name, rel.DockerfilePath(), &newFileContent)
	return
}

func (c *Client) addTagRef(tag *github.Tag) (err error) {
	refName := fmt.Sprintf("tags/%v", *tag.Tag)
	ref := github.Reference{
		Ref:    &refName,
		Object: tag.Object,
	}
	_, _, err = c.client.Git.CreateRef(c.Owner, c.Name, &ref)
	if err != nil && strings.Contains(err.Error(), "Reference already exists") {
		log.Debug("Github: Ref already existed; deleting and trying again")
		_, err = c.client.Git.DeleteRef(c.Owner, c.Name, refName)
		if err != nil {
			return wrapError(err, "deleting ref")
		}
		c.addTagRef(tag)
	}
	return wrapError(err, "creating ref")
}
func (c *Client) addTag(rel Release, rcr *github.RepositoryContentResponse) (err error) {
	tagName := rel.GitTagName()
	tagMessage := rel.GitTagMessage()
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
	tag, _, err = c.client.Git.CreateTag(c.Owner, c.Name, tag)
	if err != nil {
		return err
	}
	return c.addTagRef(tag)
}

func (c *Client) tagExists(rel Release) (bool, error) {
	_, _, err := c.client.Git.GetRef(
		c.Owner, c.Name, fmt.Sprintf("tags/%v", rel.GitTagName()))
	notFound := err != nil && strings.Contains(err.Error(), "404 Not Found")
	repoEmpty := err != nil && strings.Contains(err.Error(), "409 Git Repository is empty")
	_, multipleRefsRet := err.(*json.UnmarshalTypeError)
	if notFound || repoEmpty || multipleRefsRet {
		return false, nil
	}
	return true, err
}
