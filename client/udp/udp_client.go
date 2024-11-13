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
	Type    string `json:"type"`
	Target  string `json:"target"`
	Amount  int    `json:"amount"`
	Message string `json:"message"`
	From    string `json:"from"`
}

type BalanceRequest struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
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

	for {
		// Menampilkan menu utama
		fmt.Println("|-----------------|")
		fmt.Println("|      Menu       |")
		fmt.Println("|-----------------|")
		fmt.Println("|1. Kirim Donasi  |")
		fmt.Println("|2. Cek Saldo     |")
		fmt.Println("|3. Exit          |")
		fmt.Println("|-----------------|")

		// Menampilkan pilihan menu
		fmt.Print("\n Pilih Menu : ")
		var option int
		_, err := fmt.Scanln(&option)
		if err != nil {
			fmt.Println("Error Input....")
		}
		switch option {
		case 1:
			SendDonation(conn, username)
		case 2:
			CheckBalance(conn, username)
		case 3:
			// clear screen
			fmt.Print("\033[H\033[2J")
			// Keluar dari program
			fmt.Println("Terima kasih!")
			os.Exit(0)
		}

		// // Menanyakan apakah ingin kembali ke menu utama
		// fmt.Print("Apakah Anda ingin kembali ke menu utama? (y/n): ")
		// var choice string
		// _, err = fmt.Scanln(&choice)
		// if err != nil || strings.ToLower(choice) != "y" {
		// 	fmt.Println("Terima kasih telah menggunakan program.")
		// 	break
		// }
	}
}

// Fungsi untuk mengirim donasi
func SendDonation(conn *net.UDPConn, username string) {
	// clear screen
	fmt.Print("\033[H\033[2J")

	inputReader := bufio.NewReader(os.Stdin)

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
		return
	}

	// Membaca pesan donasi
	fmt.Print("Masukkan pesan untuk donasi: ")
	message, _ := inputReader.ReadString('\n')
	message = strings.TrimSpace(message)

	donation := Donation{
		Type:    "donation",
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

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error membaca respons dari server:", err)
		return
	}

	// Menampilkan pesan dari server
	response := string(buf[:n])
	fmt.Println(response)

}

// Fungsi untuk mengecek saldo
func CheckBalance(conn *net.UDPConn, username string) {
	// clear screen
	fmt.Print("\033[H\033[2J")
	// Membuat request untuk memeriksa saldo
	request := Donation{Type: "balance", From: username}

	// Encode data request menjadi JSON
	data, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	// Kirim request ke server UDP
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal("Error sending UDP message:", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error membaca respons dari server:", err)
		return
	}

	// Menampilkan pesan dari server
	response := string(buf[:n])
	fmt.Println(response)

}
