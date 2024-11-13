package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Struktur untuk permintaan top-up
type TopUpRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Penggunaan: go run tcp_client.go [username]")
		return
	}

	username := os.Args[1]

	for {
		// Menampilkan menu
		fmt.Println("\n|-----------------|")
		fmt.Println("|   Menu Top-Up   |")
		fmt.Println("|-----------------|")
		fmt.Println("|1. Top-Up Saldo  |")
		fmt.Println("|2. Keluar        |")
		fmt.Println("|-----------------|")

		fmt.Print("Pilih opsi: ")
		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			sendTopUpRequest(username)
		case 2:
			fmt.Println("Terima Kasih!")
			return
		default:
			fmt.Println("Opsi tidak valid. Coba lagi.")
		}
	}
}

// Fungsi untuk mengirim permintaan top-up ke server
func sendTopUpRequest(username string) {
	// clear screen
	fmt.Print("\033[H\033[2J")

	fmt.Print("Masukkan jumlah top-up: ")
	var amount int
	_, err := fmt.Scanln(&amount)
	if err != nil {
		fmt.Println("Error reading amount:", err)
		return
	}

	// Menghubungkan ke server TCP
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("Gagal menghubungkan ke server TCP:", err)
		return
	}
	defer conn.Close()

	// Menyiapkan permintaan top-up dalam format JSON
	topUpReq := TopUpRequest{
		Type:     "topup",
		Username: username,
		Balance:  amount,
	}
	data, err := json.Marshal(topUpReq)
	if err != nil {
		fmt.Println("Gagal mengonversi ke JSON:", err)
		return
	}

	// Mengirim permintaan top-up ke server
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Gagal mengirim permintaan top-up:", err)
		return
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error membaca respons dari server:", err)
		return
	}

	// Menampilkan saldo
	fmt.Println(string(buf[:n]))
}
