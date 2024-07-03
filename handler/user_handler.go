package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"tinamic/model/user"
	"tinamic/pkg/validate"
	"tinamic/service"
	"tinamic/util/response"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (uh *UserHandler) Login(ctx *fiber.Ctx) error {
	//获取参数
	var signin user.SignInDTO

	err := ctx.BodyParser(&signin)
	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}

	err = validate.ValidateRequestBody(signin)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}

	// 验证登录
	ur, err := uh.UserService.Login(&signin)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}
	// 生成JWT Token
	token, err := uh.UserService.GenerateToken(ur.UserId, ur.Role)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}

	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}
	err = response.Success(ctx, fiber.Map{"userId": ur.UserId})
	// 设置在返回头中
	ctx.Set("Authorization", "Bearer "+token)

	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}
	return nil
}

// Register 用户注册接口
// @Summary 注册一个新用户
// @Description 使用用户名和密码注册一个新用户，只有管理员有权限注册，设置角色，默认是游客角色
// @ID register-user
// @Accept  json
// @Produce  json
// @Param   loginAccount    object  true  "User Registration"
// @Success 200 {object} map[string]string
// @Router /v1/user/register [post]
func (uh *UserHandler) Register(ctx *fiber.Ctx) error {
	var reg user.RegisterDTO

	err := ctx.BodyParser(&reg)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}

	err = validate.ValidateRequestBody(reg)
	if err != nil {
		return response.Fail(ctx, err.Error())
	}

	if uh.UserService.IsRegistered(reg.UserAccount, reg.Category) {
		return response.Success(ctx, "该用户已经注册")
	}
	err = uh.UserService.Register(&reg)
	if err != nil {
		return err
	}

	return response.Success(ctx, "注册成功")
}

func (uh *UserHandler) Profile(ctx *fiber.Ctx) error {
	userId := ctx.Queries()["userId"]
	profile, err := uh.UserService.GetProfile(userId)
	if err != nil {
		log.Error().Msgf(err.Error())
		return response.Fail(ctx, err.Error())
	}

	return response.Success(ctx, profile)
}
