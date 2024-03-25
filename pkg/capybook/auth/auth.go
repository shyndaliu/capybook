package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shyndaliu/capybook/pkg/capybook/model"
)

type AuthService struct {
	signKey string
}

func NewAuthService(key string) *AuthService {
	return &AuthService{signKey: key}
}

type RefreshTokenCustomClaims struct {
	Username  string
	CustomKey string
	KeyType   string
	jwt.RegisteredClaims
}

type AccessTokenCustomClaims struct {
	Username string
	KeyType  string
	jwt.RegisteredClaims
}

func (auth *AuthService) GenerateRefreshToken(user *model.User) (string, error) {

	cusKey := auth.GenerateCustomKey(user.Username, user.TokenHash)

	claims := RefreshTokenCustomClaims{
		user.Username,
		cusKey,
		"refresh",
		jwt.RegisteredClaims{
			Issuer: "capybook.auth.service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString((auth.signKey))
}
func (auth *AuthService) GenerateAccessToken(user *model.User) (string, error) {

	claims := AccessTokenCustomClaims{
		user.Username,
		"access",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "capybook.auth.service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString((auth.signKey))
}

func (auth *AuthService) GenerateCustomKey(userID string, tokenHash string) string {

	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
func (auth *AuthService) GenerateRandomString(n int) string {
	ans := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(make([]byte, n))
	return ans
}
func (auth *AuthService) ValidateAccessToken(tokenString string) (*AccessTokenCustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.signKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.Username == "" || claims.KeyType != "access" {
		return nil, errors.New("invalid token: authentication failed")
	}
	return claims, nil
}

func (auth *AuthService) ValidateRefreshToken(tokenString string) (*RefreshTokenCustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.signKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	if !ok || !token.Valid || claims.Username == "" || claims.KeyType != "refresh" {
		return nil, errors.New("invalid token: authentication failed")
	}
	return claims, nil
}
