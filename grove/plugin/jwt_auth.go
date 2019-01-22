package plugin

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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/apache/trafficcontrol/lib/go-log"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/iancoleman/strcase"
)

type JwtSecret struct {
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Algoritm string `json:"alg"`
	Kid      string `json:"kid"`
	Value    string `json:"value"`
}

type JwtConfig struct {
	ExportClaims []string    `json:"export_claims"`
	Secrets      []JwtSecret `json:"secrets"`
}

type JwtContext struct {
	SecretsIndex map[string]interface{}
}

func init() {
	AddPlugin(10000, Funcs{onRequest: jwtAuth, startup: jwtStartup, load: jwtLoad})
}

func jwtAuth(icfg interface{}, d OnRequestData) bool {
	jwtContext := (*d.Context).(JwtContext)
	jwtConfig := icfg.(*JwtConfig)
	tokenOnRequest, err := request.OAuth2Extractor.ExtractToken(d.R)
	if err != nil {
		d.W.WriteHeader(401)
		d.W.Write([]byte(err.Error()))
		return true
	}
	parser := &jwt.Parser{}
	token, err := parser.ParseWithClaims(tokenOnRequest, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*jwt.MapClaims)
		var kid, iss, aud string
		if value, ok := (*claims)["kid"].(string); ok {
			kid = value
		}
		if value, ok := (*claims)["iss"].(string); ok {
			iss = value
		}
		if value, ok := (*claims)["aud"].(string); ok {
			aud = value
		}

		if secret := jwtContext.SecretsIndex[getSecretIndexKey(kid, iss, aud)]; secret != nil {
			return secret, nil
		}
		return nil, fmt.Errorf("Unknown key")
	})
	if err != nil {
		d.W.WriteHeader(401)
		d.W.Write([]byte(err.Error()))
		return true
	}
	if !token.Valid {
		d.W.WriteHeader(401)
		d.W.Write([]byte("Invalid token"))
		return true
	}
	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		for _, claimName := range jwtConfig.ExportClaims {
			if value := (*claims)[claimName]; value != nil {
				// What if value is not string ?
				d.R.Header.Set(fmt.Sprintf("X-Claim-%s", strcase.ToCamel(claimName)), value.(string))
			}
		}
	}
	return false
}

func getSecretIndexKey(kid, iss, aud string) string {
	if kid != "" {
		return kid
	}
	s := ""
	if iss != "" {
		s += iss
	}
	if aud != "" {
		s += ":" + aud
	}
	return s
}

func indexSecrets(jwtConfig *JwtConfig, jwtContext *JwtContext) error {
	for _, secret := range jwtConfig.Secrets {
		indexKey := getSecretIndexKey(secret.Kid, secret.Issuer, secret.Audience)
		log.Debugf("JWT loading secret for indexKey=%v", indexKey)
		switch secret.Algoritm {
		case "rs256":
			pem, err := ioutil.ReadFile(secret.Value)
			if err != nil {
				return err
			}
			pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
			if err != nil {
				return err
			}
			jwtContext.SecretsIndex[indexKey] = pubKey
		case "hs256":
			jwtContext.SecretsIndex[indexKey] = []byte(secret.Value)
		default:
			return fmt.Errorf("JWT %v algorithm is not supported", secret.Algoritm)
		}
	}
	return nil
}

func jwtStartup(icfg interface{}, d StartupData) {
	jwtConfig := icfg.(*JwtConfig)
	jwtContext := JwtContext{}
	jwtContext.SecretsIndex = make(map[string]interface{})
	err := indexSecrets(jwtConfig, &jwtContext)
	if err != nil {
		log.Errorln("JWT can't index secrets: " + err.Error())
	}
	*d.Context = jwtContext
	log.Debugf("JWT startup success")
}

func jwtLoad(b json.RawMessage) interface{} {
	jwtConfig := JwtConfig{}
	err := json.Unmarshal(b, &jwtConfig)
	if err != nil {
		log.Errorln("JWT loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf("JWT config load success")
	return &jwtConfig
}
