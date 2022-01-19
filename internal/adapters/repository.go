package adapters

import "context"

type Repository interface {
	Save(ctx context.Context, number int) error
}
