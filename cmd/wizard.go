package main

import (
	"log"
	"wizard/internal/server"
	"wizard/pkg/db"
	"wizard/pkg/redis"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	err = db.InitDB()
	if err != nil {
		return
	}
	defer db.Conn.Close()

	err = redis.InitRedis()
	if err != nil {
		return
	}
	defer redis.Cli.Close()

	server := server.NewServer()
	server.Run()
}
