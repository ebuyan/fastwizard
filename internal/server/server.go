package server

import (
	"log"
	"net/http"
	"wizard/internal/handler"

	"github.com/gorilla/mux"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/report", Handle(handler.GetReport)).Methods("POST")
	log.Println("Start server on: :8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
