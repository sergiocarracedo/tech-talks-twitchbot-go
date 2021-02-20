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

	if err != nil {
		panic(err)
	}
}
