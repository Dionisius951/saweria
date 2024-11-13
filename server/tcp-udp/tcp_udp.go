package tcp_udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"server/websocket" // Import package websocket
	"sync"
)

var mu sync.Mutex

// Struktur Donasi
type Donation struct {
	Type    string `json:"type"`
	Target  string `json:"target"`
	Amount  int    `json:"amount"`
	Message string `json:"message"`
	From    string `json:"from"`
}

type BalanceReq struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

var UserBalance = make(map[string]int)

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

	log.Println("Server UDP Running on", udpAddr.String())

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading from UDP:", err)
			continue
		}

		var request Donation
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		response := ""
		switch request.Type {
		case "balance":
			// Cek saldo pengguna
			balance := checkBalance(request.From)
			response = fmt.Sprintf("Saldo %s saat ini adalah: %d", request.From, balance)
			sendResponse(response, conn, addr)
		case "donation":
			// Tangani data donasi dan kirim ke WebSocket
			handleDonationData(request.From, request.Amount, buffer[:n], conn, addr)
		default:
			// Jika tipe tidak dikenali, kirim pesan error
			response = "Tipe permintaan tidak dikenali."
		}

		// mengirim respon server ke client
		

	}
}

func sendResponse(message string, conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte(message), addr)
	if err != nil {
		log.Println("Error sending response:", err)
	}
}

func handleDonationData(username string, amount int, data []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	var donation websocket.Donation
	err := json.Unmarshal(data, &donation)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	balance := UserBalance[username]
	if balance < amount {
		sendResponse("Saldo tidak mencukupi", conn, addr)
		return
		
	}
	sendResponse("Donasi Berhasil terkirim, saldo telah berkurang", conn, addr)
	balance -= amount
	UserBalance[username] = balance

	// Kirim donasi ke WebSocket
	websocket.Broadcast <- donation
}

func checkBalance(username string) int {
	mu.Lock()
	defer mu.Unlock()
	// Pastikan username sudah ada di map
	balance, exists := UserBalance[username]
	if !exists {
		// Jika username tidak ada, inisialisasi dengan 0
		UserBalance[username] = 0
		balance = 0
	}
	return balance
}

// Fungsi untuk memulai server TCP menangani permintaan top-up
func StartTcp() {
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on 127.0.0.1:9090")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming request
	decoder := json.NewDecoder(conn)
	var topUpReq BalanceReq
	err := decoder.Decode(&topUpReq)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	// Convert balance to integer and update the user balance
	balance := topUpReq.Balance

	// Perbarui saldo pengguna
	totalBalance := updateBalance(topUpReq.Username, balance)

	// Kirim respons kembali
	response := fmt.Sprintf("Berhasil menambah saldo menjadi : %d", totalBalance)
	conn.Write([]byte(response))
}

func updateBalance(username string, amount int) int {
	mu.Lock()
	defer mu.Unlock()

	// Jika username belum ada, inisialisasi dengan saldo 0
	if _, exists := UserBalance[username]; !exists {
		UserBalance[username] = 0
	}

	// Tambahkan saldo
	UserBalance[username] += amount
	return UserBalance[username]
}
