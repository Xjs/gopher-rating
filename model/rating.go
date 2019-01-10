package model

const (
	// NoRating means no rating has been given
	NoRating Rating = iota
	// Rating1 is 1 star, i. e. pretty bad
	Rating1
	// Rating2 is 2 stars, i. e. mediocre
	Rating2
	// Rating3 is 3 stars, i. e. okay
	Rating3
	// Rating4 is 4 stars, i. e. great
	Rating4
	// Rating5 is 5 stars, i. e. excellent
	Rating5
)

// Rating is a 1-5 star rating of something
type Rating int
