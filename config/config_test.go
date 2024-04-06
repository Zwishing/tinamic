package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	New()
}

func TestConfig_GetPgConfig(t *testing.T) {
	//Conf = New()
	Conf.GetPgConfig()
}

func TestConfig_GetPgConnString(t *testing.T) {
	//Conf = New()
	connString := Conf.GetPgConnString()
	fmt.Println(connString)
}
