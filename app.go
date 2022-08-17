package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func simpleRes(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"Hello world"}`))
}

func main() {
	fmt.Println("Hello world!")

	var dbClientOpts = options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      "username",
		Password:      "password",
	})

	dbClient, _ = mongo.Connect(context.TODO(), dbClientOpts)
	databases, _ := dbClient.ListDatabaseNames(context.TODO(), bson.M{})
	fmt.Println(databases)

	router := mux.NewRouter()
	router.HandleFunc("/", simpleRes).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
