package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"tinamic/model/user"
)

var jwtKey = []byte("a_secret_crect")

func ReleaseToken(userId string, role string) (string, error) {

	expirationTime := time.Now().Add(30 * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":             userId,
		user.GetRoleString(): role,
		"genTime":            time.Now(),
		"exp":                expirationTime,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
