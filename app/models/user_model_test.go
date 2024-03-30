package models

import (
	"fmt"
	"testing"
)

func TestNewUser(t *testing.T) {
	user := NewUser("zhang", "123456", "zhang", "zhang")
	fmt.Println(user)
}
