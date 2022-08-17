package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

type ConfigDB struct {
	Uri       string
	User      string
	Pass      string
	Db        string
	CollUsers string
	CollImgs  string
}

type Configure struct {
	Port        uint64
	Key         string
	Cert        string
	Secret      []byte
	TokenExpire uint64
	Db          ConfigDB
}

var config Configure

func simpleRes(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"Hello world"}`))
}

func main() {
	fmt.Println("Hello world!")

	e := godotenv.Load(".env")
	X(e)

	port, _ := strconv.ParseUint(os.Getenv("port"), 10, 64)
	exprire, _ := strconv.ParseUint(os.Getenv("token_expire"), 10, 64)

	config.Port = port
	config.Key = os.Getenv("server_key")
	config.Cert = os.Getenv("server_cert")
	config.Secret = []byte(os.Getenv("server_secret"))
	config.TokenExpire = exprire
	config.Db.Uri = os.Getenv("db_uri")
	config.Db.User = os.Getenv("db_user")
	config.Db.Pass = os.Getenv("db_pass")
	config.Db.Db = os.Getenv("db_name")
	config.Db.CollUsers = os.Getenv("db_coll_users")
	config.Db.CollImgs = os.Getenv("db_coll_images")

	fmt.Println(config)

	var dbClientOpts = options.Client().ApplyURI(config.Db.Uri).SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      config.Db.User,
		Password:      config.Db.Pass,
	})

	dbClient, _ = mongo.Connect(context.TODO(), dbClientOpts)
	databases, _ := dbClient.ListDatabaseNames(context.TODO(), bson.M{})
	fmt.Println(databases)

	router := mux.NewRouter()
	router.HandleFunc("/", simpleRes).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func W(w error) {
	if w != nil {
		log.Printf("Warnning: %v", w)
	}
}

func E(e error) {
	if e != nil {
		log.Printf("Error: %v", e)
	}
}

func X(e error) {
	if e != nil {
		panic(fmt.Sprintf("Critical error: %v", e))
	}
}
