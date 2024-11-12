package udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"server/websocket" // Import package websocket
)

// Struktur Donasi
type Donation struct {
	Target  string `json:"target"`
	Amount  int    `json:"amount"`
	Message string `json:"message"`
	From    string `json:"from"`
}

// Fungsi untuk menangani koneksi UDP
func HandleUDPConnection() {
	// Membuka UDP server pada port 8081
	udpAddr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		log.Fatal("Error starting UDP server:", err)
	}
	defer conn.Close()

	log.Println("Listening for UDP donations on", udpAddr.String())

	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading from UDP:", err)
			continue
		}

		var donation websocket.Donation
		err = json.Unmarshal(buffer[:n], &donation)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		fmt.Println(string(n))
		// Kirim donasi ke WebSocket berdasarkan target
		websocket.Broadcast <- donation
	}
}
