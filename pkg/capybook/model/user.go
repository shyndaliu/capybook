package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/shyndaliu/capybook/pkg/capybook/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	Activated bool     `json:"activated"`
}
type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

func (u UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username,email, password)
	VALUES ($1, $2, $3)
	RETURNING id`
	args := []interface{}{user.Username, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, password_hash
	FROM users
	WHERE email = $1`
	var user User
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) GetByVerificationCode(plaintext string) (*User, error) {
	hash := sha256.Sum256([]byte(plaintext))
	query := `
	SELECT users.id, users.username, users.email, users.password, users.activated
	FROM users
	INNER JOIN verifications
	ON users.id = verifications.user_id
	WHERE verifications.code = $1
	AND verifications.expiry > $2`
	args := []interface{}{hash[:], time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil

}

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET username = $1, email = $2, password = $3, activated = $4
	WHERE id = $5 
	RETURNING id`
	args := []interface{}{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
	}
	err := m.DB.QueryRow(query, args...).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEditConflict
	}
	return nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(validator.Matches(user.Username, validator.UsernameRX), "username", "must be a valid username address")

	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
