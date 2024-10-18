package main

import (
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить все источники (не безопасно, но подойдет для MVP)
	},
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Message defines the structure of the messages
type Message struct {
	Type       string              `json:"type"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	StatusCode int                 `json:"status_code,omitempty"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем GET запрос до WebSocket протокола
	log.Println("Upgrading HTTP connection to WebSocket")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	log.Println("WebSocket connection established")
	clients[ws] = true

	for {
		var msg Message
		// Читаем сообщение из соединения
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		log.Printf("Received via WebSocket: %s %s", msg.Method, msg.URL)

		// Обработка сообщения
		if msg.Type == "http_response" {
			// Отправляем ответ клиенту
			broadcast <- msg
		}
	}
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("New HTTP request received")
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Формируем сообщение для отправки клиенту
	message := Message{
		Type:    "http_request",
		Method:  r.Method,
		URL:     r.URL.String(),
		Headers: r.Header,
		Body:    string(body),
	}

	// Отправляем сообщение клиенту
	for client := range clients {
		log.Printf("Sending to client via WebSocket: %s %s", message.Method, message.URL)
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}

	// Ожидаем ответ от клиента
	response := <-broadcast

	// Добавляем сообщение к телу ответа
	//response.Body += " Response from server"

	log.Printf("Response body: %s", response.Body)
	// Записываем ответ в HTTP-ответ клиенту
	for key, values := range response.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
	log.Printf("Response sent to HTTP client: %d", response.StatusCode)
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/", handleRequests)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
