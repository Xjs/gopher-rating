package model

import (
	"crypto/sha256"
)

// A Hash is a SHA256 hash
type Hash [sha256.Size]byte

// Rehash rehashes the Gopher
func (gopher *Gopher) Rehash() {
	gopher.Hash = sha256.Sum256(gopher.Raw)
}

// NewGopher creates a new Gopher from bytes
func NewGopher(raw []byte) *Gopher {
	gopher := &Gopher{Raw: raw}
	gopher.Rehash()
	return gopher
}

// A Gopher describes a picture of a Gopher, together with the SHA256 digest of the content
type Gopher struct {
	Hash Hash
	Raw  []byte
}

// TODO Machine learning to detect whether pictures are gophers
