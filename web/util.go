package web

import (
	"net/http"
)

func CopyHeaderTo(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func CopyHeader(source http.Header) http.Header {
	dest := http.Header{}
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
	return dest
}
