package storage

import (
	"context"
	"log"
	_ "net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	defaultExpiryTime = time.Second * 24 * 60 * 60 // 1 day

	endpoint        string = "39.101.164.253:9000"
	accessKeyID     string = "admin"
	secretAccessKey string = "admin321"
	useSSL          bool   = false
)

type Client struct {
	cli *minio.Client
}

func NewMinioClient() *Client {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{
		cli: cli,
	}
}

func (c *Client) PostPresignedUrl(ctx context.Context, bucketName, objectName string) (string, map[string]string, error) {
	expiry := defaultExpiryTime

	policy := minio.NewPostPolicy()
	_ = policy.SetBucket(bucketName)
	_ = policy.SetKey(objectName)
	_ = policy.SetExpires(time.Now().UTC().Add(expiry))

	presignedURL, formData, err := c.cli.PresignedPostPolicy(ctx, policy)
	if err != nil {
		log.Fatalln(err)
		return "", map[string]string{}, err
	}

	return presignedURL.String(), formData, nil
}

func (c *Client) PutPresignedUrl(ctx context.Context, bucketName, objectName string) (string, error) {
	expiry := defaultExpiryTime

	presignedURL, err := c.cli.PresignedPutObject(ctx, bucketName, objectName, expiry)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return presignedURL.String(), nil
}