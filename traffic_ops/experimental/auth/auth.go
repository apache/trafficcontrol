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
	"bytes"
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
	DbName         string `json:"db-name"`
	DbUser         string `json:"db-user"`
	DbPassword     string `json:"db-password"`
	DbServer       string `json:"db-server"`
	DbPort         uint   `json:"db-port"`
	ListenPort     uint   `json:"listen-port"`
	LegacyLoginURL string `json:"legacy-login-url"`
}

type Login struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type LegacyLogin struct {
	Username    string `json:"u"`
	Password    string `json:"p"`
}

type TmUser struct {
	UserId      int    `db:"id"`
	Password    string `db:"local_passwd"`
}

type UserRole struct {
	RoleId      int    `db:"role_id"`
}

type Capability struct {
	CapName     string `db:"cap_name"`
}

type Claims struct {
    Capabilities []string     `json:"cap"`
    LegacyCookie string       `json:"legacy-cookie"`	// LEGACY: The legacy cookie to be passed to API GW
    jwt.StandardClaims
}

type TokenResponse struct {
	Token string
}

var db *sqlx.DB // global and simple

var Logger *log.Logger

func printUsage() {
	exampleConfig := `{
	"db_name":         "to_development",
	"db_user":         "username",
	"db_password":     "password",
	"db_server":       "localhost",
	"db_port":         5432,
	"listen_port":     9004,
	"legacy_to_login": "http://localhost:3000/api/1.2/user/login"
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

	handler, _ := makeHandler(&config)
	http.HandleFunc("/login", handler)
	
	if _, err := os.Stat("server.crt"); os.IsNotExist(err) {
		Logger.Fatal("server.crt file not found")
	}

	if _, err := os.Stat("server.key"); os.IsNotExist(err) {
		Logger.Fatal("server.key file not found")
	}

	Logger.Printf("Starting server on port %d...", config.ListenPort)
	Logger.Fatal(http.ListenAndServeTLS(":" + strconv.Itoa(int(config.ListenPort)), "server.crt", "server.key", nil))
}

func InitializeDatabase(username, password, dbname, server string, port uint) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", server, dbname, username, password)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func LegacyTOLogin(login Login, legacyLoginURL string, w http.ResponseWriter) (*http.Response, error) {

	// TODO(amiry) - Legacy token expiration should be longer than JWT expiration

	legacyLogin := LegacyLogin{ login.Username, login.Password }

 	body, err := json.Marshal(legacyLogin)
    if err != nil {
		Logger.Println("JSON marshal error: ", err.Error())
        return nil, err
    }

	req, err := http.NewRequest("POST", legacyLoginURL,  bytes.NewBuffer(body))
	client := &http.Client{}
    resp, err := client.Do(req)
	if err != nil {
		Logger.Println("Legacy Login error: ", err.Error(), " Legacy URL: ", legacyLoginURL)
		return nil, err;
	}

	return resp, err
}

func makeHandler(config *Config) (func(http.ResponseWriter, *http.Request), error) {

	return func (w http.ResponseWriter, r *http.Request) {

		Logger.Println(r.Method, r.URL.Scheme, r.Host, r.URL.RequestURI())

		if r.Method == "POST" {

			var login Login
			tmUserList := []TmUser{}
			userRoleList := []UserRole{}
			capList := []Capability{}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				Logger.Printf("Error reading request body: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			
			err = json.Unmarshal(body, &login)
			if err != nil {
				Logger.Printf("JSON error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Get the user id and the password from tm_user, in order to validate the user's password
			stmt, err := db.PrepareNamed("SELECT id,local_passwd FROM tm_user WHERE username=:username")
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = stmt.Select(&tmUserList, login)
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

	    	hasher := sha1.New()
	    	hasher.Write([]byte(login.Password))
	    	hashedPassword := fmt.Sprintf("%x", hasher.Sum(nil))

			if len(tmUserList) == 0 || tmUserList[0].Password != string(hashedPassword) {
				Logger.Printf("Invalid username/password. Username=%s]", login.Username)
				http.Error(w, "Invalid username/password", http.StatusUnauthorized)
				return
			}

			// We have validated the user's password, now lets get the user's roles
			stmt, err = db.PrepareNamed("SELECT role_id FROM user_role WHERE user_id=:id")
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = stmt.Select(&userRoleList, tmUserList[0])
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			rolesIds := []int{}
			for _, elem := range userRoleList {
				rolesIds = append(rolesIds, elem.RoleId)
			}

			// Get user's capabilities according to the user's roles
			sql, args, err := sqlx.In("SELECT cap_name FROM role_capability WHERE role_id IN (?)", rolesIds)
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Replace the "?" bindvar syntax with DB specific syntax ($1, $2, ... for PostgreSQL)
			sql = db.Rebind(sql)

			stmt1, err := db.Preparex(sql)
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = stmt1.Select(&capList, args...)
			if err != nil {
				Logger.Printf("DB error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			capabilities := []string{}
			for _, elem := range capList {
				capabilities = append(capabilities, elem.CapName)
			}

			/////////////////////////////////////////////////////////////////////////////////////////////////
			// LEGACY: Perform login against legacy TO. This is required until AAA is disabled in TO
			legacyResp, err := LegacyTOLogin(login, config.LegacyLoginURL, w);
			if err != nil {
				Logger.Printf("Legacy TO login error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return;
			}

			if legacyResp.StatusCode != http.StatusOK {
				Logger.Printf("Legacy TO login returned bad status: %s", legacyResp.Status)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			legacyCookies := legacyResp.Cookies()

			if (legacyCookies == nil) || (len(legacyCookies) != 1) || (legacyCookies[0].Name != "mojolicious") {
				Logger.Printf("Error parsing legacy response cookie. Legacy cookies %v ", legacyCookies)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			Logger.Printf("LEGACY LOGIN TOKEN: %s %s %s", legacyCookies[0].Name, legacyCookies[0].Value, legacyCookies[0].Expires)
			
			// LEGACY: End
			/////////////////////////////////////////////////////////////////////////////////////////////////

			Logger.Printf("User %s authenticated. Role Ids %v. Capabilities %v", login.Username, rolesIds, capabilities)

			claims := Claims {
				
				capabilities,											// Set capabilities
				legacyCookies[0].String(),								// LEGACY: Set legacy cookie

		        jwt.StandardClaims {
		        	Subject: login.Username,
		            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),	// TODO(amiry) - We will need to use shorter expiration, 
		            													// and use refresh tokens to extend access.
		            													// Expiration time should be configurable.
		        },
		    }

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenSignedString, err := token.SignedString([]byte(os.Args[2]))
			if err != nil {
				Logger.Printf("JWT error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			js, err := json.Marshal(TokenResponse{Token: tokenSignedString})
			if err != nil {
				Logger.Printf("JWT error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}, nil
}
