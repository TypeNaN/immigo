package main

import (
	"context"
	"encoding/json"
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
	Port        string
	Key         string
	Cert        string
	Secret      []byte
	TokenExpire uint64
	Db          ConfigDB
}

type Users struct {
	DisplayName string `json:"displayname" bson:"displayname"`
	UserName    string `json:"username" bson:"username"`
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
}

var config Configure

func main() {
	fmt.Println("Hello world!")

	e := godotenv.Load(".env")
	X(e)

	exprire, _ := strconv.ParseUint(os.Getenv("token_expire"), 10, 64)

	config.Port = os.Getenv("port")
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

	var dbClientOpts = options.Client().ApplyURI(config.Db.Uri).SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      config.Db.User,
		Password:      config.Db.Pass,
	})

	dbClient, e = mongo.Connect(context.TODO(), dbClientOpts)
	X(e)

	databases, _ := dbClient.ListDatabaseNames(context.TODO(), bson.M{})
	fmt.Println(databases)

	router := mux.NewRouter()
	router.HandleFunc("/", simpleRes).Methods("GET")
	router.HandleFunc("/api/user/signup", userSignUp).Methods("POST")
	log.Fatal(http.ListenAndServeTLS(":"+config.Port, config.Cert, config.Key, router))
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
		log.Printf("Critical")
		panic(fmt.Sprintf("Critical error: %v", e))
	}
}

func simpleRes(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"Hello world"}`))
}

func userSignUp(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user Users
	json.NewDecoder(request.Body).Decode(&user)

	// ครวจสอบข้อมูลอย่างง่าย
	// ยังต้องการ การตรวจสอบที่รัดกุมกว่านี้
	if user.DisplayName == "" || user.UserName == "" || user.Email == "" || user.Password == "" {
		log.Printf("Warnning userSignUp: Bad Request \n %v", user)
		http.Error(response, "Status Bad Request", http.StatusBadRequest)
		return
	}

	var doc Users
	collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollUsers)
	e := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&doc)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			log.Printf("userSignUp: User %s already register", user.Email)
		}
	}
	if doc.Email == user.Email {
		log.Printf("userSignUp: User %s already existst", user.Email)
		response.Write([]byte(`{"message":"User already existst"}`))
		return
	}
	result, _ := collection.InsertOne(context.TODO(), user)
	json.NewEncoder(response).Encode(map[string]interface{}{
		"insertedID": result.InsertedID,
		"message":    "successfully",
	})
}
