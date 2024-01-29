package main

import (
	"fmt"
	"net"
	"sync"
)

var connMar sync.Map

func getConn(userId uint64) (net.Conn, bool) {
	conn, has := connMar.Load(userId)
	if has {
		return conn.(net.Conn), true
	}
	fmt.Println("找不到连接")
	return nil, false
}

func setConn(userId uint64, room net.Conn) {
	connMar.Store(userId, room)
}
