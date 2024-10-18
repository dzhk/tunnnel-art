package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

//	func startTelegramBot() {
//		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_TOKEN"))
//		if err != nil {
//			log.Fatal("Error creating bot: " + err.Error())
//		}
//
//		bot.Debug = true
//
//		log.Printf("Authorized on account %s", bot.Self.UserName)
//
//		u := tgbotapi.NewUpdate(0)
//		u.Timeout = 60
//
//		updates := bot.GetUpdatesChan(u)
//
//		for update := range updates {
//			if update.Message != nil { // If we got a message
//				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
//				msg.ReplyToMessageID = update.Message.MessageID
//
//				bot.Send(msg)
//			}
//		}
//	}
func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет web!")
}

func startHTTPServer() {
	http.HandleFunc("/", sayhello)         // Устанавливаем роутер
	err := http.ListenAndServe(":80", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	startHTTPServer()
}
