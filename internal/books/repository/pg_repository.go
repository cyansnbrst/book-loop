package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bookloop.net/internal/books"
	"bookloop.net/internal/models"
	"bookloop.net/pkg/db"
	"bookloop.net/pkg/utils"
	"github.com/lib/pq"
)

type booksRepo struct {
	db *sql.DB
}

func newBooksRepository(db *sql.DB) books.Repository {
	return &booksRepo{db: db}
}

func (r *booksRepo) GetAll(title string, author string, genres []string, filters utils.Filters) ([]*models.Book, utils.Pagination, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, author, genres, version
		FROM books
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', author) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (genres @> $3 OR $3 = '{}')
		ORDER BY %s %s
		LIMIT $4 OFFSET $5`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		title,
		author,
		pq.Array(genres),
		filters.Limit(),
		filters.Offset(),
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	defer rows.Close()

	totalRecords := 0
	books := []*models.Book{}

	for rows.Next() {
		var book models.Book

		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Author,
			pq.Array(&book.Genres),
			&book.Version,
		)
		if err != nil {
			return nil, utils.Pagination{}, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.Pagination{}, err
	}

	metadata := utils.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return books, metadata, nil
}

func (r *booksRepo) Insert(book *models.Book) error {
	query := `
		INSERT INTO books (title, author, genres)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`
	args := []interface{}{book.Title, book.Author, pq.Array(book.Genres)}

	return r.db.QueryRow(query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (r *booksRepo) Get(id int64) (*models.Book, error) {
	if id < 1 {
		return nil, db.ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, author, genres, version
		FROM books
		WHERE id = $1`

	var book models.Book
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Author,
		pq.Array(&book.Genres),
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, db.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

func (r *booksRepo) Update(book *models.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, genres = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []interface{}{
		book.Title,
		book.Author,
		pq.Array(book.Genres),
		book.ID,
		book.Version,
	}

	err := r.db.QueryRow(query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return db.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *booksRepo) Delete(id int64) error {
	if id < 1 {
		return db.ErrEditConflict
	}

	query := `
		DELETE FROM books
		WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return db.ErrRecordNotFound
	}

	return nil
}
