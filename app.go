package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type SPAhandler struct {
	staticPath string
	indexPath  string
}

type FileImage struct {
	Path        string `json:"path" bson:"path"`
	File        string `json:"file" bson:"file"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Owner       string `json:"owner" bson:"owner"`
	User        string `json:"user" bson:"user"`
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
	router.HandleFunc("/api/user/signup", userSignUp).Methods("POST")
	router.HandleFunc("/api/user/signin", userSignIn).Methods("POST")
	router.HandleFunc("/api/user/verify", userVerify).Methods("GET")
	router.HandleFunc("/api/user/refresh", userRefresh).Methods("GET")
	router.HandleFunc("/api/image/all", imageAll).Methods("GET")
	router.HandleFunc("/api/image/upload", imageUpload).Methods("POST")
	router.HandleFunc("/api/image/edit", imageEdit).Methods("POST")
	router.HandleFunc("/api/image/remove", imageRemove).Methods("POST")
	spa := SPAhandler{staticPath: "public", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)
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

func (h SPAhandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	path, err := filepath.Abs(request.URL.Path)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(response, request, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(response, request)
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
		return &CustomerInfo{}, errors.New("403 Forbidden")
	}
	tokenHead := strings.Split(authorization, "Bearer ")[1]
	token, e := jwt.ParseWithClaims(tokenHead, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Secret, nil
	})
	if e != nil {
		return &CustomerInfo{}, e
	}
	return token.Claims, nil
}

func userVerify(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	claims, err := getVerifyToken(request.Header.Get("Authorization"))
	if err != nil {
		response.WriteHeader(http.StatusForbidden)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(&claims)
}

func getRefreshToken(response http.ResponseWriter, request *http.Request) (int64, string, error) {
	claims, err := getVerifyToken(request.Header.Get("Authorization"))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return 0, "", err
	}
	customs := claims.(*CustomClaims)
	info := customs.Info
	exprire, jwtToken, err := generateJWT(info.Email)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return exprire, jwtToken, err
	}
	return exprire, jwtToken, nil
}

func userRefresh(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	ex, token, _ := getRefreshToken(response, request)
	exprire := fmt.Sprintf("%v", ex)
	response.Write([]byte(`{"exprire":` + exprire + `, "token":"` + token + `"}`))
}

func getUserInfo(response http.ResponseWriter, request *http.Request) (Users, error) {
	var user Users
	claims, e := getVerifyToken(request.Header.Get("Authorization"))
	if e != nil {
		response.WriteHeader(http.StatusForbidden)
		response.Write([]byte(`{"message":"` + e.Error() + `"}`))
		return user, e
	}
	customs := claims.(*CustomClaims)
	info := customs.Info
	collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollUsers)
	e = collection.FindOne(context.TODO(), bson.M{"email": info.Email}).Decode(&user)
	if e != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + e.Error() + `"}`))
		return user, e
	}
	return user, e
}

func imageUpload(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user, _ := getUserInfo(response, request)

	request.ParseMultipartForm(10 << 20)

	prefix := request.FormValue("name")
	if utf8.RuneCountInString(prefix) < 1 {
		fmt.Println("Error Retrieving the Name")
		return
	}

	description := request.FormValue("description")
	if utf8.RuneCountInString(prefix) < 1 {
		fmt.Println("Error Retrieving the Description")
		return
	}

	m := request.MultipartForm
	files := m.File["myFile[]"]

	var complete = make(map[string]interface{})
	for i, file := range files {
		reader, e := file.Open()
		if e != nil {
			http.Error(response, e.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()
		path := "public/upload-images/" + user.UserName
		_, e = os.Stat(path)
		if os.IsNotExist(e) {
			e = os.MkdirAll(path, 0660|os.ModePerm)
			E(e)
		}
		tempFile, e := ioutil.TempFile(path, prefix+"-*.png")
		E(e)
		defer tempFile.Close()
		fileInfo, e := tempFile.Stat()
		E(e)
		fileBytes, e := ioutil.ReadAll(reader)
		E(e)
		tempFile.Write(fileBytes)

		var doc FileImage
		var image FileImage
		image.Path = strings.Replace(tempFile.Name(), "public/", "", -1)
		image.File = fileInfo.Name()
		image.Name = prefix
		image.Description = description
		image.Owner = user.Email
		image.User = user.UserName

		collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollImgs)
		e = collection.FindOne(context.TODO(), bson.M{"owner": user.Email, "path": image.Path}).Decode(&doc)
		if e != nil {
			if e != mongo.ErrNoDocuments {
				log.Printf("Error: imageUpload: %v", e)
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(`{"message":"` + e.Error() + `"}`))
				return
			}
		}
		if doc.Owner != "" {
			response.WriteHeader(http.StatusNotModified)
			response.Write([]byte(`{"message":"Duplicated"}`))
			return
		}
		result, _ := collection.InsertOne(context.TODO(), image)
		complete[fmt.Sprint(i)] = result
		log.Printf("LOG: user %s upload file %s", user.Email, image.Path)
	}
	exprire, token, _ := getRefreshToken(response, request)
	refrash := map[string]interface{}{
		"exprire": exprire,
		"token":   token,
	}
	json.NewEncoder(response).Encode(map[string]interface{}{
		"insertedID": complete,
		"refresh":    refrash,
		"message":    "successfully",
	})
}

func imageAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user, _ := getUserInfo(response, request)

	var results []bson.M
	collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollImgs)
	cursor, e := collection.Find(context.TODO(), bson.M{"owner": user.Email})
	E(e)
	e = cursor.All(context.TODO(), &results)
	E(e)

	var image bson.M
	var images = make(map[string]interface{})

	for _, result := range results {
		image = bson.M{
			"id":          result["_id"],
			"name":        result["name"],
			"description": result["description"],
			"path":        result["path"],
		}
		id := fmt.Sprintf("%v", result["_id"])
		id = strings.Replace(id, `ObjectID("`, "", -1)
		id = strings.Replace(id, `")`, "", -1)

		images[id] = image
	}
	json.NewEncoder(response).Encode(images)
}

func imageEdit(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user, _ := getUserInfo(response, request)

	var images = make(map[string]interface{})
	json.NewDecoder(request.Body).Decode(&images)

	id := fmt.Sprintf("%v", images["id"])
	docID, e := primitive.ObjectIDFromHex(id)
	E(e)
	var doc FileImage
	filter := bson.M{"_id": docID}
	update := bson.M{"$set": bson.M{"name": images["name"], "description": images["description"]}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollImgs)
	e = collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&doc)
	if e != nil {
		log.Printf("Error: imageEdit: %v", e)
		if e == mongo.ErrNoDocuments {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{"message":"` + e.Error() + `"}`))
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + e.Error() + `"}`))
	}
	log.Printf("LOG: user %s edit file %s", user.Email, doc.Path)
	exprire, token, _ := getRefreshToken(response, request)
	refrash := map[string]interface{}{
		"exprire": exprire,
		"token":   token,
	}
	json.NewEncoder(response).Encode(map[string]interface{}{
		"updateID": docID,
		"refresh":  refrash,
		"message":  "successfully",
	})
}

func imageRemove(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user, _ := getUserInfo(response, request)

	var images = make(map[string]interface{})
	json.NewDecoder(request.Body).Decode(&images)

	id := fmt.Sprintf("%v", images["id"])
	docID, e := primitive.ObjectIDFromHex(id)
	E(e)
	var doc FileImage
	filter := bson.M{"_id": docID}
	opts := options.FindOne()
	collection := dbClient.Database(config.Db.Db).Collection(config.Db.CollImgs)
	e = collection.FindOne(context.TODO(), filter, opts).Decode(&doc)
	if e != nil {
		log.Printf("Error: imageEdit: %v", e)
		if e == mongo.ErrNoDocuments {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{"message":"` + e.Error() + `"}`))
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + e.Error() + `"}`))
	}
	log.Printf("LOG: user %s remove file %s", user.Email, doc.Path)

	path := fmt.Sprintf("public/%v", doc.Path)
	e = os.Remove(path)
	E(e)

	res, e := collection.DeleteOne(context.TODO(), filter)
	if e != nil {
		log.Printf("Error: imageEdit: %v", e)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + e.Error() + `"}`))
	}

	log.Printf("LOG: user %s deleted %v documents and remove file %s ", user.Email, res.DeletedCount, doc.Path)

	exprire, token, _ := getRefreshToken(response, request)
	refrash := map[string]interface{}{
		"exprire": exprire,
		"token":   token,
	}
	json.NewEncoder(response).Encode(map[string]interface{}{
		"removeID": docID,
		"refresh":  refrash,
		"message":  "successfully",
	})
}
