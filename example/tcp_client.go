package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9005")

	str := `{"cmd":"hello.say","data":{"name":"beck"}}`
	msg := []byte(str)
	bodySize := len(msg)
	fmt.Printf("bodySize=%d", bodySize)

	headerData := make([]byte, 4)
	binary.BigEndian.PutUint32(headerData, uint32(bodySize))

	conn.Write(headerData)
	conn.Write(msg)

	bt := make([]byte, 1023)
	n, err := conn.Read(bt)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bt[:n]))
	}
	defer conn.Close()
}
