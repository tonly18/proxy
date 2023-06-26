package config

//config.toml 配置文件
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

//config data
var HttpConfig *httpConfigStruck = &httpConfigStruck{}

//init
func init() {
	if err := getConfig(); err != nil {
		panic(err)
	}
}

//获取config配置文件
func getConfig() error {
	if HttpConfig.HttpPort == 0 {
		if err := loadConfigFile("config", HttpConfig); err != nil {
			return err
		}
		HttpConfig.GameServerDoLoginAPI = HttpConfig.GameServerHost + HttpConfig.GameServerDoLoginAPI
		HttpConfig.GameServerCommandAPI = HttpConfig.GameServerHost + HttpConfig.GameServerCommandAPI
		HttpConfig.GameServerOnOffLineAPI = HttpConfig.GameServerHost + HttpConfig.GameServerOnOffLineAPI
	}
	return nil
}
