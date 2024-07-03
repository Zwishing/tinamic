package repository

import (
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"io"
	"tinamic/model"
	"tinamic/pkg/pg"
	"tinamic/pkg/storage"
)

type DataSourceRepository interface {
	GetStoreItems(bucketName, currentFolder string) []model.StoreNode
	CheckBucket(bucketName string) bool
	GenPutUploadPresignedUrl(bucketName, path, fileName string) (string, error)
	UploadToMinio(bucketName, objectName string, reader io.Reader, objectSize int64) error
}

type DataSourceRepositoryImpl struct {
	*pg.PGPool
	*storage.Storage
}

func NewDataSourceRepository() DataSourceRepository {
	return &DataSourceRepositoryImpl{
		GetDbPoolInstance(),
		GetMinioInstance(),
	}
}

func (dsr *DataSourceRepositoryImpl) UploadToMinio(bucketName, objectName string, reader io.Reader, objectSize int64) error {
	opts := minio.PutObjectOptions{}
	_, err := dsr.Storage.Upload(bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return err
	}
	return nil
}
func (dsr *DataSourceRepositoryImpl) GenPutUploadPresignedUrl(bucketName, path, fileName string) (string, error) {
	minio := dsr.Storage
	if minio.ObjExist(bucketName, fileName) {
		err := errors.Errorf("%s/%s exist!", bucketName, fileName)
		return "", err
	}
	presignedURL, err := minio.PutPresignedUrl(bucketName, path+fileName, 60)
	if err != nil {
		return "", err
	}
	return presignedURL, nil
}

func (dsr *DataSourceRepositoryImpl) GetStoreItems(bucketName, currentFolder string) []model.StoreNode {
	var storeItems []model.StoreNode
	items := dsr.GetStoreObjectByPath(bucketName, currentFolder)
	for _, item := range items {
		// 文件夹类型
		if item.Size == 0 && item.Key[len(item.Key)-1] == '/' {
			folder := &model.FolderNode{
				Title:        item.Key,
				Key:          item.ETag,
				Type:         "folder",
				Size:         item.Size,
				ModifiedTime: item.LastModified,
				Children:     []model.StoreNode{},
			}
			storeItems = append(storeItems, folder)
		} else {
			file := &model.FileNode{
				Title:        item.Key,
				Key:          item.ETag,
				Type:         "file",
				Size:         item.Size,
				ModifiedTime: item.LastModified,
			}
			storeItems = append(storeItems, file)
		}
	}
	return storeItems

}

func (dsr *DataSourceRepositoryImpl) CheckBucket(bucketName string) bool {
	err := dsr.Storage.CheckBucket(bucketName)
	if err != nil {
		return false
	}
	return true
}
