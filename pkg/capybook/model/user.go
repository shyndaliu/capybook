package model

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id`
	args := []interface{}{user.Username, user.Password}
	return u.DB.QueryRow(query, args...).Scan(&user.ID)
}
