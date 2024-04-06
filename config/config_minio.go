package config

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
