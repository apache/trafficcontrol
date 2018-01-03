package fetch

import (
	"sync/atomic"
	"time"
)

type fetcher struct {
	fetchers []Fetcher
	pos      *uint64
}

func NewRoundRobin(fetchers []Fetcher) Fetcher {
	p := uint64(0)
	return fetcher{fetchers: fetchers, pos: &p}
}

func (f fetcher) Fetch() ([]byte, error) {
	nextI := atomic.AddUint64(f.pos, 1)
	fetcher := f.fetchers[nextI%uint64(len(f.fetchers))]

	return fetcher.Fetch() // TODO round-robin retry on error?
}

func NewHTTPRoundRobin(hosts []string, path string, timeout time.Duration, userAgent string) Fetcher {
	fetchers := []Fetcher{}
	for _, host := range hosts {
		fetchers = append(fetchers, NewHTTP(host+path, timeout, userAgent))
	}
	return NewRoundRobin(fetchers)
}

func NewFileRoundRobin(paths []string, timeout time.Duration, userAgent string) Fetcher {
	fetchers := []Fetcher{}
	for _, path := range paths {
		fetchers = append(fetchers, NewFile(path))
	}
	return NewRoundRobin(fetchers)
}
