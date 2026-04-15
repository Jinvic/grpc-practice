package biz

import (
	commonv1 "bookstore/api/common/v1"
	"bookstore/internal/server/book/model/db"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func (b *Book) ToV1Book() *commonv1.Book {
	return &commonv1.Book{
		Id:          b.ID,
		Title:       b.Title,
		Author:      b.Author,
		Price:       b.Price,
		Isbn:        b.ISBN,
		Publisher:   b.Publisher,
		PublishedAt: timestamppb.New(b.PublishedAt),
	}
}

func FromV1Book(book *commonv1.Book) *Book {
	return &Book{
		ID:          book.Id,
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		ISBN:        book.Isbn,
		Publisher:   book.Publisher,
		PublishedAt: book.PublishedAt.AsTime(),
	}
}

func FromDBBook(book *db.Book) *Book {
	return &Book{
		ID:          book.ID,
		Status:      BookStatus(book.Status),
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		ISBN:        book.ISBN,
		Publisher:   book.Publisher,
		PublishedAt: time.Unix(book.PublishedAt, 0),
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
		DeletedAt:   book.DeletedAt.Time,
	}
}

func ToDBBook(book *Book) *db.Book {
	return &db.Book{
		ID:          book.ID,
		Status:      int(book.Status),
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		ISBN:        book.ISBN,
		Publisher:   book.Publisher,
		PublishedAt: book.PublishedAt.Unix(),
	}
}
