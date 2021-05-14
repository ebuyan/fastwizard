package main

import (
	"log"
	"wizard/psql"
	"wizard/redis"
	"wizard/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	err = psql.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	err = redis.InitRedis()
	if err != nil {
		log.Fatalln(err)
	}
	server := server.NewServer()
	server.Run()
}
