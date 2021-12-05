package middlewares

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"tinamic/app/models"
)

var jwtKey=[]byte("a_secret_crect")

type Claims struct {
	UserId uuid.UUID
	jwt.RegisteredClaims
}

func ReleaseToken(user models.User) (string, error){

	expirationTime:=time.Now().Add(24*time.Hour)

	claims:=Claims{
		UserId: user.UID,
		RegisteredClaims:jwt.RegisteredClaims{
			Issuer:"tinamic.tech",
			Subject: "user token",
			Audience:jwt.ClaimStrings{""},
			ExpiresAt:&jwt.NumericDate{Time: expirationTime},
			NotBefore: &jwt.NumericDate{Time: expirationTime},
			IssuedAt: &jwt.NumericDate{Time: time.Now()},
			ID:user.UID.String(),
		},
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString,err:=token.SignedString(jwtKey)
	if err!=nil{
		return "",err
	}
	return tokenString,nil
}