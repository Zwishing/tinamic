package model

import (
	"time"
)

type Client interface {
	ObjExist(name string, name2 string) bool
	PutPresignedUrl(name string, name2 string, second time.Duration) (string, error)
	PostPresignedUrl(name string, name2 string, second time.Duration) (string, map[string]string, error)
	GetFiles(bucketName string) ([]byte, error)
	GetFileByFolder(bucketName, folderPrefix string)
}
