package service

import (
	bookv1 "bookstore/api/book/v1"
	commonv1 "bookstore/api/common/v1"
	"bookstore/internal/server/book/usecase"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/samber/do/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookService struct {
	bookv1.UnimplementedBookServiceServer
	bu *usecase.BookUsecase
}

func NewBookService(i do.Injector) (*BookService, error) {
	bu := do.MustInvoke[*usecase.BookUsecase](i)
	return &BookService{bu: bu}, nil
}

func (s *BookService) GetBook(ctx context.Context, req *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	bizBook, err := s.bu.GetBook(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get book: %v", err)
	}

	return &bookv1.GetBookResponse{Book: BizToV1Book(bizBook)}, nil
}

func (s *BookService) CreateBook(ctx context.Context, req *bookv1.CreateBookRequest) (*bookv1.CreateBookResponse, error) {
	bizBook := V1ToBizBook(req.Book)
	bizBook, err := s.bu.CreateBook(ctx, bizBook)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create book: %v", err)
	}

	return &bookv1.CreateBookResponse{Book: BizToV1Book(bizBook)}, nil
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

	books := BizListToV1BookList(bizBooks)

	return &bookv1.ListBooksResponse{
		Books:      books,
		TotalCount: int32(len(bizBooks)),
		PageNumber: int32(pageNumber),
		PageSize:   int32(pageSize),
	}, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req *bookv1.UpdateBookRequest) (*bookv1.UpdateBookResponse, error) {
	bizBook := V1ToBizBook(req.Book)
	bizBook, err := s.bu.UpdateBook(ctx, bizBook)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update book: %v", err)
	}

	return &bookv1.UpdateBookResponse{Book: BizToV1Book(bizBook)}, nil
}

func (s *BookService) DeleteBook(ctx context.Context, req *bookv1.DeleteBookRequest) (*bookv1.DeleteBookResponse, error) {
	err := s.bu.DeleteBook(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete book: %v", err)
	}

	return &bookv1.DeleteBookResponse{Id: req.Id}, nil
}

func (s *BookService) BatchImportBooks(stream bookv1.BookService_BatchImportBooksServer) error {
	var total int32
	var success int32
	var errorMessages []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive request: %v", err)
		}
		bizBook := V1ToBizBook(req.Book)
		_, err = s.bu.CreateBook(stream.Context(), bizBook)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("failed to create book %s: %v", bizBook.Title, err))
		} else {
			success++
		}
		total++

	}
	return stream.SendAndClose(&bookv1.BatchImportBooksResponse{
		Total:         total,
		Success:       success,
		ErrorMessages: errorMessages,
	})
}

func (s *BookService) BatchExportBooks(req *bookv1.BatchExportBooksRequest, stream bookv1.BookService_BatchExportBooksServer) error {
	for _, bookId := range req.BookIds {
		book, err := s.bu.GetBook(stream.Context(), bookId)
		if err != nil {
			if sendErr := stream.Send(&bookv1.BatchExportBooksResponse{
				BookId:       bookId,
				ErrorMessage: fmt.Sprintf("failed to get book with ID %d: %v", bookId, err),
			}); sendErr != nil {
				return status.Errorf(codes.Internal, "failed to send response: %v", sendErr)
			}
			continue
		}
		if err := stream.Send(&bookv1.BatchExportBooksResponse{Book: BizToV1Book(book)}); err != nil {
			return status.Errorf(codes.Internal, "failed to send response: %v", err)
		}
	}
	return nil
}

func (s *BookService) HeartBeat(stream bookv1.BookService_HeartBeatServer) error {
	sendCh := make(chan *bookv1.HeartBeatResponse, 10)
	done := make(chan struct{})
	defer close(done)

	// 发送消息到客户端
	sender := func() {
		for {
			select {
			case msg, ok := <-sendCh:
				if !ok {
					return
				}
				if err := stream.Send(msg); err != nil {
					log.Printf("failed to send message: %v", err)
					return
				}
				log.Println("sent message:", msg.Type)
			case <-done:
				log.Println("sender goroutine stopped")
				return
			case <-stream.Context().Done():
				log.Println("client disconnected, stopping sender goroutine")
				return
			}
		}
	}

	// 定时发送心跳消息
	ticker := func() {
		jitter := time.Duration(rand.Intn(10)) * time.Second
		ticker := time.NewTicker(1*time.Minute + jitter)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sendCh <- &bookv1.HeartBeatResponse{
					Type:   commonv1.HeartBeatType_HEART_BEAT_TYPE_PING,
					SentAt: time.Now().UnixMilli(),
				}
			case <-done:
				log.Println("ticker goroutine stopped")
				return
			case <-stream.Context().Done():
				log.Println("client disconnected, stopping ticker goroutine")
				return
			}
		}
	}

	// 接收客户端的心跳消息
	receiver := func() error {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return status.Errorf(codes.Internal, "failed to receive request: %v", err)
			}

			receivedAt := time.Now().UnixMilli()
			switch req.Type {
			case commonv1.HeartBeatType_HEART_BEAT_TYPE_PING:
				sendCh <- &bookv1.HeartBeatResponse{
					Type:       commonv1.HeartBeatType_HEART_BEAT_TYPE_PONG,
					SentAt:     req.SentAt,
					ReceivedAt: receivedAt,
				}
			case commonv1.HeartBeatType_HEART_BEAT_TYPE_PONG:
				log.Println("received pong, latency:", receivedAt-req.SentAt)
			default:
				log.Println("received unknown heart beat type")
				return status.Errorf(codes.InvalidArgument, "unknown heart beat type: %v", req.Type)
			}
		}
		return nil
	}

	go ticker()
	go sender()
	return receiver()
}
