package main

import (
	"github.com/joho/godotenv"
	"log"
	"sergiocarracedo.es/streambot-go/cmd/bootstrap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = bootstrap.Run()

	//client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
	//	chat-logger.Println("User Notice message")
	//	chat-logger.Printf("%#v", message)
	//})
	//
	//client.Join(os.Getenv("CHANNEL"))

	if err != nil {
		panic(err)
	}
}
