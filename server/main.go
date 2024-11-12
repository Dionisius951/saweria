package main
import (
	"server/websocket"
	"server/udp"
)

func main() {
	go websocket.StartWebSocket()
	go udp.HandleUDPConnection()
	select{}
}
