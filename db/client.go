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
