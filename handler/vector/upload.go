package vector

import (
	"github.com/gofiber/fiber/v2"
)
import "tinamic/pkg/storage"
import "tinamic/util/response"

func UploadToMinio(ctx *fiber.Ctx) error {

	fileName := ctx.Params("name")
	presignedURL, data, err := storage.Minio.PostPresignedUrl("raster", fileName)
	data["url"] = presignedURL
	if err != nil {
		return response.Fail(ctx, data, "error")
	}
	response.Success(ctx, data, "return upload url")
	return nil
}
