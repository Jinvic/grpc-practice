package biz

import "time"

type BookStatus int

const (
	BookStatusUnspecified BookStatus = iota
	BookStatusAvailable
	BookStatusUnavailable
	BookStatusBorrowed
	BookStatusLost
	BookStatusReserved
)

type Book struct {
	ID          int64
	Status      BookStatus
	Title       string
	Author      string
	Price       float64
	ISBN        string
	Publisher   string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
