package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"wizard/pkg/db"
	"wizard/pkg/redis"
)

type App struct {
	DB    *db.DB
	Redis *redis.Client
	Request
}

func NewApp(r *http.Request) (app App, err error) {
	json.NewDecoder(r.Body).Decode(&app.Request)
	if len(app.Key) == 0 {
		err = errors.New("Empty key")
		return
	}
	db, err := db.InitDB()
	if err != nil {
		return
	}
	redis, err := redis.InitRedis()
	if err != nil {
		return
	}
	app.DB = db
	app.Redis = redis
	return
}

type Request struct {
	Key string `json:"dispatch_key"`
}
