package server

import (
	"bookstore/internal/server/book/repo"
	"bookstore/internal/server/book/service"
	"bookstore/internal/server/book/usecase"

	"github.com/samber/do/v2"
)

var injector do.Injector

func InitInjector() error {
	injector = do.New()
	do.Provide(injector, repo.NewBookRepository)
	do.Provide(injector, usecase.NewBookUsecase)
	do.Provide(injector, service.NewBookService)
	return nil
}
