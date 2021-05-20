package server

import (
	"log"
	"net/http"
	"wizard/internal/app"
)

func Handle(f func(app app.App) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app, err := app.NewApp(r)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer app.DB.Close()
		defer app.Redis.Close()

		result, err := f(app)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}
