package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tocookie"

	"github.com/jmoiron/sqlx"
)

type passwordForm struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

func LoginHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		defer r.Body.Close()
		form := passwordForm{}
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			handleErrs(http.StatusBadRequest, err)
			return
		}
		authenticated, err := checkLocalUser(form, db)
		if err != nil {
			log.Errorf("error checking local user: %s\n", err.Error())
		}
		var ldapErr error
		if !authenticated {
			if cfg.LDAPEnabled {
				authenticated, ldapErr = checkLDAPUser(form, cfg.ConfigLDAP)
				if ldapErr != nil {
					log.Errorf("error checking ldap user: %s\n", ldapErr.Error())
				}
			}
		}
		resp := struct {
			tc.Alerts
		}{}
		if authenticated {
			expiry := time.Now().Add(time.Hour * 6)
			cookie := tocookie.New(form.Username, expiry, cfg.Secrets[0])
			httpCookie := http.Cookie{Name: "mojolicious", Value: cookie, Path: "/", Expires: expiry, HttpOnly: true}
			http.SetCookie(w, &httpCookie)
			resp = struct {
				tc.Alerts
			}{tc.CreateAlerts(tc.SuccessLevel, "Successfully logged in.")}

		} else {
			resp = struct {
				tc.Alerts
			}{tc.CreateAlerts(tc.ErrorLevel, "Invalid username or password.")}
		}
		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

func checkLocalUser(form passwordForm, db *sqlx.DB) (bool, error) {
	var hashedPassword string
	err := db.Get(&hashedPassword, "SELECT local_passwd FROM tm_user WHERE username=$1", form.Username)
	if err != nil {
		return false, err
	}
	err = VerifyPassword(form.Password, hashedPassword)
	if err != nil {
		if hashedPassword == sha1Hex(form.Password) {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

func sha1Hex(s string) string {
	// SHA1 hash
	hash := sha1.New()
	hash.Write([]byte(s))
	hashBytes := hash.Sum(nil)

	// Hexadecimal conversion
	hexSha1 := hex.EncodeToString(hashBytes)
	return hexSha1
}

func checkLDAPUser(form passwordForm, cfg *config.ConfigLDAP) (bool, error) {
	userDN, valid, err := LookupUserDN(form.Username, cfg)
	if err != nil {
		return false, err
	}
	if valid {
		return AuthenticateUserDN(userDN, form.Password, cfg)
	}
	return false, errors.New("User not found in LDAP")
}
