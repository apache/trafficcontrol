package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
)

// Config holds the configuration of the server.
type Config struct {
	DbName       string `json:"dbName"`
	DbCamera     string `json:"dbUser"`
	DbPassword   string `json:"dbPassword"`
	DbServer     string `json:"dbServer,omitempty"`
	DbPort       uint   `json:"dbPort,omitempty"`
	ListenerPort string `json:"listenerPort"`
}

type Camera struct {
	Name     string `db:"name" json:"name"`
	Owner    string `db:"owner" json:"owner"`
	Type     string `db:"type" json:"type"`
	URL      string `db:"url" json:"url"`
	Location string `db:"location" json:"location"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type Response struct {
	Status     string   `json:"Status"`
	Message    string   `json:"Message,omitempty"`
	CameraData []Camera `json:"CameraData,omitempty"`
}

var db *sqlx.DB // global and simple
var Logger *log.Logger

func printUsage() {
	exampleConfig := `{
	"dbName":"my-db",
	"dbCamera":"my-user",
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

func main() {
	if len(os.Args) < 2 {
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

	db, err = InitializeDatabase(config.DbCamera, config.DbPassword, config.DbName, config.DbServer, config.DbPort)
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

func retErr(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func handler(w http.ResponseWriter, r *http.Request) {

	Logger.Println(r.Method, r.URL.Scheme, r.Host, r.URL.RequestURI())
	msg := "None"
	cameralist := []Camera{}
	if r.Method == "GET" {
		if r.URL.Path == "/cameras" || r.URL.Path == "/cameras/" {
			err := db.Select(&cameralist, "SELECT * FROM cameras")
			if err != nil {
				Logger.Println(err)
			}
		} else {
			re := regexp.MustCompile("[0-9]+")
			params := re.FindAllString(r.URL.Path, -1)
			argument := Camera{Owner: params[0], Name: params[1]}
			// argument.name = params[1]
			// argument.Owner = params[0]
			stmt, err := db.PrepareNamed("SELECT * FROM cameras WHERE owner=:Owner and name=:Name")
			err = stmt.Select(&cameralist, argument)
			if err != nil {
				Logger.Println(err)
				retErr(w, http.StatusInternalServerError)
				return
			}
			if len(cameralist) == 0 {
				retErr(w, http.StatusNotFound)
				return
			}
		}
	} else if r.Method == "POST" {
		var c Camera
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Logger.Println(err)
		}
		err = json.Unmarshal(body, &c)
		if err != nil {
			Logger.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		// TODO encrypt passwd before storing.
		sqlString := "INSERT INTO cameras (name, owner, type, url, location, username, password) " +
			" VALUES (:name, :owner, :type, :url, :location, :username, :password)"
		_, err = db.NamedExec(sqlString, c)
		if err != nil {
			Logger.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		msg = "Camera successully created"
	} else if r.Method == "PUT" {
		var c Camera
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Logger.Println(err)
		}
		err = json.Unmarshal(body, &c)
		if err != nil {
			Logger.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		// overwrite the fields in the payload - the path gets checked.
		re := regexp.MustCompile("[0-9]+")
		params := re.FindAllString(r.URL.Path, -1)
		c.Owner = params[0]
		c.Name = params[1]
		// TODO encrypt passwd before storing.
		sqlString := "UPDATE cameras SET type=:type, url=:url, location=:location, username=:username, password=:password " +
			"WHERE owner=:owner AND name=:name"
		Logger.Println(sqlString)
		_, err = db.NamedExec(sqlString, c)
		if err != nil {
			Logger.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		msg = "Camera successfully updated"
	} else if r.Method == "DELETE" {
		argument := Camera{}
		re := regexp.MustCompile("[0-9]+")
		params := re.FindAllString(r.URL.Path, -1)
		argument.Owner = params[0]
		argument.Name = params[1]
		_, err := db.NamedExec("DELETE FROM cameras WHERE name=:name and owner=:owner", argument)
		if err != nil {
			Logger.Println(err)
			retErr(w, http.StatusInternalServerError)
			return
		}
		msg = "Camera successfully deleted"
	} else {
		http.Error(w, r.Method+" "+r.URL.Path+" not valid for this microservice", http.StatusNotFound)
	}
	resp := &Response{}
	if len(cameralist) == 0 {
		resp = &Response{Status: "Success", Message: msg}
	} else {
		resp = &Response{Status: "Success", CameraData: cameralist}
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
