package main

import (
	"encoding/json"
	"fmt"
	"foog"
	"foog/server/ws"

	sjson "foog/serializer/json"
)

//定义消息的格式
type MyRequest struct {
	//命令
	Cmd string `json:"cmd"`
	//数据
	Data map[string]interface{} `json:"data"`
}

//定义router
type MyRouter struct {
}

func (this *MyRouter) HandleConnection(sess *foog.Session) {
	fmt.Printf("new client from %s, #%d\n", sess.Conn.GetRemoteAddr(), sess.Id)
}

func (this *MyRouter) HandleClose(sess *foog.Session) {
	fmt.Printf("client close #%d\n", sess.Id)
}

func (this *MyRouter) HandleMessage(sess *foog.Session, msg []byte) (string, interface{}, error) {
	req := &MyRequest{}
	//把消息解析到message里面
	json.Unmarshal(msg, req)
	return req.Cmd, req.Data, nil
}

//响应数据
type SayResponse struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

//注册服务
type Hello struct {
}

func (this *Hello) Say(
	sess *foog.Session,
	req map[string]interface{}) {
	//发送响应数据
	rsp := &SayResponse{
		Name: req["name"].(string),
		Text: fmt.Sprintf("hello %s", req["name"]),
	}
	fmt.Println(rsp.Text)
	sess.WriteMessage(rsp)
}

func main() {
	//初始化
	app := &foog.Application{}
	//设置router
	app.SetRouter(&MyRouter{})
	app.SetServer(ws.NewServer())
	app.SetSerializer(sjson.New())
	app.Register(&Hello{})
	//设置监听端口
	app.Listen("127.0.0.1:9005")

	//处理请求
	app.Start()
}
