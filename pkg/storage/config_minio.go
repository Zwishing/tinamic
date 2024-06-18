package storage

import (
	"github.com/minio/minio-go/v7"
)

// MinioConfig defines the config for storage.
type MinioConfig struct {
	// Bucket
	// Default fiber-bucket
	Bucket string

	// Endpoint is a host name or an IP address
	Endpoint string

	// Region Set this value to override region cache
	// Optional
	Region string

	// Token Set this value to provide x-amz-security-token (AWS S3 specific)
	// Optional, Default is false
	Token string

	// Secure If set to true, https is used instead of http.
	// Default is false
	Secure bool

	// Reset clears any existing keys in existing Bucket
	// Optional. Default is false
	Reset bool

	// Credentials Minio access key and Minio secret key.
	// Need to be defined
	Credentials Credentials

	// GetObjectOptions Options for GET requests specifying additional options like encryption, If-Match
	GetObjectOptions minio.GetObjectOptions

	// PutObjectOptions
	// Allows user to set optional custom metadata, content headers, encryption keys and number of threads for multipart upload operation.
	PutObjectOptions minio.PutObjectOptions

	// ListObjectsOptions Options per to list objects
	ListObjectsOptions minio.ListObjectsOptions

	// RemoveObjectOptions Allows user to set options
	RemoveObjectOptions minio.RemoveObjectOptions
}

type Credentials struct {
	// AccessKeyID is like user-id that uniquely identifies your account.
	AccessKey string
	// SecretAccessKey is the password to your account.
	SecretKey string
}

type MinioOptions func(cfg *MinioConfig)

// 初始化值
var defaultMinioConfig = MinioConfig{
	Bucket:              "",
	Endpoint:            "",
	Region:              "",
	Token:               "",
	Secure:              false,
	Reset:               false,
	Credentials:         Credentials{},
	GetObjectOptions:    minio.GetObjectOptions{},
	PutObjectOptions:    minio.PutObjectOptions{},
	ListObjectsOptions:  minio.ListObjectsOptions{},
	RemoveObjectOptions: minio.RemoveObjectOptions{},
}

func NewMinioConfig(opts ...MinioOptions) *MinioConfig {
	// 默认值初始化
	cfg := defaultMinioConfig

	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg
}

// Option functions

func WithBucket(bucket string) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Bucket = bucket
	}
}

func WithEndpoint(endpoint string) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Endpoint = endpoint
	}
}

func WithRegion(region string) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Region = region
	}
}

func WithToken(token string) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Token = token
	}
}

func WithSecure(secure bool) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Secure = secure
	}
}

func WithReset(reset bool) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Reset = reset
	}
}

func WithCredentials(accessKey, secretKey string) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.Credentials = Credentials{
			AccessKey: accessKey,
			SecretKey: secretKey,
		}
	}
}

func WithGetObjectOptions(options minio.GetObjectOptions) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.GetObjectOptions = options
	}
}

func WithPutObjectOptions(options minio.PutObjectOptions) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.PutObjectOptions = options
	}
}

func WithListObjectsOptions(options minio.ListObjectsOptions) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.ListObjectsOptions = options
	}
}

func WithRemoveObjectOptions(options minio.RemoveObjectOptions) MinioOptions {
	return func(cfg *MinioConfig) {
		cfg.RemoveObjectOptions = options
	}
}
