package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func RandomSalt() string {
	saltBytes := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, saltBytes)
	if err != nil {
		panic(err)
	}
	salt := base64.StdEncoding.EncodeToString(saltBytes)
	return salt
}

func CreateHashPassword(password string, salt string) string {
	// 加盐处理
	toHash := password + salt
	h := sha256.New()
	h.Write([]byte(toHash))
	hashPassword := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hashPassword
}

func ValidatePassword(password, salt, hashPassword string) bool {
	computedHash := CreateHashPassword(password, salt)
	return hashPassword == computedHash
}
