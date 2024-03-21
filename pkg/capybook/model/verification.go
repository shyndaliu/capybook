package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/shyndaliu/capybook/pkg/capybook/validator"
)

type VerificationModel struct {
	DB *sql.DB
}
type Verification struct {
	Code      []byte    `json:"-"`
	PlainText string    `json:"-"`
	UserID    int64     `json:"user_id"`
	Expiry    time.Time `json:"expiry"`
}

func generateVerificationCode(userID int64, ttl time.Duration) (*Verification, error) {
	verification := &Verification{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	verification.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(verification.PlainText))
	verification.Code = hash[:]
	return verification, nil

}
func (v VerificationModel) New(userId int64, ttl time.Duration) (*Verification, error) {
	newVer, err := generateVerificationCode(userId, ttl)
	if err != nil {
		return nil, err
	}

	err = v.Insert(newVer)
	if err != nil {
		return nil, err
	}
	return newVer, nil

}
func (v VerificationModel) Insert(ver *Verification) error {
	query := `
	INSERT INTO verifications (code, user_id, expiry)
	VALUES ($1, $2, $3)`
	args := []interface{}{ver.Code, ver.UserID, ver.Expiry}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := v.DB.ExecContext(ctx, query, args...)
	return err

}
func (v VerificationModel) Delete(userID int64) error {
	query := `
	DELETE FROM verifications
	WHERE user_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := v.DB.ExecContext(ctx, query, userID)

	return err
}

func ValidateVerificationCode(v *validator.Validator, plainTextCode string) {
	v.Check(plainTextCode != "", "code", "must be provided")
	v.Check(len(plainTextCode) == 26, "code", "must be 26 bytes long")
}
