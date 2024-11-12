package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Struktur untuk Donasi
type Donation struct {
	Target  string `json:"target"`
	Amount  int    `json:"amount"`
	Message string `json:"message"`
	From    string `json:"from"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [username]")
		os.Exit(1)
	}

	username := os.Args[1]
	// Mempersiapkan alamat UDP tujuan (misalnya localhost pada port 8081)
	udpAddr := net.UDPAddr{
		Port: 8081,                     // Port UDP server
		IP:   net.ParseIP("127.0.0.1"), // IP server
	}

	// Membuka koneksi UDP ke alamat tujuan
	conn, err := net.DialUDP("udp", nil, &udpAddr)
	if err != nil {
		log.Fatal("Error opening UDP connection:", err)
	}
	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	for {
		// Membaca target penerima donasi
		fmt.Print("Masukkan username target penerima donasi: ")
		target, _ := inputReader.ReadString('\n')
		target = strings.TrimSpace(target)

		// Membaca jumlah donasi
		fmt.Print("Masukkan jumlah donasi: ")
		var amount int
		_, err := fmt.Scanln(&amount)
		if err != nil {
			fmt.Println("Error reading amount:", err)
			continue
		}

		// Membaca pesan donasi
		fmt.Print("Masukkan pesan untuk donasi: ")
		message, _ := inputReader.ReadString('\n')
		message = strings.TrimSpace(message)

		donation := Donation{
			Target:  target,
			Amount:  amount,
			Message: message,
			From:    username,
		}

		// Encode struktur donasi menjadi JSON
		data, err := json.Marshal(donation)
		if err != nil {
			log.Fatal("Error encoding JSON:", err)
		}

		// Kirim data donasi ke server UDP
		_, err = conn.Write(data)
		if err != nil {
			log.Fatal("Error sending UDP message:", err)
		}

		fmt.Println("Donasi terkirim ke target:", donation.Target)

		// clear screen
		fmt.Print("\033[H\033[2J")
		
		fmt.Print("Apakah Anda ingin mengirim donasi lain? (y/n): ")
		choice, _ := inputReader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if choice != "y" {
			break // Keluar dari loop jika bukan 'y'
		}
	}

}
