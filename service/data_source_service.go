package service

import (
	"github.com/pkg/errors"
	"io"
	"tinamic/model"
	"tinamic/model/datasource"
	"tinamic/repository"
)

type DataSourceService interface {
	GetStoreItems(dst datasource.DataSourceType, path string) ([]model.StoreNode, error)
	GenPutUploadPresignedUrl(bucketName, path, fileName string) (string, error)
	UploadToMinio(bucketName, path, fileName string, reader io.Reader, fileSize int64) error
}

type DataSourceServiceImpl struct {
	repo repository.DataSourceRepository
}

func (ds *DataSourceServiceImpl) UploadToMinio(bucketName, path, fileName string,
	reader io.Reader, fileSize int64) error {
	err := ds.repo.UploadToMinio(bucketName, path+fileName, reader, fileSize)
	if err != nil {
		return err
	}
	return nil
}

func NewDataSourceService(repo repository.DataSourceRepository) DataSourceService {
	return &DataSourceServiceImpl{
		repo: repo,
	}
}

func (ds *DataSourceServiceImpl) GenPutUploadPresignedUrl(bucketName, path, fileName string) (string, error) {
	url, err := ds.repo.GenPutUploadPresignedUrl(bucketName, path, fileName)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (ds *DataSourceServiceImpl) GetStoreItems(dst datasource.DataSourceType, path string) ([]model.StoreNode, error) {
	bucketName := dst.String()
	if !ds.repo.CheckBucket(bucketName) {
		// 不存在时
		return make([]model.StoreNode, 0), errors.Errorf("存储桶%s不存在", bucketName)
	}
	storeItems := ds.repo.GetStoreItems(bucketName, path)
	return storeItems, nil
}
