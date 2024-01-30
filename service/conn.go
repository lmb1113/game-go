package main

import (
	"fmt"
	"net"
	"sync"
)

var connMar sync.Map
var remoteAddrMar sync.Map

func getConn(userId uint64) (net.Conn, bool) {
	conn, has := connMar.Load(userId)
	if has {
		return conn.(net.Conn), true
	}
	fmt.Println("找不到连接")
	return nil, false
}

func setConn(userId uint64, conn net.Conn) {
	connMar.Store(userId, conn)
}

func deleteConn(userId uint64) {
	connMar.Delete(userId)
}

func getRemoteAddr(addr string) (uint64, bool) {
	value, has := remoteAddrMar.Load(addr)
	if has {
		return value.(uint64), true
	}
	return 0, false
}

func setRemoteAddr(addr string, userId uint64) {
	remoteAddrMar.Store(addr, userId)
}

func deleteRemoteAddr(addr string) {
	remoteAddrMar.Delete(addr)
}
