package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/util/response"
)

func Protected(ctx *fiber.Ctx) error {
	// 从请求头中获取 Authorization 字段
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return response.Fail(ctx, "缺少Authorization")
	}

	// 检查 Bearer token
	tokenString := authHeader[len("Bearer "):]
	err := validateToken(tokenString)
	if err != nil {
		return response.Fail(ctx, "Authorization 验证失败")
	}
	if err != nil {
		return response.Fail(ctx, "Authorization 不正确")
	}

	return ctx.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
