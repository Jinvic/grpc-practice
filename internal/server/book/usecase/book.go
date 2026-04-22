package usecase

import (
	"bookstore/internal/server/book/model/biz"
	"bookstore/internal/server/book/repo"
	"context"

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
	return u.br.GetBook(ctx, id)
}

func (u *BookUsecase) CreateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	return u.br.CreateBook(ctx, book)
}

func (u *BookUsecase) ListBooks(ctx context.Context, page int, pageSize int) ([]*biz.Book, error) {
	return u.br.ListBooks(ctx, page, pageSize)
}

func (u *BookUsecase) UpdateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	return u.br.UpdateBook(ctx, book)
}

func (u *BookUsecase) DeleteBook(ctx context.Context, id int64) error {
	return u.br.DeleteBook(ctx, id)
}
