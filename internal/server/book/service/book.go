package service

import (
	bookv1 "bookstore/api/book/v1"
	commonv1 "bookstore/api/common/v1"
	"bookstore/internal/server/book/model/biz"
	"bookstore/internal/server/book/usecase"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookService struct {
	bookv1.UnimplementedBookServiceServer
	bu *usecase.BookUsecase
}

func NewBookService(bu *usecase.BookUsecase) *BookService {
	return &BookService{bu: bu}
}

func (s *BookService) GetBook(ctx context.Context, req *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	bizBook, err := s.bu.GetBook(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get book: %v", err)
	}

	return &bookv1.GetBookResponse{Book: bizBook.ToV1Book()}, nil
}

func (s *BookService) CreateBook(ctx context.Context, req *bookv1.CreateBookRequest) (*bookv1.CreateBookResponse, error) {
	bizBook := biz.FromV1Book(req.Book)
	bizBook, err := s.bu.CreateBook(ctx, bizBook)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create book: %v", err)
	}

	return &bookv1.CreateBookResponse{Book: bizBook.ToV1Book()}, nil
}

func (s *BookService) ListBooks(ctx context.Context, req *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	var pageNumber int
	var pageSize int
	if req.PageNumber == nil {
		pageNumber = 1
	} else {
		pageNumber = int(*req.PageNumber)
	}
	if req.PageSize == nil {
		pageSize = 10
	} else {
		pageSize = int(*req.PageSize)
	}

	bizBooks, err := s.bu.ListBooks(ctx, pageNumber, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list books: %v", err)
	}

	books := make([]*commonv1.Book, len(bizBooks))
	for i, bizBook := range bizBooks {
		books[i] = bizBook.ToV1Book()
	}

	return &bookv1.ListBooksResponse{
		Books:      books,
		TotalCount: int32(len(bizBooks)),
		PageNumber: int32(pageNumber),
		PageSize:   int32(pageSize),
	}, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req *bookv1.UpdateBookRequest) (*bookv1.UpdateBookResponse, error) {
	bizBook := biz.FromV1Book(req.Book)
	bizBook, err := s.bu.UpdateBook(ctx, bizBook)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update book: %v", err)
	}

	return &bookv1.UpdateBookResponse{Book: bizBook.ToV1Book()}, nil
}

func (s *BookService) DeleteBook(ctx context.Context, req *bookv1.DeleteBookRequest) (*bookv1.DeleteBookResponse, error) {
	err := s.bu.DeleteBook(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete book: %v", err)
	}

	return &bookv1.DeleteBookResponse{Id: req.Id}, nil
}
