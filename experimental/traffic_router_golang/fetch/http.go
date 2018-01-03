package fetch

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type httpFetcher struct {
	url       string
	timeout   time.Duration
	userAgent string
}

func NewHTTP(url string, timeout time.Duration, userAgent string) Fetcher {
	return httpFetcher{url: url, timeout: timeout}
}

func (f httpFetcher) Fetch() ([]byte, error) {
	client := http.Client{
		Timeout: f.timeout,
	}
	req, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		// TODO round-robin retry on error?
		return nil, errors.New("HTTP creating request '" + f.url + "': " + err.Error())
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := client.Do(req)
	if err != nil {
		// TODO round-robin retry on error?
		return nil, errors.New("HTTP request '" + f.url + "': " + err.Error())
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("HTTP reading response body '" + f.url + "': " + err.Error())
	}
	return b, nil
}
