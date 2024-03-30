package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"time"
	"tinamic/app/middlewares"
	"tinamic/app/models"
	"tinamic/app/queries"
	"tinamic/common/database"
	"tinamic/common/response"
)

type UserController struct {
}

func Login(ctx *fiber.Ctx) error {

	//获取参数
	var signin models.SignIn
	err := ctx.BodyParser(&signin)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}

	user, err := queries.QueryUser(database.Db, signin.Name)
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

	var user models.User
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
	tag, err := queries.InsertUser(database.Db, user)
	if err != nil {
		return response.Fail(ctx, "", err.Error())
	}
	return response.Success(ctx, "", string(tag))
}
