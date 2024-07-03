package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"tinamic/model/datasource"
	"tinamic/service"
	"tinamic/util/response"
)

type DataSourceHandler struct {
	DataSource service.DataSourceService
}

func NewDataSourceHandler(dataSource service.DataSourceService) *DataSourceHandler {
	return &DataSourceHandler{
		DataSource: dataSource,
	}
}

// GetStoreItems / 请求获取文件和文件夹
// /  sourceType, path=当前的文件夹路径
// /
func (ds *DataSourceHandler) GetStoreItems(ctx *fiber.Ctx) error {
	dataSourceType := datasource.DataSourceTypeId(ctx.Params("sourceType"))
	path := ctx.Query("path", "")
	items, err := ds.DataSource.GetStoreItems(dataSourceType, path)
	if err != nil {
		return err
	}
	return response.Success(ctx, items)
}

func (ds *DataSourceHandler) NewFolder(ctx *fiber.Ctx) error {
	path := ctx.Query("path", "")
	fmt.Println(path)
	return nil
}

func (ds *DataSourceHandler) GeneratePutPreSignedUrl(ctx *fiber.Ctx) error {
	dataSourceType := ctx.Params("sourceType")
	path := ctx.Query("path", "")
	fileName := ctx.Query("fileName", "")
	url, err := ds.DataSource.GenPutUploadPresignedUrl(dataSourceType, path, fileName)
	if err != nil {
		return err
	}
	return response.Success(ctx, fiber.Map{
		"uploadUrl": url,
	})
}

func (ds *DataSourceHandler) Upload(ctx *fiber.Ctx) error {
	dataSourceType := ctx.Params("sourceType")
	fileName := ctx.Query("filename")
	if fileName == "" {
		return response.Fail(ctx, "Filename is required")
	}
	path := ctx.Query("path", "")
	// 读取Content-Length头，获取文件大小
	fileSizeStr := ctx.Get("Content-Length")
	if fileSizeStr == "" {
		return response.Fail(ctx, "Content-Length header is required")
	}
	fileSize, _ := strconv.ParseInt(fileSizeStr, 10, 64)
	// 获取请求体
	bodyStream := ctx.Request().BodyStream()
	if bodyStream == nil {
		return response.Fail(ctx, "Request body is empty")
	}

	err := ds.DataSource.UploadToMinio(dataSourceType, path, fileName, bodyStream, fileSize)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	fmt.Println(dataSourceType, path, fileSize)
	return response.Success(ctx, "文件上传成功")
}

//func (ds *DataSourceHandler) ChunkUpload(ctx *fiber.Ctx) error {
//	dataSourceType := ctx.Params("sourceType")
//	filename := ctx.FormValue("filename")
//	if filename == "" {
//		return response.Fail(ctx, "Filename is required")
//	}
//
//	chunkIndexStr := ctx.FormValue("chunkIndex")
//	chunkIndex, err := strconv.Atoi(chunkIndexStr)
//	if err != nil {
//		return response.Fail(ctx, "Invalid chunkIndex")
//	}
//
//	totalChunksStr := ctx.FormValue("totalChunks")
//	totalChunks, err := strconv.Atoi(totalChunksStr)
//	if err != nil {
//		return response.Fail(ctx, "Invalid totalChunks")
//	}
//
//	chunkFile, err := ctx.FormFile("chunk")
//	if err != nil {
//		return response.Fail(ctx, "Failed to get chunk file")
//	}
//
//	file, err := chunkFile.Open()
//	if err != nil {
//		return response.Fail(ctx, "Failed to open chunk file")
//	}
//	defer file.Close()
//	return response.Success(ctx, "文件上传成功")
//}
