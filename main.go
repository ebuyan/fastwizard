package main

import (
	"log"
	"wizard/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	server := server.NewServer()
	server.Run()
}
