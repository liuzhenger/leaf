package gate

import (
	"net"
)

/*
 Agent 客户端代理
 主要是给业务层使用
 获取客户端信息和发送消息
*/

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}
