package test

import (
	"fmt"
	"testing"
	"tinamic/config"
	"tinamic/database"
)

var db, _ =database.DbConnect(config.New().GetPgConfig())

func TestDbConnect(t *testing.T) {
	fmt.Println(db)
}
