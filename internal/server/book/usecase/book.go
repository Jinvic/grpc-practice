package usecase

import (
	"bookstore/internal/server/book/model/biz"
	"bookstore/internal/server/book/repo"
	"context"
	"time"

	"github.com/samber/do/v2"
)

type BookUsecase struct {
	br *repo.BookRepository
}

func NewBookUsecase(i do.Injector) (*BookUsecase, error) {
	br := do.MustInvoke[*repo.BookRepository](i)
	return &BookUsecase{br: br}, nil
}

func (u *BookUsecase) GetBook(ctx context.Context, id int64) (*biz.Book, error) {
	// debug
	return &biz.Book{
		ID:          id,
		Title:       "The Great Gatsby",
		Author:      "F. Scott Fitzgerald",
		Price:       10.99,
		ISBN:        "978-0-7432-1967-1",
		Publisher:   "Scribner",
		PublishedAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   time.Now(),
	}, nil
}

func (u *BookUsecase) CreateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	return nil, nil
}

func (u *BookUsecase) ListBooks(ctx context.Context, page int, pageSize int) ([]*biz.Book, error) {
	return nil, nil
}

func (u *BookUsecase) UpdateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	return nil, nil
}

func (u *BookUsecase) DeleteBook(ctx context.Context, id int64) error {
	return nil
}
