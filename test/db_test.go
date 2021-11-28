package test

import (
	"fmt"
	"testing"

	. "tinamic/database"
)

func TestDbConnect(t *testing.T) {
	db,err:=DbConnect()
	fmt.Println(db,err)
}
