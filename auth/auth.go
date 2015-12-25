// Copyright 2015 Comcast Cable Communications Management, LLC

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
	"../todb"

	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"

	ctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Flash struct {
	Type    string
	Message string
}

const loginPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="username" name="username">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))

func LoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, loginPage)
}

// GetContext wraps each request in a function which fills in the context for a given request.
// This includes setting the User and Session keys and values as necessary for use in later functions.
func GetContext(handler http.Handler) http.HandlerFunc {
	// Set the context here
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing request", http.StatusInternalServerError)
		}
		// Set the context appropriately here.
		// Set the session
		session, _ := Store.Get(r, "trafficOps")
		// Put the session in the context so that
		ctx.Set(r, "session", session)
		if id, ok := session.Values["id"]; ok {
			fmt.Println("userid ", id)
			if err != nil {
				ctx.Set(r, "user", nil)
			} else {
				ctx.Set(r, "user", id)
			}
		} else {
			ctx.Set(r, "user", nil)
		}
		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}

// Login attempts to login the user given a request. Only works for local passwd at this time
func Login(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	session, _ := Store.Get(r, "trafficOps")
	u, err := todb.GetTmUser(username)
	redirectTarget := "/hello/jvd"
	encBytes := sha1.Sum([]byte(password))
	encString := hex.EncodeToString(encBytes[:])
	// fmt.Println("sha1:", hex.EncodeToString(encBytes[:]), " localpasswd:", u.LocalPasswd.String, "err:", err)
	if err != nil || u.LocalPasswd.String != encString {
		ctx.Set(r, "user", nil)
		fmt.Println("Invalid passwd")
		redirectTarget = "/login"
		delete(session.Values, "id")
	} else {
		ctx.Set(r, "user", u)
		session.Values["id"] = u.Id
	}
	session.Save(r, w)
	http.Redirect(w, r, redirectTarget, 302)
}

// Logout destroys the current user session
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "trafficOps")
	delete(session.Values, "id")
	http.Redirect(w, r, "/login", 302)
}

// RequireLogin is a simple middleware which checks to see if the user is currently logged in.
// If not, the function returns a 302 redirect to the login page.
func RequireLogin(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := ctx.Get(r, "user")
		fmt.Println(">", u)
		if u != nil {
			handler.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
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
