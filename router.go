package foog

//路由处理接口
type IRouter interface {
	//处理连接
	HandleConnection(*Session)
	//处理关闭
	HandleClose(*Session)
	//处理消息
	HandleMessage(*Session, []byte) (string, interface{}, error)
}
