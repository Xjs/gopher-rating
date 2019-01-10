package storage

import (
	"context"

	"github.com/Xjs/gopher-rating/model"
)

// Interface is a storage interface for Gopher pictures
type Interface interface {
	Save(context.Context, *model.Gopher) error
	Load(context.Context, model.Hash) (*model.Gopher, error)
	List(ctx context.Context, start, count int) ([]model.Hash, error)
	Count(context.Context) (int, error)
	Rating(context.Context, model.Hash) (int, error)
	Rate(context.Context, model.Hash, int) error
}
