package storage

import "github.com/Xjs/gopher-rating/model"

// Interface is a storage interface for Gopher pictures
type Interface interface {
	Save(*model.Gopher) error
	Load(model.Hash) (*model.Gopher, error)
	List(start, count int) ([]model.Hash, error)
	Count() (int, error)
}
