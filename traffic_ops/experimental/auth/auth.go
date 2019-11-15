/*

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package main

import (
	"crypto/sha1"
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
	"strconv"
	"time"
)

// Config holds the configuration of the server.
type Config struct {
	DbName      string `json:"db-name"`
	DbUser      string `json:"db-user"`
	DbPassword  string `json:"db-password"`
	DbServer    string `json:"db-server"`
	DbPort      uint   `json:"db-port"`
	ListenPort  uint   `json:"listen-port"`
}

type Login struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type TmUser struct {
	Role        uint   `db:"role"`
	Password    string `db:"local_passwd"`
}

type Claims struct {
    Capabilities []string `json:"cap"`
    jwt.StandardClaims
}

type TokenResponse struct {
	Token string
}

var db *sqlx.DB // global and simple

var Logger *log.Logger

func printUsage() {
	exampleConfig := `{
	"db_name":     "to_development",
	"db_user":     "username",
	"db_password": "password",
	"db_server":   "localhost",
	"db_port":     5432,
	"listen_port": 9004
}`
	fmt.Println("Usage: " + path.Base(os.Args[0]) + " config-file secret")
	fmt.Println("")
	fmt.Println("Example config file:")
	fmt.Println(exampleConfig)
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

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

	http.HandleFunc("/login", handler)
	
	if _, err := os.Stat("server.pem"); os.IsNotExist(err) {
		Logger.Fatal("server.pem file not found")
	}

	if _, err := os.Stat("server.key"); os.IsNotExist(err) {
		Logger.Fatal("server.key file not found")
	}

	Logger.Printf("Starting server on port %d...", config.ListenPort)
	Logger.Fatal(http.ListenAndServeTLS(":" + strconv.Itoa(int(config.ListenPort)), "server.pem", "server.key", nil))
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
		var login Login
		tmUserlist := []TmUser{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Logger.Println("Error reading body: ", err.Error())
			http.Error(w, "Error reading body: "+err.Error(), http.StatusBadRequest)
			return
		}
		
		err = json.Unmarshal(body, &login)
		if err != nil {
			Logger.Println("Invalid JSON: ", err.Error())
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		
		stmt, err := db.PrepareNamed("SELECT role,local_passwd FROM tm_user WHERE username=:username")
		if err != nil {
			Logger.Println("Database error: ", err.Error())
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = stmt.Select(&tmUserlist, login)
		if err != nil {
			Logger.Println("Database error: ", err.Error())
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

    	hasher := sha1.New()
    	hasher.Write([]byte(login.Password))
    	hashedPassword := fmt.Sprintf("%x", hasher.Sum(nil))

		if len(tmUserlist) == 0 || tmUserlist[0].Password != string(hashedPassword) {
			Logger.Printf("Invalid username/password, username %s", login.Username)
			http.Error(w, "Invalid username/password", http.StatusUnauthorized)
			return
		}

		Logger.Printf("User %s authenticated", login.Username)

		claims := Claims {
	        []string{"read-ds", "write-ds", "read-cg"},	// TODO(amiry) - Adding hardcoded capabilities as a POC. 
	        											// Need to read from TO role tables when tables are ready
	        jwt.StandardClaims {
	        	Subject: login.Username,
	            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),	// TODO(amiry) - We will need to use shorter expiration, 
	            													// and use refresh tokens to extend access
	        },
	    }

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(os.Args[2]))
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

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	http.Error(w, r.Method+" "+r.URL.Path+" not valid for this microservice", http.StatusNotFound)
}