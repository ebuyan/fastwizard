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

	app.DB = db.Conn
	app.Redis = redis.Cli
	return
}

type Request struct {
	Key string `json:"dispatch_key"`
}
