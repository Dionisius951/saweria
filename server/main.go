package main

import (
	"server/tcp-udp"
	"server/websocket"
)

func main() {
	go websocket.StartWebSocket()
	go tcp_udp.HandleUDPConnection()
	go tcp_udp.StartTcp()
	select {}
}
