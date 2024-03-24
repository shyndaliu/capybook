package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/shyndaliu/capybook/pkg/capybook/validator"
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

func (b BookModel) GetAll(title string, author string, genres []string, filters Filters) ([]*Book, error) {
	query := fmt.Sprintf(`
	SELECT id,  title, author,  year, description, genres
	FROM books
	WHERE (LOWER(title) = LOWER($1) OR $1 = '')
	AND (LOWER(author) = LOWER($2) OR $2 = '')
	AND (genres @> $3 OR $3 = '{}')
	ORDER BY %s %s, id ASC
	LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Pass the title and genres as the placeholder parameter values.
	rows, err := b.DB.QueryContext(ctx, query, title, author, pq.Array(genres), filters.Limit, filters.offset())

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := []*Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.Description,
			pq.Array(&book.Genres),
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
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

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 200, "title", "must not be more than 200 bytes long")
	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 200, "author", "must not be more than 200 bytes long")
	v.Check(book.Description != "", "description", "must be provided")
	v.Check(len(book.Description) <= 500, "description", "must not be more than 500 bytes long")
	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}
