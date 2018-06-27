package foog

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

//基本的app类
type Application struct {
	id         int            //当前的id
	listenAddr string         //监听地址
	server     IServer        //server：连接的读，写，关闭
	router     IRouter        //router：数据处理连接，数据处理关闭，处理消息
	serializer ISerializer    //serialize：编码和解码
	logFile    string         //日志文件
	logLevel   int            //日志级别
	handler    handlerManager // 处理器
}

//注册c
func (this *Application) Register(c IObject) {
	this.handler.register(c)
}

//设置监听地址
func (this *Application) Listen(addr string) {
	this.listenAddr = addr
}

//设置server
func (this *Application) SetServer(s IServer) {
	this.server = s
}

//设置router
func (this *Application) SetRouter(r IRouter) {
	this.router = r
}

//设置日志级别
func (this *Application) SetLogLevel(level int) {
	this.logLevel = level
}

//设置日志文件
func (this *Application) SetLogFile(filename string) {
	this.logFile = filename
}

//设置序列胡函数
func (this *Application) SetSerializer(s ISerializer) {
	this.serializer = s
}

//设置id
func (this *Application) SetId(id int) {
	this.id = id
}

//开始服务
func (this *Application) Start() {
	//初始化tcp连接
	ls, err := net.Listen("tcp", this.listenAddr)
	if err != nil {
		fmt.Println("listen server failed", err)
		return
	}

	//初始化日志处理
	if len(this.logFile) > 0 {
		w, err := os.OpenFile(this.logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println("open log file error", err)
			return
		}

		log.SetOutput(w)
	}
	//运行服务
	log.Println("server started", this.listenAddr)
	//处理连接，其实是把ls交给this.handleConnection去处理
	this.server.Run(ls, this.handleConnection)
}

//实际的处理连接的代码
func (this *Application) handleConnection(conn IConn) {
	//新建一个session
	sess := NewSession(conn, this.id)
	sess.serializer = this.serializer

	//关闭连接，关闭router
	defer conn.Close()
	//router关闭
	defer this.router.HandleClose(sess)

	//路由去处理连接
	this.router.HandleConnection(sess)
	for {
		//读取消息
		msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message failed", err)
			break
		}
		//router 处理消息
		name, data, err := this.router.HandleMessage(sess, msg)
		if err != nil {
			log.Println("handle message failed", err)
			break
		}
		//设置最后时间
		sess.LastTime = time.Now().Unix()

		//分发请求,到合适的方法去处理
		log.Printf("name=%s,data=%s", name, data)
		this.handler.dispatch(name, sess, data)
	}
}
