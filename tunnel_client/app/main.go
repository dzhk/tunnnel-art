package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

var addr = "tunnel_server:8080" // Адрес сервера
var wsEndpoint = "ws://" + addr + "/ws"

func main() {
	// Устанавливаем соединение с сервером
	log.Printf("Connecting to WebSocket server at %s", wsEndpoint)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	log.Println("WebSocket connection established")

	for {
		// Читаем запрос от сервера
		var request Message
		err := conn.ReadJSON(&request)
		if err != nil {
			log.Println("read:", err)
			return
		}

		// Отправляем запрос на локальный домен
		localURL := "http://tunnel-web.local" + request.URL
		req, err := http.NewRequest(request.Method, localURL, bytes.NewBuffer([]byte(request.Body)))
		if err != nil {
			log.Println("new request:", err)
			continue
		}
		for key, values := range request.Headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("do request:", err)
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("read body:", err)
			continue
		}

		// Формируем ответное сообщение
		response := Message{
			Type:       "http_response",
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       string(body),
		}

		// Отправляем ответ на сервер
		err = conn.WriteJSON(response)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

// Message defines the structure of the messages
type Message struct {
	Type       string              `json:"type"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	StatusCode int                 `json:"status_code,omitempty"`
}
