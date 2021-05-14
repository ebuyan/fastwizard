package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Handler
}

func NewServer() Server {
	return Server{NewHandler()}
}

func (s Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/report", s.report).Methods("POST")
	log.Println("Start server on: :8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
