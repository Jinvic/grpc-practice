package repo

import (
	"bookstore/internal/server/book/model/biz"
	"context"
)

type BookRepository struct{}

func NewBookRepository() *BookRepository {
	return &BookRepository{}
}

func (r *BookRepository) GetBook(ctx context.Context, id int64) (*biz.Book, error) {
	return nil, nil
}

func (r *BookRepository) CreateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	return nil, nil
}

func (r *BookRepository) ListBooks(ctx context.Context, page int, pageSize int) ([]*biz.Book, error) {
	return nil, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	return nil, nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id int64) error {
	return nil
}
