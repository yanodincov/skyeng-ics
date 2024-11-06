package httphelper

import (
	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
)

var ErrUserIDNotFound = errors.New("userID not found in token claims")

func ParseJWTClaimsInt(tokenString string, fieldName string) (int, error) {
	claims := jwt.MapClaims{}

	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse token")
	}

	res, ok := claims[fieldName].(float64) // JWT library often decodes numbers as float64
	if !ok {
		return 0, ErrUserIDNotFound
	}

	return int(res), nil
}
