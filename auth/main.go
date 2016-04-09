package main

import (
	// "encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	// null "gopkg.in/guregu/null.v3"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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
	log.Println("Usage: " + path.Base(os.Args[0]) + " configfile")
	log.Println("")
	log.Println("Example config file:")
	log.Println(exampleConfig)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	log.SetOutput(os.Stdout)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println("Error opening config file:", err)
		return
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("Error reading config file:", err)
		return
	}

	db, err = InitializeDatabase(config.DbUser, config.DbPassword, config.DbName, config.DbServer, config.DbPort)
	if err != nil {
		log.Println("Error initializing database:", err)
		return
	}

	var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Printf("Starting server on port " + config.ListenerPort + "...")

	http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)
	http.ListenAndServeTLS(":8080", "server.pem", "server.key", nil)

	if err != nil {
		log.Println(err)
	}
}
func validateToken(tokenString string) (*jwt.Token, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("CAmeRAFiveSevenNineNine"), nil
	})

	if err == nil && token.Valid {
		log.Println("TOKEN IS GOOD -- user:", token.Claims["userid"], " role:", token.Claims["role"])
	} else {
		log.Println("TOKEN IS BAD", err)
	}
	return token, err
}

func InitializeDatabase(username, password, dbname, server string, port uint) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", server, dbname, username, password)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func retErr(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method, r.URL.Scheme, r.Host, r.URL.RequestURI())
	if r.Method == "POST" {
		var u User
		userlist := []User{}
		username := ""
		password := ""
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body: ", err.Error())
			http.Error(w, "Error reading body: "+err.Error(), http.StatusBadRequest)
			return
		}
		// var lj loginJson
		// log.Println(body)
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Println("Error unmarshalling JSON: ", err.Error())
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		// username = lj.User
		// password = lj.Password

		stmt, err := db.PrepareNamed("SELECT * FROM users WHERE username=:username")
		err = stmt.Select(&userlist, u)
		if err != nil {
			log.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		if len(userlist) == 0 || userlist[0].Password != u.Password {
			retErr(w, http.StatusUnauthorized)
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["User"] = username
		token.Claims["Password"] = password
		token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, err := token.SignedString([]byte("CAmeRAFiveSevenNineNine"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(TokenResponse{Token: tokenString})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		// var u User
		// userlist := []User{}
		// body, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	log.Println(err)
		// }
		// err = json.Unmarshal(body, &u)
		// if err != nil {
		// 	log.Println(err)
		// 	retErr(w, http.StatusInternalServerError)
		// 	return
		// }

	}
	retErr(w, http.StatusNotFound)
}
