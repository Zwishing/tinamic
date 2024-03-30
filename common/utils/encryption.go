package utils

import (
	"crypto/rand"
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
