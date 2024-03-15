package quote

import "context"

type Storage interface {
	FindOne(ctx context.Context, id string) (Quote, error)
}