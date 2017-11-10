package fetch

import (
	"errors"
	"io/ioutil"
)

type fileFetcher struct {
	path string
}

func NewFile(path string) Fetcher {
	return fileFetcher{path: path}
}

func (f fileFetcher) Fetch() ([]byte, error) {
	b, err := ioutil.ReadFile(f.path)
	if err != nil {
		// TODO round-robin retry on error?
		return nil, errors.New("reading file: " + err.Error())
	}
	return b, nil
}
