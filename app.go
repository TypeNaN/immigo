package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

type UserSign struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type CustomerInfo struct {
	Email string
	Kind  string
}

type CustomClaims struct {
	*jwt.StandardClaims
	TokenLevel string
	Info       CustomerInfo
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
	router.HandleFunc("/api/user/signin", userSignIn).Methods("POST")
	router.HandleFunc("/api/user/verify", userVerify).Methods("GET")
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

	user.Password = getHash([]byte(user.Password))

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

func userSignIn(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user UserSign
	json.NewDecoder(request.Body).Decode(&user)
	var doc Users
	collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollUsers)
	e := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&doc)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			log.Printf("userSignIp: User %s not found", user.Email)
			response.Write([]byte(`{"message":"User not found"}`))
			return
		}
	}

	Upass := []byte(user.Password)
	Dpass := []byte(doc.Password)
	e = bcrypt.CompareHashAndPassword(Dpass, Upass)
	if e != nil {
		log.Printf("userSignIp: Wrong password!, User %s", user.Email)
		response.Write([]byte(`{"message":"Wrong Password!"}`))
		return
	}

	exprire, jwtToken, e := generateJWT(user.Email)
	if e != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + e.Error() + `"}`))
		return
	}

	log.Printf("userSignIp: User %s Signed in", user.Email)
	json.NewEncoder(response).Encode(map[string]interface{}{
		"exprire": exprire,
		"token":   jwtToken,
	})
}

func getHash(password []byte) string {
	hash, e := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	E(e)
	return string(hash)
}

func generateJWT(email string) (int64, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	exprire := &jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.TokenExpire)).Unix()}
	token.Claims = &CustomClaims{exprire, "0", CustomerInfo{email, "human"}}
	tokenString, e := token.SignedString(config.Secret)
	if e != nil {
		log.Println("GenerateJWT: Error in JWT token generation")
		return 0, "", e
	}
	return exprire.ExpiresAt, tokenString, nil
}

func getVerifyToken(authorization string) (interface{}, error) {
	if authorization == "" {
		return &CustomerInfo{Email: "", Kind: ""}, errors.New("403 Forbidden")
	}
	tokenHead := strings.Split(authorization, "Bearer ")[1]
	token, e := jwt.ParseWithClaims(tokenHead, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Secret, nil
	})
	if e != nil {
		return &CustomerInfo{Email: "", Kind: ""}, e
	}
	return token.Claims, nil
}

func userVerify(response http.ResponseWriter, request *http.Request) {
	claims, err := getVerifyToken(request.Header.Get("Authorization"))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	customs := claims.(*CustomClaims)
	info := customs.Info
	response.Write([]byte(`{
		"verify":{
			"Level":"` + customs.TokenLevel + `",
			"Email":"` + info.Email + `",
			"Kind":"` + info.Kind + `"
		}
	}`))
}
