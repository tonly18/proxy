package zconf

import (
	"encoding/json"
	"fmt"
	"os"
	"proxy/core/zinx/zlog"
	"time"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type Config struct {
	/*
		Server
	*/
	Host     string //当前服务器主机IP
	TCPPort  int    //当前服务器主机监听端口号
	Name     string //当前服务器名称
	ServerID uint32 //服务器ID

	/*
		Zinx
	*/
	Version          string //当前Zinx版本号
	MaxConn          int    //当前服务器主机允许的最大链接个数
	MaxPacketSize    uint32 //都需数据包的最大值
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	/*
		config file path
	*/
	ConfFilePath string

	/*
		logger
	*/
	LogDir  string //日志所在文件夹 默认"./log"
	LogFile string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr

	/*
		conn 读写时间
	*/
	MaxConnReadTime  int //conn 读时间：单位秒
	MaxConnWriteTime int //conn 写时间：单位毫秒

	//conn 最长心跳检测间隔时间(单位:秒),超过改时间间隔,则认为超时,从配置文件读取
	HeartbeatMax int
}

/*
定义一个全局的对象
*/
var GlobalObject *Config

// PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (g *Config) Reload() {
	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists == false {
		panic(fmt.Sprintf(`Config File:%v is not exist!`, g.ConfFilePath))
	}

	data, err := os.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	if err := json.Unmarshal(data, g); err != nil {
		panic(err)
	}

	//Logger 设置
	zlog.SetLogFile(g.LogDir, g.LogFile)
}

func (g *Config) HeartbeatMaxDuration() time.Duration {
	return time.Duration(g.HeartbeatMax) * time.Second
}

/*
提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &Config{
		Host:    "0.0.0.0",
		TCPPort: 7000,
		Name:    "Zinx TCP Server",

		Version:          "V1.0",
		MaxConn:          1000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,

		ConfFilePath: fmt.Sprintf(`%v/zinx_%v.json`, ZINX_CONFIG_PATH, ZINX_ENV),

		LogDir:  ZINX_LOG_PATH,
		LogFile: "zinx.log",

		MaxConnReadTime:  0,
		MaxConnWriteTime: 0,

		HeartbeatMax: 0,
	}

	//NOTE: 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
