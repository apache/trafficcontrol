package web

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
		CopyHeader(t.realWriter.Header(), &t.WrittenHeader) // if we didn't copy, someone could call Write() then change the headers, and we'd get the changed but unwritten headers.
		t.Code = http.StatusOK                              // emulate http.ResponseWriter.Write behavior
	}
	return wrote, err
}

func (t *HTTPResponseWriterTee) WriteHeader(code int) {
	t.Code = code
	t.realWriter.WriteHeader(code)
	CopyHeader(t.realWriter.Header(), &t.WrittenHeader) // if we didn't copy, someone could call Write() then change the headers, and we'd get the changed but unwritten headers.
}
