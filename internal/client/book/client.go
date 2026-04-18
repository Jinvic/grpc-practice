package bookclient

import (
	bookv1 "bookstore/api/book/v1"
	commonv1 "bookstore/api/common/v1"
	"bookstore/internal/pkg/config"
	"context"
	"fmt"

	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
