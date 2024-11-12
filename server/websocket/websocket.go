package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Struktur untuk Donasi
type Donation struct {
	Target  string `json:"target"`
	Amount  int    `json:"amount"`
	Message string `json:"message"`
	From    string `json:"from"`
}

var (
	Clients   = make(map[string]*websocket.Conn) // Memetakan username ke WebSocket client
	Broadcast = make(chan Donation)              // Channel untuk mengirimkan donasi
	mu        sync.Mutex
)

// Fungsi untuk menghubungkan WebSocket
func handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Ambil username dari query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	// Tambahkan klien dengan username ke daftar
	mu.Lock()
	Clients[username] = conn
	mu.Unlock()

	log.Printf("%s has joined the server\n", username)

	// Terima pesan dari klien
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			mu.Lock()
			delete(Clients, username)
			mu.Unlock()
			break
		}
	}
}

// Fungsi untuk mengirim donasi ke target tertentu
func sendDonationToTarget(donation Donation) {
	mu.Lock()
	defer mu.Unlock()

	// Cari klien yang sesuai dengan target donasi
	client, exists := Clients[donation.Target]
	if exists {
		err := client.WriteJSON(donation)
		if err != nil {
			log.Println("Error sending donation to target:", err)
			client.Close()
			delete(Clients, donation.Target)
		}
	} else {
		log.Println("Client not found for target:", donation.Target)
	}
}

// Fungsi untuk memulai WebSocket server
func StartWebSocket() {
	http.HandleFunc("/", handleConnection)

	// Jalankan server WebSocket
	log.Println("WebSocket server running at :8080")
	go http.ListenAndServe(":8080", nil)

	// Mendengarkan channel untuk donasi yang masuk
	for donation := range Broadcast {
		sendDonationToTarget(donation)
	}
}
