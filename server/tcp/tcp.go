package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/scgywx/foog"
	"log"
	"net"
)

//头部的长度
type TcpServer struct {
	headSize int
}

//连接，实现了IConn接口：ReadMessage，WriteMessage，Close，GetRemoteAddr
type TcpConn struct {
	server     *TcpServer    //头部长度
	conn       net.Conn      //连接
	br         *bufio.Reader //读取数据
	bw         *bytes.Buffer //写的buffer
	remoteAddr string        //远程地址
}

//新建一个tcpserver headsize
func NewServer() *TcpServer {
	return &TcpServer{
		headSize: 4,
	}
}

//设置长度
func (this *TcpServer) SetHeadSize(n int) {
	this.headSize = n
}

//实现了Server方法
func (this *TcpServer) Run(ls net.Listener, fn func(foog.IConn)) {
	if this.headSize != 2 && this.headSize != 4 {
		this.headSize = 4
	}

	for {
		//接收连接
		c, err := ls.Accept()
		if err != nil {
			log.Println("Accept failed", err)
			break
		}
		//
		go fn(&TcpConn{
			server: this,                         //当前server
			conn:   c,                            //连接
			br:     bufio.NewReaderSize(c, 1024), //readersize
			//bufersize
			bw: bytes.NewBuffer(make([]byte, 0, 1024)),
			//地址
			remoteAddr: c.RemoteAddr().String(),
		})
	}
}

func (this *TcpConn) ReadMessage() ([]byte, error) {
	//获取长度数据
	headSize := this.server.headSize
	head, err := this.br.Peek(headSize)
	if err != nil {
		return nil, err
	}

	//丢弃长度数据
	this.br.Discard(headSize)
	bodySize := 0
	// 默认是大端，获取数据的长度
	if headSize == 2 {
		bodySize = int(binary.BigEndian.Uint16(head))
	} else {
		bodySize = int(binary.BigEndian.Uint32(head) & 0x7fffffff)
	}

	off := 0
	//创建数据的长度
	bytes := make([]byte, bodySize)
	//如果小于bodysize，持续去读取
	for off < bodySize {
		n, err := this.br.Read(bytes[off:])
		//获取错误
		if err != nil {
			return nil, err
		}
		//增加偏移量
		off += n
	}

	return bytes, nil
}

//写入数据
func (this *TcpConn) WriteMessage(msg []byte) error {
	headSize := this.server.headSize
	bodySize := len(msg)
	hdr := make([]byte, headSize)

	if headSize == 2 {
		binary.BigEndian.PutUint16(hdr, uint16(bodySize))
	} else {
		binary.BigEndian.PutUint32(hdr, uint32(bodySize))
	}

	//写入数据
	this.bw.Reset()
	this.bw.Write(hdr)
	this.bw.Write(msg)
	_, err := this.conn.Write(this.bw.Bytes())
	if err != nil {
		return err
	}

	return nil
}

//关闭连接
func (this *TcpConn) Close() {
	this.conn.Close()
}

//返回地址
func (this *TcpConn) GetRemoteAddr() string {
	return this.remoteAddr
}
