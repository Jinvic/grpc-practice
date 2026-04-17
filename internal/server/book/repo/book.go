package repo

import (
	"bookstore/internal/server/book/model/biz"
	"context"

	"github.com/samber/do/v2"
)

type BookRepository struct{}

func NewBookRepository(i do.Injector) (*BookRepository, error) {
	return &BookRepository{}, nil
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
