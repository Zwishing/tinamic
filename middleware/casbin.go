package middleware

import (
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"tinamic/model/user"
	icasbin "tinamic/pkg/casbin"
	"tinamic/repository"
)

func Authz() *casbin.Middleware {
	adapter, _ := icasbin.NewAdapter(repository.GetDbPoolInstance().Pool)
	authz := casbin.New(casbin.Config{
		ModelFilePath: "./conf/model.conf",
		PolicyAdapter: adapter,
		Lookup: func(ctx *fiber.Ctx) string {
			role := ctx.Locals(user.GetRoleString())
			if role != nil {
				return role.(string)
			}
			return ""
		},
	})
	return authz
}
