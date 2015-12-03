package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type LastUpdateTime struct {
	ID   uint `gorm:"primary_key"`
	Time time.Time
}

func (c *Client) whereLastUpdateTime() *gorm.DB {
	return c.DB.Where(LastUpdateTime{ID: 1})
}

func (c *Client) SetLastUpdateTime(t time.Time) error {
	var lat LastUpdateTime
	db := c.whereLastUpdateTime().Assign(LastUpdateTime{Time: t}).FirstOrCreate(&lat)
	return db.Error
}

func (c *Client) GetLastUpdateTime() (time.Time, error) {
	var lut LastUpdateTime
	db := c.whereLastUpdateTime().First(&lut)
	return lut.Time, db.Error
}
