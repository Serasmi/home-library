package books

import "context"

type Storage interface {
	Find(ctx context.Context) ([]Book, error)
}
