package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lucaserm/go-sinmos/internal/env"
)

func CreateToken(userId string) (string, error) {
	secretKey := env.GetString("JWT_SECRET", "somerandomsecret")
	expInSeconds := env.GetInt("JWT_EXP", 60*60*24*7) // a week
	exp := time.Second * time.Duration(expInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(exp).Unix(),
	})
	key := []byte(secretKey)
	return token.SignedString(key)
}

func ValidateToken(payloadToken string, withClaimValidation bool) (string, error) {
	var (
		secretKey = []byte(env.GetString("JWT_SECRET", "somerandomsecret"))
		claims    = jwt.MapClaims{}
		token     *jwt.Token
		err       error
	)

	if withClaimValidation {
		token, err = jwt.ParseWithClaims(payloadToken, claims, func(t *jwt.Token) (any, error) {
			return secretKey, nil
		})
	} else {
		token, err = jwt.ParseWithClaims(payloadToken, claims, func(t *jwt.Token) (any, error) {
			return secretKey, nil
		}, jwt.WithoutClaimsValidation())
	}

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	userId := claims["userId"].(string)

	return userId, nil
}
