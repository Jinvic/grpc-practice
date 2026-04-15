package usecase

import (
	"bookstore/internal/server/book/model/biz"
	"bookstore/internal/server/book/repo"
	"context"
)

type BookUsecase struct {
	br *repo.BookRepository
}

func NewBookUsecase(br *repo.BookRepository) *BookUsecase {
	return &BookUsecase{br: br}
}

func (u *BookUsecase) GetBook(ctx context.Context, id int64) (*biz.Book, error) {
	return nil, nil
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
