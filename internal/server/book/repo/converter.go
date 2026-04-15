package repo

import (
	"bookstore/internal/server/book/model/biz"
	"bookstore/internal/server/book/model/db"
	"time"
)

func BizToDBBook(b *biz.Book) *db.Book {
	if b == nil {
		return nil
	}
	return &db.Book{
		ID:          b.ID,
		Status:      int(b.Status),
		Title:       b.Title,
		Author:      b.Author,
		Price:       b.Price,
		ISBN:        b.ISBN,
		Publisher:   b.Publisher,
		PublishedAt: b.PublishedAt.Unix(),
	}
}

func DBToBizBook(d *db.Book) *biz.Book {
	if d == nil {
		return nil
	}
	return &biz.Book{
		ID:          d.ID,
		Status:      biz.BookStatus(d.Status),
		Title:       d.Title,
		Author:      d.Author,
		Price:       d.Price,
		ISBN:        d.ISBN,
		Publisher:   d.Publisher,
		PublishedAt: time.Unix(d.PublishedAt, 0),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt.Time,
	}
}

func BizListToDBBookList(books []*biz.Book) []*db.Book {
	if books == nil {
		return nil
	}
	result := make([]*db.Book, len(books))
	for i, b := range books {
		result[i] = BizToDBBook(b)
	}
	return result
}

func DBListToBizBookList(books []*db.Book) []*biz.Book {
	if books == nil {
		return nil
	}
	result := make([]*biz.Book, len(books))
	for i, d := range books {
		result[i] = DBToBizBook(d)
	}
	return result
}
