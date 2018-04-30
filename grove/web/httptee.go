package web

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
)

func NewHTTPResponseWriterTee(w http.ResponseWriter) *HTTPResponseWriterTee {
	return &HTTPResponseWriterTee{realWriter: w, WrittenHeader: http.Header{}}
}

// HTTPResponseWriterTee fulfills `http.ResponseWriter`
type HTTPResponseWriterTee struct {
	realWriter http.ResponseWriter
	// headers       http.Header
	Bytes         []byte
	Code          int
	WrittenHeader http.Header
}

func (t *HTTPResponseWriterTee) Header() http.Header {
	return t.realWriter.Header()
}

func (t *HTTPResponseWriterTee) Write(b []byte) (int, error) {
	t.Bytes = append(t.Bytes, b...)
	wrote, err := t.realWriter.Write(b)
	// copy, because after Write is called, no more headers are written to the real connection.
	if t.Code == 0 {
		CopyHeaderTo(t.realWriter.Header(), &t.WrittenHeader) // if we didn't copy, someone could call Write() then change the headers, and we'd get the changed but unwritten headers.
		t.Code = http.StatusOK                                // emulate http.ResponseWriter.Write behavior
	}
	return wrote, err
}

func (t *HTTPResponseWriterTee) WriteHeader(code int) {
	t.Code = code
	t.realWriter.WriteHeader(code)
	CopyHeaderTo(t.realWriter.Header(), &t.WrittenHeader) // if we didn't copy, someone could call Write() then change the headers, and we'd get the changed but unwritten headers.
}
