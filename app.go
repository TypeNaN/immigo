package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func simpleRes(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"Hello world"}`))
}

func main() {
	fmt.Println("Hello world!")
	router := mux.NewRouter()
	router.HandleFunc("/", simpleRes).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
