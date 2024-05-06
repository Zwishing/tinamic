package user

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"tinamic/pkg/database"
	"tinamic/util"
	//"tinamic/util"
)

func ValidateUser(name, password string) bool {
	// 按照用户名从数据库中获取salt
	sql := fmt.Sprintf("SELECT salt,password FROM user_info.user u "+
		"INNER JOIN user_info.account a on u.id =a.user_id WHERE a.login_account = '%s'", name)
	var salt, hashPassword string
	err := database.Db.QueryRow(context.Background(), sql).Scan(&salt, &hashPassword)
	if err != nil {
		log.Error().Msgf("%s", err)
		return false
	}
	//验证密码
	return util.ValidatePassword(password, salt, hashPassword)

}
