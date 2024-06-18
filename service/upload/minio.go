package upload

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
	"tinamic/pkg/storage"
)

func PutUploadPresignedUrl(bucketName, fileName string, expirySecond time.Duration) (string, error) {
	minio := storage.GetMinioInstance()
	if minio.ObjExist(bucketName, fileName) {
		err := errors.New(fmt.Sprintf("%s/%s exist!", bucketName, fileName))
		log.Error().Msgf(err.Error())
		return "", err

	}
	presignedURL, err := minio.PutPresignedUrl(bucketName, fileName, expirySecond)
	if err != nil {
		return "", err
	}
	return presignedURL, nil
}

func PostUploadPresignedUrl(bucketName, fileName string, expirySecond time.Duration) (string, map[string]string, error) {
	minio := storage.GetMinioInstance()
	if minio.ObjExist(bucketName, fileName) {
		err := errors.New(fmt.Sprintf("%s/%s exist!", bucketName, fileName))
		log.Error().Msgf(err.Error())
		return "", nil, err

	}
	presignedURL, formData, err := minio.PostPresignedUrl(bucketName, fileName, expirySecond)
	if err != nil {
		return "", nil, err
	}
	return presignedURL, formData, nil
}
