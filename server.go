package foog

import (
	"net"
)

/**
运行服务：
	net.Listener,处理连接的接口
**/
type IServer interface {
	//连接，连接处理函数
	Run(net.Listener, func(IConn))
}
