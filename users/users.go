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

	db, err = InitializeDatabase(config.DbUser, config.DbPassword, config.DbName, config.DbServer, config.DbPort)
	if err != nil {
		log.Println("Error initializing database:", err)
		return
	}

	var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Printf("Starting server on port " + config.ListenerPort + "...")

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

func retErr(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method, r.URL.Scheme, r.Host, r.URL.RequestURI())
	if r.Method == "GET" {
		if r.URL.Path == "/" {
			userlist := []User{}
			err := db.Select(&userlist, "SELECT * FROM users")
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.Encode(userlist)
		} else {
			username := strings.Replace(r.URL.Path, "/", "", 1)
			userlist := []User{}
			argument := User{}
			argument.Username = username
			stmt, err := db.PrepareNamed("SELECT * FROM users WHERE username=:username")
			err = stmt.Select(&userlist, argument)
			if err != nil {
				log.Println(err)
				retErr(w, http.StatusInternalServerError)
				return
			}
			if len(userlist) == 0 {
				retErr(w, http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.Encode(userlist[0])
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
			retErr(w, http.StatusInternalServerError)
			return
		}
		// TODO encrypt passwd before storing.
		sqlString := "INSERT INTO users (username, last_name, first_name, password) VALUES (:username, :last_name, :first_name, :password)"
		result, err := db.NamedExec(sqlString, u)
		if err != nil {
			log.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		rows, _ := result.RowsAffected()
		fmt.Fprintf(w, "Done! (%s Rows Affected)", rows)
	} else if r.Method == "PUT" {
		var u User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		u.Username = strings.Replace(r.URL.Path, "/", "", 1) // overwrite the username in the json, the path gets checked.
		// TODO encrypt passwd before storing.
		sqlString := "UPDATE users SET last_name=:last_name, first_name=:first_name, password=:password WHERE username=:username"
		result, err := db.NamedExec(sqlString, u)
		if err != nil {
			log.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		rows, _ := result.RowsAffected()
		fmt.Fprintf(w, "Done! (%s Rows Affected)", rows)
	} else if r.Method == "DELETE" {
		argument := User{}
		argument.Username = strings.Replace(r.URL.Path, "/", "", 1)
		result, err := db.NamedExec("DELETE FROM users WHERE username=:username", argument)
		if err != nil {
			log.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		rows, _ := result.RowsAffected()
		fmt.Fprintf(w, "Done! (%s Rows Affected)", rows)
	}
	retErr(w, http.StatusNotFound)
}
