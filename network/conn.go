package network

import (
	"net"
)

// 网络链接对象接口
type Conn interface {
	// 读取消息
	ReadMsg() ([]byte, error)
	// 发送消息
	WriteMsg(args ...[]byte) error
	// 本地地址
	LocalAddr() net.Addr
	// 客户端地址
	RemoteAddr() net.Addr
	// 关闭链接，等待发送队列为空后断开链接
	Close()
	// 关闭链接，立即断开
	Destroy()
}
