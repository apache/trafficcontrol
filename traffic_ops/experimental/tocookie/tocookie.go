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
	"strings"
	"time"
)

const Name = "mojolicious"

type Cookie struct {
	AuthData    string `json:"auth_data"`
	ExpiresUnix int64  `json:"expires"`
}

func checkHmac(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func Parse(secret, cookie string) (*Cookie, error) {
	dashPos := strings.Index(cookie, "-")
	if dashPos == -1 {
		return nil, fmt.Errorf("malformed cookie '%s' - no dashes", cookie)
	}
	if len(cookie) < dashPos+4 {
		return nil, fmt.Errorf("malformed cookie '%s' - no signature", cookie)
	}

	base64Txt := cookie[:dashPos]

	txtBytes, err := base64.RawURLEncoding.DecodeString(base64Txt)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 data: %v", err)
	}

	base64Sig := cookie[dashPos+4:]
	sigBytes, err := hex.DecodeString(base64Sig)
	if err != nil {
		return nil, fmt.Errorf("error decoding signature: %v", err)
	}

	if !checkHmac([]byte(base64Txt+"--"), sigBytes, []byte(secret)) {
		return nil, fmt.Errorf("bad signature")
	}

	cookieData := Cookie{}
	if err := json.Unmarshal(txtBytes, &cookieData); err != nil {
		return nil, fmt.Errorf("error decoding base64 text '%s' to JSON: %v", string(txtBytes), err)
	}

	if cookieData.ExpiresUnix-time.Now().Unix() < 0 {
		return nil, fmt.Errorf("signature expired")
	}

	return &cookieData, nil
}
