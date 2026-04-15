package db

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID          int64   `gorm:"column:id;primaryKey"`
	Status      int     `gorm:"column:status"`
	Title       string  `gorm:"column:title"`
	Author      string  `gorm:"column:author"`
	Price       float64 `gorm:"column:price"`
	ISBN        string  `gorm:"column:isbn"`
	Publisher   string  `gorm:"column:publisher"`
	PublishedAt int64   `gorm:"column:published_at"`
}

func (Book) TableName() string {
	return "books"
}
