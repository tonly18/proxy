package ziface

type IHeartbeatChecker interface {
	Start()
	Stop()
	BindConn(IConnection)
	Clone() IHeartbeatChecker
}
