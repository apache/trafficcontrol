// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tocookie

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const GeneratedByStr = "trafficcontrol-go-tocookie"
const Name = "mojolicious"
const MojoCookie = "mojoCookie"
const AccessToken = "access_token"
const BearerToken = "Bearer"
const DefaultDuration = time.Hour

type Cookie struct {
	AuthData    string `json:"auth_data"`
	ExpiresUnix int64  `json:"expires"`
	By          string `json:"by"`
}

func checkHmac(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func Parse(secret, cookie string) (*Cookie, error, error) {
	if cookie == "" {
		return nil, nil, nil
	}

	dashPos := strings.Index(cookie, "-")
	if dashPos == -1 {
		return nil, fmt.Errorf("error parsing cookie: malformed cookie '%s' - no dashes", cookie), nil
	}

	lastDashPos := strings.LastIndex(cookie, "-")
	if lastDashPos == -1 {
		return nil, fmt.Errorf("error parsing cookie: malformed cookie '%s' - no dashes", cookie), nil
	}

	if len(cookie) < lastDashPos+1 {
		return nil, fmt.Errorf("error parsing cookie: malformed cookie '%s' -- no signature", cookie), nil
	}

	base64Txt := cookie[:dashPos]
	txtBytes, err := base64.RawURLEncoding.DecodeString(base64Txt)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing cookie: error decoding base64 data: %v", err)
	}
	base64TxtSig := cookie[:lastDashPos-1] // the signature signs the base64 including trailing hyphens, but the Go base64 decoder doesn't want the trailing hyphens.

	base64Sig := cookie[lastDashPos+1:]
	sigBytes, err := hex.DecodeString(base64Sig)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing cookie: error decoding signature: %v", err)
	}

	if !checkHmac([]byte(base64TxtSig), sigBytes, []byte(secret)) {
		return nil, fmt.Errorf("bad signature - unauthorized, please log in"), nil
	}

	cookieData := Cookie{}
	if err := json.Unmarshal(txtBytes, &cookieData); err != nil {
		return nil, nil, fmt.Errorf("error parsing cookie: error decoding base64 text '%s' to JSON: %v", string(txtBytes), err)
	}

	if cookieData.ExpiresUnix-time.Now().Unix() < 0 {
		return nil, fmt.Errorf("signature expired - unauthorized, please log in"), nil
	}

	return &cookieData, nil, nil
}

func NewRawMsg(msg, key []byte) string {
	base64Msg := base64.RawURLEncoding.EncodeToString(msg)
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base64Msg))
	encMac := mac.Sum(nil)
	base64Sig := hex.EncodeToString(encMac)
	return base64Msg + "--" + base64Sig
}

func GetCookie(authData string, duration time.Duration, secret string) *http.Cookie {
	expiry := time.Now().Add(duration)
	maxAge := int(duration.Seconds())
	c := Cookie{By: GeneratedByStr, AuthData: authData, ExpiresUnix: expiry.Unix()}
	m, _ := json.Marshal(c)
	msg := NewRawMsg(m, []byte(secret))
	httpCookie := http.Cookie{Name: "mojolicious", Value: msg, Path: "/", Expires: expiry, MaxAge: maxAge, HttpOnly: true}
	return &httpCookie
}
