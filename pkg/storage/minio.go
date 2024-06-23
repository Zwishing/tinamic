package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"tinamic/conf"
	"tinamic/model"
)

var (
	m    *storage
	once sync.Once
)

// Storage interface that is implemented by storage providers
type storage struct {
	minio *minio.Client
	cfg   *MinioConfig
	ctx   context.Context
	mu    sync.Mutex
}

// NewStorage creates a new storage
func NewStorage(cfg *MinioConfig) (*storage, error) {

	// Minio instance
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Credentials.AccessKey, cfg.Credentials.SecretKey, cfg.Token),
		Secure: cfg.Secure,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	storage := &storage{minio: minioClient, cfg: cfg, ctx: context.Background()}

	// Reset all entries if set to true
	if cfg.Reset {
		if err = storage.Reset(); err != nil {
			return nil, err
		}
	}

	log.Info().Msgf("Connected to minio @ '%s' bucket '%s'", cfg.Endpoint, cfg.Bucket)
	return storage, nil
}

// Get value by key
func (s *storage) Get(key string) ([]byte, error) {

	if len(key) <= 0 {
		return nil, errors.New("the key value is required")
	}

	// get object
	object, err := s.minio.GetObject(s.ctx, s.cfg.Bucket, key, s.cfg.GetObjectOptions)
	if err != nil {
		return nil, err
	}

	// convert to byte
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, err = bb.ReadFrom(object)
	if err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

// Set key with value
func (s *storage) Set(key string, val []byte, exp time.Duration) error {

	if len(key) <= 0 {
		return errors.New("the key value is required")
	}

	// create Reader
	file := bytes.NewReader(val)

	// set content type
	s.mu.Lock()
	s.cfg.PutObjectOptions.ContentType = http.DetectContentType(val)

	// put object
	_, err := s.minio.PutObject(s.ctx, s.cfg.Bucket, key, file, file.Size(), s.cfg.PutObjectOptions)
	s.mu.Unlock()

	return err
}

// Delete entry by key
func (s *storage) Delete(key string) error {

	if len(key) <= 0 {
		return errors.New("the key value is required")
	}

	// remove
	err := s.minio.RemoveObject(s.ctx, s.cfg.Bucket, key, s.cfg.RemoveObjectOptions)

	return err
}

// Reset all entries, including unexpired
func (s *storage) Reset() error {

	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range s.minio.ListObjects(s.ctx, s.cfg.Bucket, s.cfg.ListObjectsOptions) {
			if object.Err != nil {
				log.Error().Msgf("object %s", object.Err)
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for err := range s.minio.RemoveObjects(s.ctx, s.cfg.Bucket, objectsCh, opts) {
		log.Error().Msgf("Error detected during deletion: %s", err)
	}

	return nil
}

// Close the storage
func (s *storage) Close() error {
	return nil
}

// CheckBucket Check to see if bucket already exists
func (s *storage) CheckBucket() error {
	exists, err := s.minio.BucketExists(s.ctx, s.cfg.Bucket)
	if !exists || err != nil {
		return errors.New("the specified bucket does not exist")
	}
	return nil
}

// CreateBucket Bucket not found so Make a new bucket
func (s *storage) CreateBucket() error {
	return s.minio.MakeBucket(s.ctx, s.cfg.Bucket, minio.MakeBucketOptions{Region: s.cfg.Region})
}

// RemoveBucket Bucket remove if bucket is empty
func (s *storage) RemoveBucket() error {
	return s.minio.RemoveBucket(s.ctx, s.cfg.Bucket)
}

// Conn Return minio client
func (s *storage) Conn() *minio.Client {
	return s.minio
}

func (s *storage) Upload(bucketName, objectName string, reader io.Reader,
	objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	object, err := s.minio.PutObject(s.ctx, bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return object, nil
}
func (s *storage) PostPresignedUrl(bucketName, objectName string, expirySecond time.Duration) (string, map[string]string, error) {
	expiry := time.Second * expirySecond

	policy := minio.NewPostPolicy()
	policy.SetBucket(bucketName)
	policy.SetKey(objectName)
	policy.SetExpires(time.Now().UTC().Add(expiry))

	presignedURL, formData, err := s.minio.PresignedPostPolicy(s.ctx, policy)
	if err != nil {
		log.Error().Msgf("%s", err)
		return "", map[string]string{}, err
	}

	return presignedURL.String(), formData, nil
}

func (s *storage) PutPresignedUrl(bucketName, objectName string, expirySecond time.Duration) (string, error) {
	expiry := time.Second * expirySecond

	presignedURL, err := s.minio.PresignedPutObject(s.ctx, bucketName, objectName, expiry)
	if err != nil {
		log.Fatal().Msgf("%s", err)
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *storage) ObjExist(bucketName, objectName string) bool {
	_, err := s.minio.StatObject(s.ctx, bucketName, objectName, s.cfg.GetObjectOptions)
	if err != nil {
		return false
	}
	return true
}
func (s *storage) GetFiles(bucketName string) ([]byte, error) {
	opts := minio.ListObjectsOptions{
		Recursive: true,
	}
	tree := model.FolderFileTree{
		Root: &model.FolderNode{Key: bucketName},
	}
	ctx := context.Background()
	for object := range s.minio.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		if strings.Contains(object.Key, "/") {
			keyParts := strings.Split(object.Key, "/")
			for index, key := range keyParts {
				node := model.FindNode(tree.Root, key)
				if node == nil {
					var child model.Node
					if index == len(keyParts)-1 {
						child = &model.FileNode{
							Key:   object.ETag,
							Type:  "file",
							Title: key,
							Size:  object.Size,
						}
					} else {
						child = &model.FolderNode{
							Key:   key,
							Type:  "folder",
							Title: key,
						}
					}
					if index == 0 {
						tree.AddNode(bucketName, child)
					} else {
						tree.AddNode(keyParts[index-1], child)
					}
				}
			}

		} else {
			child := &model.FileNode{
				Title:        object.Key,
				Key:          object.ETag,
				Type:         "file",
				Size:         object.Size,
				ModifiedTime: object.LastModified,
			}
			tree.AddNode(bucketName, child)
		}
	}
	jsonData, err := json.MarshalIndent(tree.Root.Children, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (s *storage) GetFileByFolder(bucketName, folderPrefix string) {
	ctx := context.Background()
	opts := minio.ListObjectsOptions{
		Recursive: false,
		Prefix:    folderPrefix,
	}
	for object := range s.minio.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
		}
		if object.Size == 0 && object.Key[len(object.Key)-1] == '/' {
			fmt.Printf("Found folder: %s\n", object.Key)
		} else {
			fmt.Printf("Found file: %s\n", object.Key)
		}
	}

}

func GetMinioInstance() model.Client {
	once.Do(func() {
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
