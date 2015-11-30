package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // so we can open the databse
)

type Client struct {
	DB *gorm.DB
}

func NewClient(url string) (*Client, error) {
	db, err := gorm.Open("postgres", url)
	return &Client{DB: &db}, err
}

// GetToProcess returns a list of releases that still need to be processed
// in some way
func (c *Client) GetToProcess() (rels []Release, err error) {
	err = c.DB.Where(Release{AddedGithub: false}).
		Or(Release{AddedGithub: false}).
		Or(Release{AddedDockerHub: false}).
		Or(Release{TriggeredDockerHub: false}).
		Find(&rels).Error
	return
}

// GetToProcess returns a list of all releases
func (c *Client) GetAll() (rels []Release, err error) {
	err = c.DB.Find(&rels).Error
	return
}

func (c *Client) AddReleases(rels []Release) error {
	for _, rel := range rels {
		if err := c.DB.Create(rel).Error; err != nil {
			return err
		}
	}
	return nil
}
