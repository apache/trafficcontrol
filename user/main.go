package main

import (
	// "encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	// null "gopkg.in/guregu/null.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Password  string `db:"password" json:"Password"`
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

	// gob.Register(auth.SessionUser{}) // this is needed to pass the SessionUser struct around in the gorilla session.

	log.Println(config.DbUser, config.DbPassword, config.DbName, config.DbServer, config.DbPort)
	db, err = InitializeDatabase(config.DbUser, config.DbPassword, config.DbName, config.DbServer, config.DbPort)
	if err != nil {
		log.Println("Error initializing database:", err)
		return
	}

	var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Printf("Starting server on port " + config.ListenerPort + "...")

	// err = http.ListenAndServe(":"+config.ListenerPort, handlers.CombinedLoggingHandler(os.Stdout, routes.CreateRouter(dbb)))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println(err)
	}
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

	if r.Method == "GET" {
		if r.URL.Path == "/" {
			// TODO return list
		} else {
			username := strings.Replace(r.URL.Path, "/", "", 1)
			userlist := []User{}
			argument := User{}
			argument.Username = username
			stmt, err := db.PrepareNamed("SELECT * FROM users WHERE username=:username")
			err = stmt.Select(&userlist, argument)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.Encode(userlist)
		}
	} else if r.Method == "POST" {
		var u User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Println(err)
			// TODO return error.
		}
		sqlString := "INSERT INTO users (username, last_name, first_name, password) VALUES (:username, :last_name, :first_name, :password)"
		result, err := db.NamedExec(sqlString, u)
		if err != nil {
			log.Println(err)
			// TODO return error.
		}
		fmt.Fprintf(w, "Done! (%s)", result)
	} else if r.Method == "PUT" {

	} else if r.Method == "DELETE" {

	}
}
