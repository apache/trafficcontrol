package main

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

/*
This module it provides a server creating a RESTful API towards Traffic-Ops DB.
It was created in order to POC the concept of API Gateway (added on another PR) between the different modules of Traffic-Ops, specifically old perl Traffic-Ops, and new (go?) one.

At first step, only "tenant" table is covered (CRUD).
Other tables can be added, by adding modules (see tenant foe example), and seeding their enpoints below.
To Run this module, call: go run api.go --server --db-config-file 
e.g. go run api.go --server :8888 --db-config-file ../../app/conf/test/database.conf
Note that the below go modules are required, so you'll might need to "go get" them:)
"github.com/gorilla/mux"
"github.com/lib/pq"
*/

 
import (
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
  _ "github.com/lib/pq"

//modules to seed endpoints  
    "./tenant"
)

type DbConfig struct {
    Hostname  string    `json:"hostname"`
    Port      string    `json:"port"` //String for now, as the conf files are like this
    User      string    `json:"user"`
    Password  string    `json:"password"`
    DbName    string    `json:"dbname"`
}

func main() {

    server := flag.String("server", "", "IP:port to listen on")
    dbConfigFileName := flag.String("db-config-file", "", "DB to connect to config file")

    flag.Parse()

    if *server == ""{
        fmt.Println("Missing server address")
        os.Exit(1)
    }
    
    if *dbConfigFileName == ""{
        fmt.Println("Missing DB config file")
        os.Exit(1)
    }

    dbConfigFD, err := os.Open(*dbConfigFileName)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    
    var dbConfig DbConfig
    jsonParser := json.NewDecoder(dbConfigFD)
    if err := jsonParser.Decode(&dbConfig); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",dbConfig.Hostname, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err= db.Ping()
    if err != nil {
      panic(err)
    }

    fmt.Println("Successfully connected!")


    apiPrefix := "/api/2.0"
    router := mux.NewRouter()
    

    //list of elements to populat the api
    tenantEndpoints := tenant.NewEndpointSeeder(db)
    tenantEndpoints.Seed(router, apiPrefix)
    
    log.Fatal(http.ListenAndServe(*server, router))
}



