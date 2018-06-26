package foog

import (
	"time"
)

//定义基本的session
type Session struct {
	serializer ISerializer //序列化对象
	Id         int64       //定义id
	Conn       IConn       //连接管理
	LastTime   int64       //最后时间
	Data       interface{} //数据
}

//定义全局数，计算数量
var counter int64 = 0

/**
 * 1位符号
 * 31位时间戳(最大可表示到2038年)
 * 10位毫秒
 * 10位服务器ID(最大可表示1024)
 * 12位自增id(最大值是4096)
 * 共64位，每秒可生成400w条不同ID
 */
func NewSession(conn IConn, appId int) *Session {
	counter++
	//定义一个session
	sess := &Session{
		//????
		//计算一个id
		Id:   ((time.Now().UnixNano() / 1000000) << 22) | int64((appId&0x3ff)<<12) | (counter & 0xfff),
		Conn: conn,
	}
	return sess
}

func (this *Session) WriteMessage(data interface{}) error {
	if msg, ok := data.([]byte); ok || this.serializer == nil {
		//如果是byte或者serializer为nil直接写
		return this.Conn.WriteMessage(msg)
	} else {
		//否则编码后写
		bytes, err := this.serializer.Encode(data)
		if err != nil {
			return err
		}

		return this.Conn.WriteMessage(bytes)
	}
}
