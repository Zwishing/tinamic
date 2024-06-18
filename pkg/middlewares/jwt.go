package middlewares

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
	"tinamic/model/user"
)

var jwtKey = []byte("a_secret_crect")

func ReleaseToken(user *user.User) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(user.Id),
		"exp":    expirationTime,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("User ID: %s\n", claims["user_id"])
		return nil
	} else {
		return fmt.Errorf("invalid token")
	}
}
