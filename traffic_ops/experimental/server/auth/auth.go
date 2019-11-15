
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Started from https://raw.githubusercontent.com/jordan-wright/gophish/master/auth/auth.go

package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	api "github.com/apache/trafficcontrol/traffic_ops/experimental/server/api"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	ctx "github.com/gorilla/context"
)

type loginJson struct {
	U string `json:"u"`
	P string `json:"p"`
}

type SessionUser struct {
	User string
	Role int64
}

type TokenResponse struct {
	Token string
}

// to get a token:
//  curl --header "Content-Type:application/json" -XPOST http://host:port/login -d'{"u":"yourusername", "p":"yourpassword}'

func validateToken(tokenString string) (*jwt.Token, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("mySigningKey"), nil // JvD
	})

	if err == nil && token.Valid {
		log.Println("TOKEN IS GOOD -- user:", token.Claims["userid"], " role:", token.Claims["role"])
	} else {
		log.Println("TOKEN IS BAD", err)
	}
	return token, err
}

// GetContext wraps each request in a function which fills in the context for a given request.
// This includes setting the User and Session keys and values as necessary for use in later functions.
func GetContext(handler http.Handler) http.HandlerFunc {
	// Set the context here

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := validateToken(r.Header.Get("Authorization"))
		if err != nil {
			log.Println("No valid token found!")
		} else {
			ctx.Set(r, "user", token.Claims["userid"])
			ctx.Set(r, "role", token.Claims["role"])
		}
		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}

// GetLoginOptionsFunc returns a func which handles the OPTIONS request for the login endpoint.
func GetLoginOptionsFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// GetLoginFunc returns a func which attempts to login the user given a request.
// Only works for local password at this time.
func GetLoginFunc(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := ""
		password := ""
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body: ", err.Error())
			http.Error(w, "Error reading body: "+err.Error(), http.StatusBadRequest)
			return
		}
		var lj loginJson
		log.Println(body)
		err = json.Unmarshal(body, &lj)
		if err != nil {
			log.Println("Error unmarshalling JSON: ", err.Error())
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		username = lj.U
		password = lj.P
		userInterface, err := api.GetUser(username, db)
		if err != nil {
			http.Error(w, "Invalid user: "+err.Error(), http.StatusUnauthorized)
			return
		}
		u, ok := userInterface.(api.Users)
		if !ok {
			http.Error(w, "Error GetUser returned a non-user.", http.StatusInternalServerError)
			return
		}

		encBytes := sha1.Sum([]byte(password))
		encString := hex.EncodeToString(encBytes[:])
		if err != nil {
			ctx.Set(r, "user", nil)
			log.Println("Invalid password")
			http.Error(w, "Invalid password: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if u.LocalPassword.String != encString {
			ctx.Set(r, "user", nil)
			log.Println("Invalid password")
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Create the token
		token := jwt.New(jwt.SigningMethodHS256)
		// Set some claims
		token.Claims["userid"] = u.Username
		token.Claims["role"] = u.Links.RolesLink.ID
		token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte("mySigningKey")) // TODO JvD
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
	}
}

// Logout destroys the current user session
func Logout(w http.ResponseWriter, r *http.Request) {
	// TODO JvD: revoke the token?
	http.Redirect(w, r, "/login", http.StatusFound)
}

func DONTRequireLogin(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

// RequireLogin is a simple middleware which checks to see if the user is currently logged in.
// If not, the function returns a 302 redirect to the login page.
func RequireLogin(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := ctx.Get(r, "user")
		role := ctx.Get(r, "role")
		if user != nil {
			log.Println("userId:", user, " userRole:", role)
			handler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

// Use allows us to stack middleware to process the request
// Example taken from https://github.com/gorilla/mux/pull/36#issuecomment-25849172
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
