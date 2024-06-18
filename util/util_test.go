package util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestBoxStringToArray(t *testing.T) {
	BoxStringToArray("BOX(111.964430789849 30.9798174714992,116.727509619973 34.8957377544938)")
}

func TestExeCommand(t *testing.T) {
	err := ExeCommand("shp2pgsql")
	if err != nil {
		fmt.Println(err)
	}
}

func TestDeCompress(t *testing.T) {
	err, filename := Decompress("D:\\Code\\go-web\\tinamic\\data\\九段线.zip", "D:\\Code\\go-web\\tinamic\\data")
	if err != nil {
		return
	}
	fmt.Println(filename)
}

func TestGetDir(t *testing.T) {
	p := GetDir("spatialdata/temp.cpg")
	fmt.Println(p)
}

func TestGetFileName(t *testing.T) {
	//name:=GetFileName("spatialdata/aka/temp.cpg")
	//fmt.Println(name)
	dir, file := filepath.Split("D:/Code\\go-web\\tinamic\\uplaod/1212121.zip")
	fmt.Println(filepath.Base(dir), file)
}

func TestGetFileSuffix(t *testing.T) {

}

func TestSearchFile(t *testing.T) {
	fileList, _ := SearchFile("D:\\Code\\go-web\\tinamic\\temp\\data")
	fmt.Println(fileList)
}

func TestVerifyShapeFile(t *testing.T) {
	name := VerifyShapeFile("D:/Code/go-web\\tinamic\\temp\\data\\river.shp")
	fmt.Println(name)
	fi, _ := os.Stat("D:/Code/go-web\\tinamic\\temp\\data\\river.shp")
	fmt.Println(filepath.ToSlash(fi.Name()))
}

func TestRandomSalt(t *testing.T) {
	salt := RandomSalt()
	fmt.Println(salt)
}

func TestCreateHashPassword(t *testing.T) {
	hashPassword := CreateHashPassword("admin123", "")
	fmt.Println(hashPassword)
}

func TestValidatePassword(t *testing.T) {
	isValid := ValidatePassword("admin123", "", "admin123")
	fmt.Println(isValid)

}
