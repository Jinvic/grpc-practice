package repo

import (
	"bookstore/internal/server/book/model/biz"
	"bookstore/internal/server/book/model/db"
	"context"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(i do.Injector) (*BookRepository, error) {
	gormDB := do.MustInvoke[*gorm.DB](i)
	if err := gormDB.AutoMigrate(&db.Book{}); err != nil {
		return nil, err
	}
	return &BookRepository{db: gormDB}, nil
}

func (r *BookRepository) GetBook(ctx context.Context, id int64) (*biz.Book, error) {
	var book db.Book
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&book).Error; err != nil {
		return nil, err
	}
	return DBToBizBook(&book), nil
}

func (r *BookRepository) CreateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	dbBook := BizToDBBook(book)
	if err := r.db.WithContext(ctx).Create(dbBook).Error; err != nil {
		return nil, err
	}
	return DBToBizBook(dbBook), nil
}

func (r *BookRepository) ListBooks(ctx context.Context, page int, pageSize int) ([]*biz.Book, error) {
	offset := (page - 1) * pageSize
	var books []*db.Book
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, err
	}
	return DBListToBizBookList(books), nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *biz.Book) (*biz.Book, error) {
	dbBook := BizToDBBook(book)
	if err := r.db.WithContext(ctx).Save(dbBook).Error; err != nil {
		return nil, err
	}
	return DBToBizBook(dbBook), nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&db.Book{}).Error; err != nil {
		return err
	}
	return nil
}
