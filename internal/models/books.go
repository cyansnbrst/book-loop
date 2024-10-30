package models

import (
	"time"

	"bookloop.net/pkg/utils"
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateBook(v *utils.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")

	v.Check(book.Author != "", "author", "must be provided")

	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(utils.Unique(book.Genres), "genres", "must not contain duplicate values")
}
