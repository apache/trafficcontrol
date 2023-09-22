package urisigning

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
	"bytes"
	"encoding/json"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

const (
	goodKeyset = `
{
  "Kabletown URI Authority 1": {
		"renewal_kid": "Second Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Second Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  },
  "Kabletown URI Authority 2": {
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Third Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
	badJSONKeySet = `

  "Kabletown URI Authority 1": {
		"renewal_kid": "Second Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Second Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
	noRenewalKidKeyset = `
{
  "Kabletown URI Authority 1": {
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Second Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
	noMatchingRenewalKidKeyset = `
{
  "Kabletown URI Authority 1": {
		"renewal_kid": "Second Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Other Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
)

func TestValidateURIKeyset(t *testing.T) {
	var keyset tc.JWKSMap

	// unmarshal a good URISignerKeyset
	if err := json.Unmarshal([]byte(goodKeyset), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	// now validate it
	if err := validateURIKeyset(keyset); err != nil {
		t.Errorf("validateURIKeyset(): expected nil error, actual: %v", err)
	}

	// unmarshal a bad URISignerKeySet
	if err := json.Unmarshal([]byte(badJSONKeySet), &keyset); err == nil {
		t.Errorf("json.UnMarshal(): expected an error")
	}

	// unmarshal a good URISignerKeyset with a missing renewal_kid
	if err := json.Unmarshal([]byte(noRenewalKidKeyset), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	// now validate it, expect an erro due to missing renewal_kid
	if err := validateURIKeyset(keyset); err == nil {
		t.Errorf("validateURIKeyset(): expected an error")
	}

	// unmarshal a good URISignerKeyset with no matching kid to the renewal_kid
	if err := json.Unmarshal([]byte(noMatchingRenewalKidKeyset), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	// now validate it, expect an erro due to missing a matching key kid to renewal_kid
	if err := validateURIKeyset(keyset); err == nil {
		t.Errorf("validateURIKeyset(): expected an error")
	}
}

func TestEmptyURISigningKey(t *testing.T) {
	buf, err := json.Marshal(emptyURISigningKey)
	if err != nil {
		t.Errorf(`failed to marshal empty URISingingKey`)
	}
	if !bytes.Equal([]byte(`{"renewal_kid":null,"keys":null}`), buf) {
		t.Errorf(`unexpected JSON: %q`, buf)
	}
}
