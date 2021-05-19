package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Handler
}

func NewServer() Server {
	return Server{}
}

func (s Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/report", s.handle(s.report)).Methods("POST")
	r.HandleFunc("/test", s.handle(s.test)).Methods("GET")
	log.Println("Start server on: :8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}

func (s Server) handle(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
