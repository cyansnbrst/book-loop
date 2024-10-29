package data

import (
	"database/sql"
	"time"

	"bookloop.net/internal/validator"
)

type Advertisment struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	BookID      int64     `json:"book_id"`
	UserID      int64     `json:"user_id"`
	Publisher   string    `json:"publisher,omitempty"`
	State       string    `json:"state"`
	Description string    `json:"description,omitempty"`
}

func ValidateAdvertisment(v *validator.Validator, advertisment *Advertisment) {
	v.Check(advertisment.BookID > 0, "book_id", "must not be zero value")

	v.Check(advertisment.UserID > 0, "user_id", "must not be zero value")

	v.Check(advertisment.State != "", "state", "must be provided")
	v.Check(validator.In(advertisment.State, "новая", "б/у"), "state", "value must be one of the following: \"новая\", \"б/у\"")
}

type AdvertismentModel struct {
	DB *sql.DB
}

func (m AdvertismentModel) Insert(advertisment *Advertisment) error {
	return nil
}

func (m AdvertismentModel) Get(id int64) (*Advertisment, error) {
	return nil, nil
}

func (m AdvertismentModel) Update(advertisment *Advertisment) error {
	return nil
}

func (m AdvertismentModel) Delete(id int64) error {
	return nil
}
