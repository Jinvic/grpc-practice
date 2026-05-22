package bookclient

import (
	bookv1 "bookstore/api/book/v1"
	commonv1 "bookstore/api/common/v1"
	"bookstore/internal/pkg/config"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Client struct {
	bookv1.BookServiceClient
	conn *grpc.ClientConn
}

func NewClient(i do.Injector) (*Client, error) {
	cfg := do.MustInvoke[*config.Config](i)
	serverAddr := fmt.Sprintf("%s:%d", cfg.Services.Book.Host, cfg.Services.Book.Port)

	opts := []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(5*time.Second),
			retry.UnaryClientInterceptor(
				retry.WithMax(3),
				retry.WithPerRetryTimeout(2*time.Second),
			),
		),
	}
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &Client{
		BookServiceClient: bookv1.NewBookServiceClient(conn),
		conn:              conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetBook(ctx context.Context, id int64) (*commonv1.Book, error) {
	resp, err := c.BookServiceClient.GetBook(ctx, &bookv1.GetBookRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return resp.Book, nil
}

func (c *Client) CreateBook(ctx context.Context, book *commonv1.Book) (*commonv1.Book, error) {
	resp, err := c.BookServiceClient.CreateBook(ctx, &bookv1.CreateBookRequest{Book: book})
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}
	return resp.Book, nil
}

func (c *Client) ListBooks(ctx context.Context, page *int32, pageSize *int32) ([]*commonv1.Book, int32, error) {
	resp, err := c.BookServiceClient.ListBooks(ctx, &bookv1.ListBooksRequest{PageNumber: page, PageSize: pageSize})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}
	return resp.Books, resp.TotalCount, nil
}

func (c *Client) UpdateBook(ctx context.Context, book *commonv1.Book, updateMask *fieldmaskpb.FieldMask) (*commonv1.Book, error) {
	resp, err := c.BookServiceClient.UpdateBook(ctx, &bookv1.UpdateBookRequest{Book: book, UpdateMask: updateMask})
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}
	return resp.Book, nil
}

func (c *Client) DeleteBook(ctx context.Context, id int64) error {
	_, err := c.BookServiceClient.DeleteBook(ctx, &bookv1.DeleteBookRequest{Id: id})
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

func (c *Client) BatchImportBooks(ctx context.Context, books []*commonv1.Book) (*bookv1.BatchImportBooksResponse, error) {
	stream, err := c.BookServiceClient.BatchImportBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to import books: %w", err)
	}
	for _, book := range books {
		err := stream.Send(&bookv1.BatchImportBooksRequest{Book: book})
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("failed to close stream: %w", err)
	}
	return resp, nil
}

func (c *Client) BatchExportBooks(ctx context.Context, bookIds []int64) ([]*commonv1.Book, []string, error) {
	req := &bookv1.BatchExportBooksRequest{BookIds: bookIds}
	stream, err := c.BookServiceClient.BatchExportBooks(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to export books: %w", err)
	}

	books := make([]*commonv1.Book, 0)
	errorMessages := make([]string, 0)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to receive response: %w", err)
		}
		if resp.ErrorMessage != "" {
			errorMessages = append(errorMessages, fmt.Sprintf("failed to export book with ID %d: %s", resp.BookId, resp.ErrorMessage))
			continue
		}
		books = append(books, resp.Book)
	}
	return books, errorMessages, nil
}

func (c *Client) HeartBeat(ctx context.Context) error {
	stream, err := c.BookServiceClient.HeartBeat(ctx)
	if err != nil {
		return fmt.Errorf("failed to create heart beat stream: %w", err)
	}

	sendCh := make(chan *bookv1.HeartBeatRequest, 10)
	done := make(chan struct{})
	defer close(done)

	// 发送消息到服务端
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
				sendCh <- &bookv1.HeartBeatRequest{
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

	// 接收服务端的心跳消息
	receiver := func() error {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return status.Errorf(codes.Internal, "failed to receive request: %v", err)
			}

			switch req.Type {
			case commonv1.HeartBeatType_HEART_BEAT_TYPE_PING:
				sendCh <- &bookv1.HeartBeatRequest{
					Type:   commonv1.HeartBeatType_HEART_BEAT_TYPE_PONG,
					SentAt: req.SentAt,
				}
			case commonv1.HeartBeatType_HEART_BEAT_TYPE_PONG:
				receivedAt := time.Now().UnixMilli()
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
