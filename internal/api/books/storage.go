package books

import "context"

type Storage interface {
	Find(ctx context.Context) ([]Book, error)
	FindOne(ctx context.Context, id string) (Book, error)
	Remove(ctx context.Context, id string) error
}
