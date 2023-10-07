package config

import (
	"fmt"
	"github.com/spf13/viper"
	"proxy/server/global"
)

//获取配置文件并解析到指定的struck
func loadConfigFile(fname string, configStruck any) error {
	//viper
	viper.AddConfigPath(global.PROXY_SERVER_WORK_PATH_ENV + "/conf")
	viper.SetConfigName(getConfigFileName(fname))
	viper.SetConfigType("toml")

	//reade
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal Error Config File, File Name: %s, Error: %s", fname, err)
	}

	//parse
	viper.Unmarshal(configStruck)

	return nil
}

//获取配置文件
func getConfigFileName(fname string) string {
	if global.PROXY_SERVER_ENV == "" {
		return fname
	}

	return fname + "_" + global.PROXY_SERVER_ENV
}
