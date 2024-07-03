package middleware

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/model/user"
	"tinamic/util"
	"tinamic/util/response"
)

func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 从请求头中获取 Authorization 字段
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return response.Fail(ctx, "缺少 Authorization")
		}

		// 检查 Bearer token
		tokenString := authHeader[len("Bearer "):]
		claims, err := util.ValidateToken(tokenString)
		if err != nil {
			return response.Fail(ctx, "Authorization 验证失败")
		}

		// 将角色信息存储到 Fiber 上下文中，以便后续处理函数使用
		ctx.Locals(user.GetRoleString(), claims[user.GetRoleString()])

		return ctx.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
