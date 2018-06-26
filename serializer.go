package foog

//序列化
type ISerializer interface {
	//编码
	Encode(interface{}) ([]byte, error)
	//解码
	Decode([]byte, interface{}) error
}
