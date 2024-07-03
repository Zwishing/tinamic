package repository

import (
	"sync"
	"tinamic/conf"
	. "tinamic/pkg/storage"
)

var (
	m         *Storage
	minioOnce sync.Once
)

func GetMinioInstance() *Storage {
	minioOnce.Do(func() {
		cfg := conf.GetConfigInstance()
		minioConfig := NewMinioConfig(
			WithBucket(cfg.GetString("storage.minio.bucket")),
			WithEndpoint(cfg.GetString("storage.minio.endpoint")),
			WithRegion(cfg.GetString("storage.minio.region")),
			WithToken(cfg.GetString("storage.minio.token")),
			WithSecure(cfg.GetBool("storage.minio.secure")),
			WithCredentials(
				cfg.GetString("storage.minio.accessKey"),
				cfg.GetString("storage.minio.secretKey")),
		)
		m, _ = NewStorage(minioConfig)
	})
	return m
}
