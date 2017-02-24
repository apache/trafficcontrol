package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

// Config holds the configuration of the server.
type Config struct {
	DbName       string `json:"dbName"`
	DbUser       string `json:"dbUser"`
	DbPassword   string `json:"dbPassword"`
	DbServer     string `json:"dbServer,omitempty"`
	DbPort       uint   `json:"dbPort,omitempty"`
	ListenerPort string `json:"listenerPort"`
}

type User struct {
	Username  string `db:"username" json:"username"`
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Password  string `db:"password" json:"Password"`
}

type TokenResponse struct {
	Token string
}

var db *sqlx.DB // global and simple

func printUsage() {
	exampleConfig := `{
	"dbName":"my-db",
	"dbUser":"my-user",
	"dbPassword":"secret",
	"dbServer":"localhost",
	"dbPort":5432,
	"listenerPort":"8080"
}`
	Logger.Println("Usage: " + path.Base(os.Args[0]) + " configfile")
	Logger.Println("")
	Logger.Println("Example config file:")
	Logger.Println(exampleConfig)
}

var Logger *log.Logger

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// log.SetOutput(os.Stdout)
	Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	file, err := os.Open(os.Args[1])
	if err != nil {
		Logger.Println("Error opening config file:", err)
		return
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		Logger.Println("Error reading config file:", err)
		return
	}

	db, err = InitializeDatabase(config.DbUser, config.DbPassword, config.DbName, config.DbServer, config.DbPort)
	if err != nil {
		Logger.Println("Error initializing database:", err)
		return
	}

	http.HandleFunc("/", handler)
	if _, err := os.Stat("server.pem"); os.IsNotExist(err) {
		Logger.Fatal("server.pem file not found")
	}
	if _, err := os.Stat("server.key"); os.IsNotExist(err) {
		Logger.Fatal("server.key file not found")
	}
	Logger.Printf("Starting server on port " + config.ListenerPort + "...")
	Logger.Fatal(http.ListenAndServeTLS(":"+config.ListenerPort, "server.pem", "server.key", nil))
}

func InitializeDatabase(username, password, dbname, server string, port uint) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", server, dbname, username, password)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	Logger.Println(r.Method, r.URL.Scheme, r.Host, r.URL.RequestURI())
	if r.Method == "POST" {
		var u User
		userlist := []User{}
		body, err := ioutil.ReadAll(r.Body)
		log.Println(string(body))
		if err != nil {
			Logger.Println("Error reading body: ", err.Error())
			http.Error(w, "Error reading body: "+err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &u)
		if err != nil {
			Logger.Println("Error unmarshalling JSON: ", err.Error())
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		stmt, err := db.PrepareNamed("SELECT * FROM users WHERE username=:username")
		err = stmt.Select(&userlist, u)
		if err != nil {
			Logger.Println(err.Error())
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(userlist) == 0 || userlist[0].Password != u.Password {
			http.Error(w, "Invalid username/password ", http.StatusUnauthorized)
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["User"] = u.Username
		token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, err := token.SignedString([]byte("CAmeRAFiveSevenNineNine"))
		if err != nil {
			Logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(TokenResponse{Token: tokenString})
		if err != nil {
			Logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	http.Error(w, r.Method+" "+r.URL.Path+" not valid for this microservice", http.StatusNotFound)
}
