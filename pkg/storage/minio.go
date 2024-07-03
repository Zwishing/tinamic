package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"tinamic/model"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	minio *minio.Client
	cfg   *MinioConfig
	ctx   context.Context
	mu    sync.Mutex
}

// NewStorage creates a new storage
func NewStorage(cfg *MinioConfig) (*Storage, error) {

	// Minio instance
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Credentials.AccessKey, cfg.Credentials.SecretKey, cfg.Token),
		Secure: cfg.Secure,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	storage := &Storage{minio: minioClient, cfg: cfg, ctx: context.Background()}

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
func (s *Storage) Get(key string) ([]byte, error) {

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
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {

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
func (s *Storage) Delete(key string) error {

	if len(key) <= 0 {
		return errors.New("the key value is required")
	}

	// remove
	err := s.minio.RemoveObject(s.ctx, s.cfg.Bucket, key, s.cfg.RemoveObjectOptions)

	return err
}

// Reset all entries, including unexpired
func (s *Storage) Reset() error {

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
func (s *Storage) Close() error {
	return nil
}

// CheckBucket Check to see if bucket already exists
func (s *Storage) CheckBucket(bucketName string) error {
	exists, err := s.minio.BucketExists(s.ctx, bucketName)
	if !exists || err != nil {
		return errors.New("the specified bucket does not exist")
	}
	return nil
}

// CreateBucket Bucket not found so Make a new bucket
func (s *Storage) CreateBucket() error {
	return s.minio.MakeBucket(s.ctx, s.cfg.Bucket, minio.MakeBucketOptions{Region: s.cfg.Region})
}

// RemoveBucket Bucket remove if bucket is empty
func (s *Storage) RemoveBucket() error {
	return s.minio.RemoveBucket(s.ctx, s.cfg.Bucket)
}

// Conn Return minio client
func (s *Storage) Conn() *minio.Client {
	return s.minio
}

func (s *Storage) Upload(bucketName, objectName string, reader io.Reader,
	objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	object, err := s.minio.PutObject(s.ctx, bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return object, nil
}

func (s *Storage) PostPresignedUrl(bucketName, objectName string, expirySecond time.Duration) (string, map[string]string, error) {
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

func (s *Storage) PutPresignedUrl(bucketName, objectName string, expirySecond time.Duration) (string, error) {
	expiry := time.Second * expirySecond

	presignedURL, err := s.minio.PresignedPutObject(s.ctx, bucketName, objectName, expiry)
	if err != nil {
		log.Fatal().Msgf("%s", err)
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *Storage) ObjExist(bucketName, objectName string) bool {
	_, err := s.minio.StatObject(s.ctx, bucketName, objectName, s.cfg.GetObjectOptions)
	if err != nil {
		return false
	}
	return true
}

func (s *Storage) GetFiles(bucketName string) ([]byte, error) {
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
					var child model.StoreNode
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

func (s *Storage) GetStoreObjectByPath(bucketName, folderPrefix string) []minio.ObjectInfo {
	var storeObjects []minio.ObjectInfo
	ctx := context.Background()
	opts := minio.ListObjectsOptions{
		Recursive: false,
		Prefix:    folderPrefix,
	}
	for object := range s.minio.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
			return nil
		}
		storeObjects = append(storeObjects, object)
	}
	return storeObjects
}
