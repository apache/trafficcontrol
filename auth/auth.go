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
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	api "github.com/Comcast/traffic_control/traffic_ops/goto2/api"
	output "github.com/Comcast/traffic_control/traffic_ops/goto2/output_format"
	// jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strings"

	ctx "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// type Flash struct {
// 	Type    string
// 	Message string
// }

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

type loginJson struct {
	U string `json:"u"`
	P string `json:"p"`
}

type SessionUser struct {
	UserId int64
	Role   int64
}

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
		session, err := Store.Get(r, "trafficOps")
		if err != nil {
			fmt.Println(err)
		}
		// Put the session in the context so that
		ctx.Set(r, "session", session)
		val := session.Values["user"]
		if user, ok := val.(SessionUser); !ok {
			fmt.Println("No valid SessionUser found!")
		} else {
			ctx.Set(r, "user", user)
		}
		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}

// Login attempts to login the user given a request. Only works for local passwd at this time
func Login(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Starting Login...")
	username := ""
	password := ""
	htmlSession := true
	if r.FormValue("username") != "" {
		username, password = r.FormValue("username"), r.FormValue("password")
	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic("booboo")
		}
		var lj loginJson
		err = json.Unmarshal(body, &lj)
		if err != nil {
			panic("boo")
		}
		username = lj.U
		password = lj.P
		htmlSession = false
	}
	session, _ := Store.Get(r, "trafficOps")
	u := api.TmUser{}
	u, err := api.GetTmUserByName(username)
	redirectTarget := "/"
	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, flashMsg := range flashes {
			if strings.HasPrefix(flashMsg.(string), "pathDenied:") {
				redirectTarget = strings.Replace(flashMsg.(string), "pathDenied:", "", 1)
			}
		}
	}
	encBytes := sha1.Sum([]byte(password))
	encString := hex.EncodeToString(encBytes[:])
	// fmt.Println("sha1:", hex.EncodeToString(encBytes[:]), " localpasswd:", u.LocalPasswd.String, "err:", err)
	if err != nil || u.LocalPasswd.String != encString {
		ctx.Set(r, "user", nil)
		fmt.Println("Invalid passwd")
		redirectTarget = "/login"
		delete(session.Values, "id")
	} else {
		session.Values["user"] = SessionUser{UserId: u.Id, Role: u.Role.Int64}
	}
	session.Save(r, w)
	if htmlSession {
		http.Redirect(w, r, redirectTarget, 302)
	} else {
		respTxt := output.MakeApiResponse(nil, output.MakeAlert("Successfully logged in.", "success"), nil)
		js, err := json.Marshal(respTxt)
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
	session, _ := Store.Get(r, "trafficOps")
	delete(session.Values, "id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
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
		u := ctx.Get(r, "user")
		if u != nil {
			fmt.Println("userId:", u.(SessionUser).UserId, " userRole:", u.(SessionUser).Role)
			handler.ServeHTTP(w, r)
		} else {
			session, _ := Store.Get(r, "trafficOps")
			session.AddFlash("pathDenied:" + r.URL.EscapedPath())
			session.Save(r, w)
			http.Redirect(w, r, "/login", 302)
		}
	}
}

// func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
// 	authBackend := InitJWTAuthenticationBackend()

// 	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		} else {
// 			return authBackend.PublicKey, nil
// 		}
// 	})

// 	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
// 		next(rw, req)
// 	} else {
// 		rw.WriteHeader(http.StatusUnauthorized)
// 	}
// }

// Use allows us to stack middleware to process the request
// Example taken from https://github.com/gorilla/mux/pull/36#issuecomment-25849172
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
