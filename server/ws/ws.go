package ws

import (
	"github.com/gorilla/websocket"
	"github.com/scgywx/foog"
	"log"
	"net"
	"net/http"
)

//websocker处理
type WebSocketServer struct {
	handle   func(foog.IConn)   //处理函数
	msgType  int                //消息类型
	upgrader websocket.Upgrader //更新
}

//websocket连接
type WebSocketConn struct {
	conn       *websocket.Conn
	msgType    int
	remoteAddr string
}

//新建一个server
func NewServer() *WebSocketServer {
	s := &WebSocketServer{}
	s.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return s
}

//设置变量
func (this *WebSocketServer) SetCheckOriginFunc(fn func(r *http.Request) bool) {
	this.upgrader.CheckOrigin = fn
}

//设置变量
func (this *WebSocketServer) SetMessageType(msgType int) {
	this.msgType = msgType
}

//运行websocket服务
func (this *WebSocketServer) Run(ls net.Listener, fn func(foog.IConn)) {
	this.handle = fn
	http.HandleFunc("/", this.handleConnection)
	http.Serve(ls, nil)
}

func (this *WebSocketServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	//处理websocket的更新
	c, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("websocket upgrade:", err)
		return
	}
	//把更新交给WebSocketConn去处理
	this.handle(&WebSocketConn{
		conn:       c,
		msgType:    this.msgType,
		remoteAddr: r.RemoteAddr,
	})
}

//读取消息：类型+消息+err
func (this *WebSocketConn) ReadMessage() ([]byte, error) {
	mt, msg, err := this.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if this.msgType == 0 {
		this.msgType = mt
	}

	return msg, nil
}

//写入消息类型+消息
func (this *WebSocketConn) WriteMessage(msg []byte) error {
	return this.conn.WriteMessage(this.msgType, msg)
}

//关闭
func (this *WebSocketConn) Close() {
	this.conn.Close()
}

//返回remoteAddr
func (this *WebSocketConn) GetRemoteAddr() string {
	return this.remoteAddr
}
