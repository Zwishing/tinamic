package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"time"
	"tinamic/model"
	user "tinamic/model/user"
	"tinamic/pkg/database"
	"tinamic/router/middlewares"
	"tinamic/util/response"
)

type UserController struct {
}

func Login(ctx *fiber.Ctx) error {

	//获取参数
	var signin user.SignIn
	err := ctx.BodyParser(&signin)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}

	user, err := model.QueryUser(database.Db, signin.Name)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}

	token, err := middlewares.ReleaseToken(user)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}

	err = response.Success(ctx, fiber.Map{"token": token}, "返回成功")
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}
	return nil
}

func Register(ctx *fiber.Ctx) error {

	var user user.User
	err := ctx.BodyParser(&user)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}
	v4, err := uuid.NewV4()
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}
	user.UID = v4
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	tag, err := model.InsertUser(database.Db, user)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}
	return response.Success(ctx, "", string(tag))
}
