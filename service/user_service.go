package service

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"tinamic/model/user"
	"tinamic/pkg/database"
	"tinamic/pkg/middlewares"
	"tinamic/util"
)

func Login(signin *user.SignIn) (map[string]string, string) {
	database := database.GetDbPoolInstance()
	// 通过Account查询用户id
	var account = new(user.Account)
	accountSql := fmt.Sprintf(
		"SELECT user_id,login_account,category FROM user_info.account WHERE login_account='%s' and category=%d",
		signin.LoginAccount, signin.Category)
	err := database.SelectRow(accountSql, &account.UserId, &account.LoginAccount, &account.Category)
	if err != nil {
		log.Error().Msgf(err.Error())
		return nil, ""
	}

	// 查询user，获取密码
	var usr = new(user.User)
	userSql := fmt.Sprintf("SELECT id,salt,password,name,cell_phone FROM user_info.user WHERE id=%d", account.UserId)
	err = database.SelectRow(userSql, &usr.Id, &usr.Salt, &usr.Password, &usr.Name, &usr.CellPhone)
	if err != nil {
		log.Error().Msgf(err.Error())
		return nil, ""
	}

	// 验证密码是否正确
	if util.ValidatePassword(signin.Password, usr.Salt, usr.Password) {
		//生成token
		token, err := middlewares.ReleaseToken(usr)
		if err != nil {
			log.Error().Msgf(err.Error())
			return nil, ""
		}

		profile := map[string]string{
			"userId":    strconv.Itoa(usr.Id),
			"name":      usr.Name,
			"cellphone": usr.CellPhone,
		}
		return profile, token
	}
	log.Error().Msgf("login fail")
	return nil, ""
}
