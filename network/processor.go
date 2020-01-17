package network

// 消息处理器接口
// 1.消息注册
// 2.设置消息处理方法
// 3.设置消息处理服务端
// 4.消息序列化反序列化
type Processor interface {
	// must goroutine safe
	Route(msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
}
