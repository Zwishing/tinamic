package vector

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/pkg/storage"
	"tinamic/util/response"
)

func UploadToMinio(ctx *fiber.Ctx) error {
	fileName := ctx.Params("fileName")
	//category := ctx.Params("category")
	presignedURL, data, err := storage.Minio.PostPresignedUrl("raster", fileName)
	data["url"] = presignedURL
	if err != nil {
		return response.Fail(ctx, data, "error")
	}
	response.Success(ctx, data, "return upload url")
	return nil
}

func AddVectorSource(ctx *fiber.Ctx) error {

	// 解析上传到minio的数据

	// 添加到数据库

	return nil
}
