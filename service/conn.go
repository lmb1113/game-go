package main

import (
	"net"
	"sync"
)

var connMar sync.Map

func getConn(userId string) (net.Conn, bool) {
	conn, has := connMar.Load(userId)
	if has {
		return conn.(net.Conn), true
	}
	return nil, false
}

func setConn(userId string, room net.Conn) {
	connMar.Store(userId, room)
}
