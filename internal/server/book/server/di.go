package server

import (
	"bookstore/internal/server/book/repo"
	"bookstore/internal/server/book/service"
	"bookstore/internal/server/book/usecase"
)

func BuildBookService() *service.BookService {
	br := repo.NewBookRepository()
	bu := usecase.NewBookUsecase(br)
	return service.NewBookService(bu)
}
