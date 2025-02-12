package jwt

import (
	"avito_shop/internal/domain"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidAuthToken = errors.New("invalid auth token")

func NewToken(user domain.User, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(token, secret string) (any, error) {
	var userID domain.UserID

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %w", ErrInvalidAuthToken)
		}

		return []byte(secret), nil
	})
	if err != nil {
		return userID, ErrInvalidAuthToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return userID, ErrInvalidAuthToken
	}

	return claims["user_id"], nil
}
