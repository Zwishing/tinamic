package utils

import "github.com/spf13/viper"

func LoadConfig(){
	//通过配置文件初始化数据库连接
	viper.AddConfigPath("../config")
	viper.SetConfigName("tinamic")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return
	}
}

func UpdateConfig()  {

}

