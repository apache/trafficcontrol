package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	_ "github.com/lib/pq"
)

const TrafficOpsDomain = "localhost"
const CookieName = "access_token"
const DBTable = "traffic_ops_auth_users"

// AllowedCreateUserRoles returns the database roles which are allowed to create users.
// This should be treated as a constant, and only exists because Go doesn't allow constant map literals.
func AllowedCreateUserRoles() map[string]struct{} {
	return map[string]struct{}{
		"admin": struct{}{},
	}
}

// Config holds the configuration of the server.
type Config struct {
	DBServer  string `json:"database_server"`
	DBPort    uint   `json:"database_port"`
	DBName    string `json:"database_name"`
	DBUser    string `json:"database_user"`
	DBPass    string `json:"database_password"`
	AdminUser string `json:"admin_user"`
	AdminPass string `json:"admin_pass"`
	Port      uint   `json:"port"`
	TokenKey  string `json:"token_key"`
}

func createConnectionStringPostgres(server, database, user, pass string, port uint) (string, error) {
	connString := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", database, user, pass)
	if server != "" {
		connString += fmt.Sprintf(" host=%s", server)
	}
	if port > 0 && port < 65536 {
		connString += fmt.Sprintf(" port=%d", port)
	}
	return connString, nil
}

// authRolesStr returns a string of roles allowed to create users, as expected by SQL `IN` clauses.
func authRolesStr() string {
	roles := AllowedCreateUserRoles()
	s := "("
	for role, _ := range roles {
		s += "'" + role + "'" + ","
	}
	s = s[:len(s)-1] // strip trailing ,
	s += ")"
	return s
}

const CreateUserPath = "/create_user/"

// handleCreateUser is the HTTP handler for the create_user endpoint.
// It checks that the current logged in user is authorized to create users, and if so, creates the user.
// hasAdmin is a mutable pointer, so we don't have to check every time. Once an admin exists, the pointer is set to true
// TODO change to POST
func HandleCreateUser(db *sql.DB, jwtSigningKey string, w http.ResponseWriter, r *http.Request) {
	token, err := getTokenData(jwtSigningKey, r)
	if err != nil {
		log.Printf("%s ERROR unauthorized: %s\n", r.RemoteAddr, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if _, ok := AllowedCreateUserRoles()[token.Role]; !ok {
		log.Printf("%s ERROR unauthorized: %s role %s is not allowed to create users\n", r.RemoteAddr, token.User, token.Role)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	path := r.URL.String()[len(CreateUserPath):]
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		log.Printf("%s ERROR invalid request: not enough parts: '%s'. Syntax is '/create_user/name/pass/role\n", r.RemoteAddr, path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := parts[0]
	pass := parts[1]
	role := parts[2]

	if err := createUser(db, user, pass, role); err != nil {
		// TODO return 400 if the error is that the user exists
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s ERROR failed to insert into database: '%v'\n", r.RemoteAddr, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type TokenData struct {
	User string
	Role string
}

func getTokenData(jwtSigningKey string, r *http.Request) (*TokenData, error) {
	encToken, err := r.Cookie(CookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(
		[]byte(encToken.Value),
		jwt.WithVerify(jwa.HS256, []byte(jwtSigningKey)),
	)
	if err != nil {
		return nil, err
	}

	userInterface, hasUser := token.Get("user")
	roleInterface, hasRole := token.Get("role")
	user, userIsStr := userInterface.(string)
	role, roleIsStr := roleInterface.(string)
	if !hasUser || !hasRole || !userIsStr || !roleIsStr {
		// we signed it: this should never happen
		return nil, fmt.Errorf("token missing claims")
	}

	return &TokenData{User: user, Role: role}, nil
}

const LoginPath = "/login/"

func HandleLogin(db *sql.DB, jwtSigningKey string, w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "text/plain")
	path := r.URL.String()[len(LoginPath):]
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s ERROR invalid request: not enough parts: '%s'. Syntax is '/get_user/name/pass'\n", r.RemoteAddr, path)
		return
	}
	user := parts[0]
	pass := parts[1]

	role, err := getUser(db, user, pass)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("%s ERROR failed to get user '%s' with %v\n", r.RemoteAddr, path, err)
		return
	}

	token, err := jwt.NewBuilder().
		Claim(`user`, user).
		Claim(`role`, role).
		Build()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR creating token: %s", err)
		return
	}
	signed, err := jwt.Sign(token, jwa.HS256, []byte(jwtSigningKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s ERROR creating token for '%s': %v\n", r.RemoteAddr, user, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  string(signed),
		Domain: TrafficOpsDomain,
		Path:   "/",
		//Secure: true, // TODO uncomment when https is implemented
		HttpOnly: true, // prevents the cookie being accessed by Javascript. DO NOT remove, security vulnerability
	})

	w.WriteHeader(http.StatusOK)
	log.Printf("%s got user %s role %s\n", r.RemoteAddr, user, role)
}

type UserInfo struct {
	User string `json:"user"`
	Role string `json:"role"`
}

const UserInfoPath = "/user_info"

func HandleUserInfo(db *sql.DB, jwtSigningKey string, w http.ResponseWriter, r *http.Request) {
	token, err := getTokenData(jwtSigningKey, r)
	if err != nil {
		log.Printf("%s ERROR unauthorized: %s\n", r.RemoteAddr, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userInfo := UserInfo{User: token.User, Role: token.Role}
	jsonUserInfo, err := json.Marshal(userInfo)
	if err != nil {
		log.Printf("%s ERROR marshalling userinfo %v: %s\n", r.RemoteAddr, userInfo, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", jsonUserInfo)
}

func getConfig(file string) (Config, error) {
	configJson, err := ioutil.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = json.Unmarshal(configJson, &config)
	return config, err
}

func getDB(server string, port uint, name string, user string, pass string) (*sql.DB, error) {
	connStr, err := createConnectionStringPostgres(server, name, user, pass, port)
	if err != nil {
		return nil, err
	}
	return sql.Open("postgres", connStr)
}

// Generates a base64-encoded 128-bit random number.
func randGUIDStr() string {
	guid := make([]byte, 16, 16)
	rand.Read(guid)
	return base64.RawURLEncoding.EncodeToString(guid)
}

func adminExists(db *sql.DB) (bool, error) {
	authRolesCount := 0
	err := db.QueryRow(`SELECT COUNT(1) FROM "` + DBTable + `" WHERE role IN ` + authRolesStr() + `;`).Scan(&authRolesCount)
	if err != nil {
		return false, err
	}
	return authRolesCount > 0, nil
}

// TODO prepare query
func createUser(db *sql.DB, user, pass, role string) error {
	salt := randGUIDStr()
	hashedPassBytes := sha512.New().Sum([]byte(pass + salt))
	hashedPass := base64.RawURLEncoding.EncodeToString(hashedPassBytes)
	_, err := db.Exec(`insert into "`+DBTable+`" (username, hash, salt, role) VALUES ($1, $2, $3, $4);`, user, hashedPass, salt, role)
	return err
}

// getUser returns the user's role, if authentication was successful. If authentication fails, err != nil
func getUser(db *sql.DB, user, pass string) (string, error) {
	salt := ""
	hashedPass := ""
	role := ""
	err := db.QueryRow(`select salt, hash, role from "`+DBTable+`" where username = $1;`, user).Scan(&salt, &hashedPass, &role)
	if err != nil {
		return "", err
	}

	hashedPassBytes := sha512.New().Sum([]byte(pass + salt))
	sentHashedPass := base64.RawURLEncoding.EncodeToString(hashedPassBytes)

	if hashedPass != sentHashedPass {
		return "", fmt.Errorf("authentication failed")
	}

	return role, nil
}

func GetRoutes() map[string]func(db *sql.DB, jwtSigningKey string, w http.ResponseWriter, r *http.Request) {
	return map[string]func(db *sql.DB, jwtSigningKey string, w http.ResponseWriter, r *http.Request){
		LoginPath:      HandleLogin,
		CreateUserPath: HandleCreateUser,
		UserInfoPath:   HandleUserInfo,
	}
}

func main() {
	config, err := getConfig("config.json")
	if err != nil {
		log.Println(err)
		return
	}

	if len(config.TokenKey) < 22 {
		log.Printf("Token key is too short. Must be at least 128 bits. Suggested alternative: '%s'\n", randGUIDStr())
		return
	}

	db, err := getDB(config.DBServer, config.DBPort, config.DBName, config.DBUser, config.DBPass)
	if err != nil {
		log.Println(err)
		return
	}

	dbHasAdmin, err := adminExists(db)
	if err != nil {
		log.Printf("Error checking for admin: %v\n", err)
		return
	}
	if !dbHasAdmin && (config.AdminUser == "" || config.AdminPass == "") {
		log.Println("Error: no admin in config or database")
		return
	}

	if !dbHasAdmin {
		firstAllowedCreateRole := ""
		for role, _ := range AllowedCreateUserRoles() {
			firstAllowedCreateRole = role
			break
		}
		if firstAllowedCreateRole == "" {
			log.Printf("ERROR creating initial admin user: no roles are allowed to create users\n")
			return
		}
		err := createUser(db, config.AdminUser, config.AdminPass, firstAllowedCreateRole)
		if err != nil {
			log.Printf("ERROR creating initial admin user: %v\n", err)
			return
		}
		log.Printf("INFORMATION created initial admin user: %s\n", config.AdminUser)
	}

	if dbHasAdmin && (config.AdminUser != "" || config.AdminPass != "") {
		log.Printf("WARNING admin exists in both database and config: ignoring config admin\n")
	}

	wrapHandleFunc := func(f func(db *sql.DB, jwtSigningKey string, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			f(db, config.TokenKey, w, r)
		}
	}

	routes := GetRoutes()
	for path, handleFunc := range routes {
		http.HandleFunc(path, wrapHandleFunc(handleFunc))
	}

	if err := http.ListenAndServe(":"+strconv.Itoa(int(config.Port)), nil); err != nil {
		log.Println(err)
		return
	}
}
