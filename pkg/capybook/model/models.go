package model

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Books         BookModel
	Users         UserModel
	Verifications VerificationModel
	Permissions   PermissionModel
	Reviews       ReviewModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books:         BookModel{DB: db},
		Users:         UserModel{DB: db},
		Verifications: VerificationModel{DB: db},
		Permissions:   PermissionModel{DB: db},
		Reviews:       ReviewModel{DB: db},
	}
}
