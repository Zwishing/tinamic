package response

import "github.com/gofiber/fiber/v2"

// Response Return response.
func Response(ctx *fiber.Ctx, httpStatus int, code int, data interface{}, msg string) error {
	err := ctx.Status(httpStatus).JSON(fiber.Map{"code": code, "data": data, "msg": msg})
	if err != nil {
		return err
	}
	return nil
}

// Success Return status 200 OK.
func Success(ctx *fiber.Ctx, data interface{}) error {
	err := Response(ctx, fiber.StatusOK, 200, data, "返回成功")
	if err != nil {
		return err
	}
	return nil
}

// Fail Return status 404.
func Fail(ctx *fiber.Ctx, msg string) error {
	err := Response(ctx, fiber.StatusNotFound, 404, "", msg)
	if err != nil {
		return err
	}
	return nil
}
