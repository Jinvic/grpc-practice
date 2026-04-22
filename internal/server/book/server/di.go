package server

import (
	"bookstore/internal/pkg/config"
	"bookstore/internal/pkg/database"
	"bookstore/internal/server/book/repo"
	"bookstore/internal/server/book/service"
	"bookstore/internal/server/book/usecase"

	"github.com/samber/do/v2"
)

func InitInjector(i do.Injector) error {
	do.Provide(i, config.GetConfig)
	do.Provide(i, database.NewDB)
	do.Provide(i, repo.NewBookRepository)
	do.Provide(i, usecase.NewBookUsecase)
	do.Provide(i, service.NewBookService)
	return nil
}
