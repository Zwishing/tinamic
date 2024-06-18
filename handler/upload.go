package handler

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"tinamic/service/upload"
	"tinamic/util/response"
)

type DataType int8

const (
	Vector DataType = iota
	Raster
)

func CreatePostPresignedUrl(ctx *fiber.Ctx) error {

	bucketName, fileName := parseDataTypeUrl(ctx)
	presignedURL, data, err := upload.PostUploadPresignedUrl(bucketName, fileName, 1000)

	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	data["url"] = presignedURL
	err = response.Success(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func CreatePutPresignedUrl(ctx *fiber.Ctx) error {
	bucketName, fileName := parseDataTypeUrl(ctx)

	//dataType := ctx.Queries()["dataType"]

	presignedURL, err := upload.PutUploadPresignedUrl(bucketName, fileName, 1000)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	var data = make(map[string]string)
	data["url"] = presignedURL
	err = response.Success(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func parseDataTypeUrl(ctx *fiber.Ctx) (string, string) {

	dataType, _ := strconv.Atoi(ctx.Params("dtype"))
	fileName := ctx.Queries()["file"]
	var bucketName string
	switch dataType {
	case int(Raster):
		bucketName = "raster"
	case int(Vector):
		bucketName = "vector"
	}
	return bucketName, fileName
}
