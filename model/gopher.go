package model

import (
	"crypto/sha256"
)

// A Hash is a SHA256 hash
type Hash [sha256.Size]byte

// A Gopher describes a picture of a Gopher, together with the SHA256 digest of the content
type Gopher struct {
	Hash Hash
	Raw  []byte
}
