package user

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/service/user"
)

type UserController struct {
}

func Login(ctx *fiber.Ctx) error {

	//获取参数
	//var signin model.SignIn
	//
	//err := ctx.BodyParser(&signin)
	user.ValidateUser("admin", "1234")
	//if err != nil {
	//	return response.Fail(ctx, "", err.Error())
	//}
	//
	//user, err := model.QueryUser(database.Db, signin.Name)
	//if err != nil {
	//	return response.Fail(ctx, "", err.Error())
	//}
	//
	//token, err := middlewares.ReleaseToken(user)
	//if err != nil {
	//	return response.Fail(ctx, "", err.Error())
	//}
	//
	//err = response.Success(ctx, fiber.Map{"token": token}, "返回成功")
	//if err != nil {
	//	return response.Fail(ctx, "", err.Error())
	//}
	return nil
}

//func Register(ctx *fiber.Ctx) error {
//
//	//var user model.User
//	//err := ctx.BodyParser(&user)
//	//if err != nil {
//	//	return response.Fail(ctx, "", err.Error())
//	//}
//	////v4, err := uuid.NewV4()
//	//if err != nil {
//	//	return response.Fail(ctx, "", err.Error())
//	//}
//	//user.Id = 1
//	//user.Created = time.Now()
//	//user.Edited = time.Now()
//	//tag, err := model.InsertUser(database.Db, user)
//	//if err != nil {
//	//	return response.Fail(ctx, "", err.Error())
//	//}
//	return response.Success(ctx, "", string(tag))
//}
