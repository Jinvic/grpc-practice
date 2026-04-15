package service

import (
	commonv1 "bookstore/api/common/v1"
	"bookstore/internal/server/book/model/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func BizToV1Book(b *biz.Book) *commonv1.Book {
	if b == nil {
		return nil
	}
	return &commonv1.Book{
		Id:          b.ID,
		Title:       b.Title,
		Author:      b.Author,
		Price:       b.Price,
		Isbn:        b.ISBN,
		Publisher:   b.Publisher,
		PublishedAt: timestamppb.New(b.PublishedAt),
	}
}

func V1ToBizBook(v *commonv1.Book) *biz.Book {
	if v == nil {
		return nil
	}
	return &biz.Book{
		ID:          v.Id,
		Title:       v.Title,
		Author:      v.Author,
		Price:       v.Price,
		ISBN:        v.Isbn,
		Publisher:   v.Publisher,
		PublishedAt: v.PublishedAt.AsTime(),
	}
}

func BizListToV1BookList(books []*biz.Book) []*commonv1.Book {
	if books == nil {
		return nil
	}
	result := make([]*commonv1.Book, len(books))
	for i, b := range books {
		result[i] = BizToV1Book(b)
	}
	return result
}
