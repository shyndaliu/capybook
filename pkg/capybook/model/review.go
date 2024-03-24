package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shyndaliu/capybook/pkg/capybook/validator"
)

type ReviewModel struct {
	DB *sql.DB
}
type Review struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"posted_at"`
	AuthorId       int64     `json:"-"`
	AuthorUsername string    `json:"author"`
	BookId         int64     `json:"-"`
	BookTitle      string    `json:"book"`
	Content        string    `json:"content"`
	Rating         int       `json:"rating"`
}

func (r ReviewModel) Insert(review *Review) error {
	query := `
	INSERT INTO reviews (user_id,book_id, content, rating)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`
	args := []interface{}{review.AuthorId, review.BookId, review.Content, review.Rating}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.CreatedAt)
	if err != nil {
		return err
	}
	return nil

}
func (r ReviewModel) Get(book_id int64, user_id int64) (*Review, error) {
	query := `
	select reviews.id, created_at,username, title, content, rating from reviews
    join books on book_id=books.id
    join users on user_id=users.id
	where book_id=$1 and user_id=$2
	limit 1;`
	var review Review
	err := r.DB.QueryRow(query, book_id, user_id).Scan(

		&review.ID,
		&review.CreatedAt,
		&review.AuthorUsername,
		&review.BookTitle,
		&review.Content,
		&review.Rating,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &review, nil
}
func (r ReviewModel) GetAll(book_id int64, filters Filters) ([]*Review, error) {
	query := fmt.Sprintf(`
	select reviews.id, created_at,username, title, content, rating from reviews
    join books on book_id=books.id
    join users on user_id=users.id
	where book_id=$1
	order by %s %s, reviews.id asc 
	limit $2 offset $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Pass the title and genres as the placeholder parameter values.
	rows, err := r.DB.QueryContext(ctx, query, book_id, filters.Limit, filters.offset())

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []*Review{}
	for rows.Next() {
		var review Review
		err := rows.Scan(
			&review.ID,
			&review.CreatedAt,
			&review.AuthorUsername,
			&review.BookTitle,
			&review.Content,
			&review.Rating,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r ReviewModel) Update(review *Review) error {
	query := `
	UPDATE reviews
	SET content=$1, rating=$2
	where id=$3
	returning id`
	err := r.DB.QueryRow(query, review.Content, review.Rating, review.ID).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEditConflict
	}
	return nil
}
func (r ReviewModel) Delete(book_id int64, user_id int64) error {
	query := `
		DELETE FROM reviews
		WHERE book_id = $1 and user_id=$2`
	result, err := r.DB.Exec(query, book_id, user_id)
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

func ValidateRating(v *validator.Validator, rating int) {
	v.Check(rating > 0, "raing", "must be an integer between 1 and 5 ")
	v.Check(rating <= 5, "raing", "must be an integer between 1 and 5 ")

}

func ValidateContent(v *validator.Validator, content string) {
	v.Check(len(content) >= 50, "content", "must be at least 50 bytes long")
	v.Check(len(content) <= 1000, "password", "must not be more than 1000 bytes long")
}

func ValidateReview(v *validator.Validator, review *Review) {
	ValidateContent(v, review.Content)
	ValidateRating(v, review.Rating)
}
