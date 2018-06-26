package foog

//定义连接的接口
type IConn interface {
	//读消息
	ReadMessage() ([]byte, error)
	//写消息
	WriteMessage([]byte) error
	//关闭连接
	Close()
	//获取远程地址
	GetRemoteAddr() string
}
