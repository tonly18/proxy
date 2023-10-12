package config

import "proxy/server/library"

// config.toml 配置文件
type httpConfigStruck struct {
	HttpHost               string
	HttpPort               int
	GameServerHost         string
	GameServerDoLoginAPI   string
	GameServerCommandAPI   string
	GameServerOnOffLineAPI string

	//PHP api
	GameServerPHPCommandAPI string
}

// config data
var HttpConfig *httpConfigStruck = &httpConfigStruck{}

// init
func init() {
	if err := LoadProxyConfig(); err != nil {
		panic(err)
	}
}

// 获取config配置文件
func LoadProxyConfig() error {
	if err := loadConfigFile("config", HttpConfig); err != nil {
		return err
	}

	//设置url轮询初始数据
	library.NewPoll().Set(HttpConfig.GameServerHost)

	//return
	return nil
}
