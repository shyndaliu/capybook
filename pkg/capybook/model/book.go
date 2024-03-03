package model

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type BookModel struct {
	DB *sql.DB
}

type Book struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Year        int32    `json:"year"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
}

func (b BookModel) Insert(book *Book) error {
	query := `
	INSERT INTO books (title, author, year, description, genres)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	args := []interface{}{book.Title, book.Author, book.Year, book.Description, pq.Array(book.Genres)}
	return b.DB.QueryRow(query, args...).Scan(&book.ID)
}

func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT * FROM books
	WHERE id = $1`
	var book Book
	err := b.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Year,
		&book.Description,
		pq.Array(&book.Genres),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &book, nil
}

func (b BookModel) Update(book *Book) (*Book, error) {
	query := `
UPDATE books
SET title = $1, author=$2, year = $3, description = $4, genres = $5
WHERE id = $6
RETURNING *`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		book.Title,
		book.Author,
		book.Year,
		book.Description,
		pq.Array(book.Genres),
		book.ID,
	}
	var newbook Book
	err := b.DB.QueryRow(query, args...).Scan(
		&newbook.ID,
		&newbook.Title,
		&newbook.Author,
		&newbook.Year,
		&newbook.Description,
		pq.Array(&newbook.Genres),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &newbook, nil
}

func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM books
		WHERE id = $1`
	result, err := b.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
