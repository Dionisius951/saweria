package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Donation struct {
	Message string `json:"message"`
	Amount  int    `json:"amount"`
	From    string `json:"from"`
}

func main() {
	// Membuat koneksi ke server UDP di alamat 127.0.0.1:8081
	serverAddr := "127.0.0.1:8081"
	remoteAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatal("Error connecting to UDP server:", err)
	}
	defer conn.Close()

	// Memasukkan data donasi dari pengguna
	var donation Donation
	fmt.Print("Enter your name: ")
	fmt.Scanln(&donation.From)

	fmt.Print("Enter your message: ")
	fmt.Scanln(&donation.Message)

	fmt.Print("Enter donation amount: ")
	fmt.Scanln(&donation.Amount)

	// Mengonversi data donasi ke format JSON
	data, err := json.Marshal(donation)
	if err != nil {
		log.Fatal("Error encoding donation:", err)
	}

	// Mengirim pesan donasi ke server UDP
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal("Error sending donation:", err)
	}

	fmt.Println("Donation sent successfully!")
}
