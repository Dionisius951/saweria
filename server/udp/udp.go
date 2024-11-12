// udp_server.go
package udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// Struktur donasi
type Donation struct {
	Message string `json:"message"`
	Amount  int    `json:"amount"`
	From    string `json:"from"`
}

// Channel untuk mengirim donasi ke WebSocket
var donationChannel = make(chan Donation)

func HandleUDPConnection() {
	// Membuka UDP server pada port 8081
	addr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal("Error starting UDP server:", err)
	}
	defer conn.Close()

	fmt.Println("Listening for UDP donations on", addr.String())

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error receiving UDP message:", err)
			continue
		}

		var donation Donation
		err = json.Unmarshal(buffer[:n], &donation)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		// Mengirim donasi ke channel yang terhubung dengan WebSocket
		donationChannel <- donation
	}
}
