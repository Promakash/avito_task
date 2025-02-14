package jwt

import (
	"avito_shop/internal/domain"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user domain.User, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error while casting claims")
	}

	claims["user_id"] = user.ID

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(token, secret string) (domain.UserID, error) {
	var userID domain.UserID

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %w", domain.ErrInvalidAuthToken)
		}

		return []byte(secret), nil
	})
	if err != nil {
		return userID, domain.ErrInvalidAuthToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return userID, domain.ErrInvalidAuthToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("missed user_id claim: %w", domain.ErrInvalidAuthToken)
	}
	userID = domain.UserID(userIDFloat)

	return userID, nil
}
