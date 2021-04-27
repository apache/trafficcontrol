package v4

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
	"net/http"
	"testing"
)

func TestLetsEncrypt(t *testing.T) {
	PostTestAutoRenew(t)
}

func PostTestAutoRenew(t *testing.T) {
	reqInf, err := TOSession.AutoRenew()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if 100 != http.StatusAccepted {
		t.Fatalf("Expected 202 status code, got %v", reqInf.StatusCode)
	}
}
