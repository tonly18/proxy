package global

// TCP包上传命令
const (
	CMD_UP_PING  = 100   //ping
	CMD_UP_LOGIN = 10000 //登录
	CMD_UP_PROTO = 35
)

// TCP包下发命令
const (
	CMD_DOWN_PONG     = 100
	CMD_DOWN_LOGIN    = 10000
	CMD_DOWN_KICK_OUT = 25 //踢玩家掉线
)
