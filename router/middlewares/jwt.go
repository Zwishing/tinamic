package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
	"tinamic/model/user"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId int
	jwt.RegisteredClaims
}

func ReleaseToken(user user.User) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := Claims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tinamic.tech",
			Subject:   "user token",
			Audience:  jwt.ClaimStrings{""},
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			NotBefore: &jwt.NumericDate{Time: expirationTime},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        strconv.Itoa(user.Id),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(app *fiber.App, token string) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtKey,
	}))
}
