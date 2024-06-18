package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"time"
	"tinamic/model"
	"tinamic/model/user"
	"tinamic/pkg/cache"
	"tinamic/pkg/database"
	"tinamic/service"
	"tinamic/util/response"
)

type UserHandler struct {
}

func Login(ctx *fiber.Ctx) error {
	//获取参数
	//signin := new(user.SignIn)
	signin := &user.SignIn{
		Category: 1,
	}
	err := ctx.BodyParser(signin)
	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}

	// 验证登录
	profile, token := service.Login(signin)

	//token写入redis
	err = cache.RedisClient.SetMap(profile["userId"], profile, time.Second*10000)

	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}
	err = response.Success(ctx, "登录成功")
	// 设置在返回头中
	ctx.Set("Authorization", "Bearer "+token)

	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}
	return nil
}

func Register(ctx *fiber.Ctx) error {

	var user user.User
	err := ctx.BodyParser(&user)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	//v4, err := uuid.NewV4()
	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	user.Id = 1
	user.Created = time.Now()
	user.Edited = time.Now()
	_, err = model.InsertUser(database.GetDbPoolInstance().Pool, user)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	return response.Success(ctx, "")
}

func Profile(ctx *fiber.Ctx) error {
	userId := ctx.Queries()["userId"]
	profile, err := cache.RedisClient.Get(userId)
	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}
	err = response.Success(ctx, fiber.Map{"profile": profile})
	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}
	return nil
}
