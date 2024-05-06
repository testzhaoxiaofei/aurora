package setting

import (
	"fmt"
	"github.com/spf13/viper"
)

var Type string

// Setup 主进程初始化配置
func Setup() {
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err)
	}

}
