package main

import (
	"github.com/gorilla/mux"
	"WebApi_Go/service"
	"log"
	"net/http"
	"os"
)

func main() {
	service.OpenTaskDb()
	defer service.CloseTaskDb()

	r := mux.NewRouter()

	r.HandleFunc("/task", service.GetAllMethod).Methods("GET")
	r.HandleFunc("/task/{id}", service.GetByIdMethod).Methods("GET")
	r.HandleFunc("/task", service.CreateMethod).Methods("POST")
	r.HandleFunc("/task/{id}", service.DeleteMethod).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), r))
}
