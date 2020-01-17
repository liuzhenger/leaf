package network

/*
 Agent 链接代理
*/

type Agent interface {
	// 消息解析和分发
	Run()
	// 断开连接，对应Conn接口中的Close方法
	OnClose()
}
