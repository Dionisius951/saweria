package websocket

import (
	"encoding/json"
	"flag"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn)
var mu sync.Mutex

// Channel untuk komunikasi antara UDP dan WebSocket
var donationChannel = make(chan Donation)

type Donation struct {
	Message string `json:"message"`
	Amount  int    `json:"amount"`
	From    string `json:"from"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection ke WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Mendapatkan username dari query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	// Menyimpan koneksi WebSocket untuk username
	mu.Lock()
	clients[username] = conn
	mu.Unlock()

	fmt.Printf("User %s connected\n", username)

	// Mendengarkan pesan WebSocket dari client (ini opsional, tergantung kebutuhan)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("User %s disconnected\n", username)
			mu.Lock()
			delete(clients, username)
			mu.Unlock()
			break
		}
	}
}

// Fungsi untuk mengirim pesan donasi ke WebSocket client yang sesuai
func handleDonations() {
	for donation := range donationChannel {
		mu.Lock()
		client, exists := clients[donation.From]
		mu.Unlock()
		if exists {
			data, _ := json.Marshal(donation)
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("Error sending message to WebSocket:", err)
				client.Close()
				mu.Lock()
				delete(clients, donation.From)
				mu.Unlock()
			}
		} else {
			log.Printf("No WebSocket client found for username: %s\n", donation.From)
		}
	}
}

func StartWebSocket() {
	flag.Parse()
	log.SetFlags(0)
	// http.HandleFunc("/echo", echo)
	http.HandleFunc("/donate", wsHandler)
	go handleDonations()

	fmt.Println("server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
	
