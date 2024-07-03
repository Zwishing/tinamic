package service

import (
	"encoding/json"
	uuid2 "github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
	"tinamic/model/user"
	"tinamic/repository"
	"tinamic/util"
)

type UserService interface {
	Register(s *user.RegisterDTO) error
	Login(s *user.SignInDTO) (*user.UserDTO, error)
	GenerateToken(userId string, userRole user.RoleCode) (string, error)
	GetProfile(userId string) (map[string]any, error)
	IsRegistered(userAccount string, category string) bool
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepository: repo,
	}
}

func (us *UserServiceImpl) Register(s *user.RegisterDTO) error {
	code := user.RoleCode(s.Role)
	roleId := user.RoleId(&code)
	// 随机盐
	salt := util.RandomSalt()
	// 加密密码
	uuid, _ := uuid2.NewUUID()
	u := user.NewUser(
		user.WithUserId(uuid.String()),
		user.WithSalt(salt),
		user.WithPassword(util.CreateHashPassword(s.Password, salt)))

	account := user.NewAccount(
		user.WithAccountUserId(uuid.String()),
		user.WithLoginAccount(s.UserAccount),
		user.WithCategory(user.AccountCategoryId(s.Category)),
	)

	userRole := &user.UserRole{
		UserId: uuid.String(),
		RoleId: roleId,
	}
	err := us.UserRepository.AddUser(account, u, userRole)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserServiceImpl) Login(signin *user.SignInDTO) (*user.UserDTO, error) {
	account := user.NewAccount(
		user.WithLoginAccount(signin.UserAccount),
		user.WithCategory(user.AccountCategoryId(signin.Category)))
	u := user.NewUser()
	role := &user.UserRole{}
	err := us.UserRepository.QueryUser(account, u, role)
	if err != nil {
		return nil, err
	}
	// 验证密码是否正确
	if !util.ValidatePassword(signin.Password, u.Salt, u.Password) {
		log.Error().Msgf("password is error,login fail")
		return nil, errors.New("password is error")
	}
	userDTO := &user.UserDTO{
		UserId:      account.UserId,
		UserAccount: account.UserAccount,
		Category:    account.Category.String(),
		Role:        user.RoleId2Code(role.RoleId),
		Name:        u.Name,
		Avatar:      u.Avatar,
		CellPhone:   u.CellPhone,
	}
	err = us.UserRepository.SaveToRedis(userDTO.UserId, userDTO, time.Hour*1)
	if err != nil {
		return nil, err
	}
	return userDTO, nil

}

func (us *UserServiceImpl) GenerateToken(userId string, userRole user.RoleCode) (string, error) {
	//生成带角色信息的token
	token, err := util.ReleaseToken(userId, string(userRole))
	if err != nil {
		log.Error().Msgf(err.Error())
		return "", err
	}
	return token, nil
}

func (us *UserServiceImpl) GetProfile(userId string) (map[string]any, error) {
	profileString, err := us.UserRepository.GetFromRedis(userId)
	if err != nil {
		return nil, err
	}
	var profile map[string]any
	err = json.Unmarshal([]byte(profileString), &profile)
	return profile, nil
}

func (us *UserServiceImpl) IsRegistered(userAccount string, category string) bool {
	account := user.NewAccount(
		user.WithLoginAccount(userAccount),
		user.WithCategory(user.AccountCategoryId(category)),
	)
	err := us.UserRepository.QueryAccount(account)
	if err != nil {
		return false
	}
	return true
}
