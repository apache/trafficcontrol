package gmx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

const GMX_VERSION = 0

var (
	r = &registry{
		entries: make(map[string]func() interface{}),
	}
)

func init() {
	s, err := localSocket()
	if err != nil {
		log.Printf("gmx: unable to open local socket: %v", err)
		return
	}

	// register the registries keys for discovery
	Publish("keys", func() interface{} {
		return r.keys()
	})
	go serve(s, r)
}

func localSocket() (net.Listener, error) {
	return net.ListenUnix("unix", localSocketAddr())
}

func localSocketAddr() *net.UnixAddr {
	return &net.UnixAddr{
		filepath.Join(os.TempDir(), fmt.Sprintf(".gmx.%d.%d", os.Getpid(), GMX_VERSION)),
		"unix",
	}
}

// Publish registers the function f with the supplied key.
func Publish(key string, f func() interface{}) {
	r.register(key, f)
}

func serve(l net.Listener, r *registry) {
	// if listener is a unix socket, try to delete it on shutdown
	if l, ok := l.(*net.UnixListener); ok {
		if a, ok := l.Addr().(*net.UnixAddr); ok {
			defer os.Remove(a.Name)
		}
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			return
		}
		go handle(c, r)
	}
}

func handle(nc net.Conn, reg *registry) {
	// conn makes it easier to send and receive json
	type conn struct {
		net.Conn
		*json.Encoder
		*json.Decoder
	}
	c := conn{
		nc,
		json.NewEncoder(nc),
		json.NewDecoder(nc),
	}
	defer c.Close()
	for {
		var keys []string
		if err := c.Decode(&keys); err != nil {
			if err != io.EOF {
				log.Printf("gmx: client %v sent invalid json request: %v", c.RemoteAddr(), err)
			}
			return
		}
		var result = make(map[string]interface{})
		for _, key := range keys {
			if f, ok := reg.value(key); ok {
				// invoke the function for key and store the result
				result[key] = f()
			}
		}
		if err := c.Encode(result); err != nil {
			log.Printf("gmx: could not send response to client %v: %v", c.RemoteAddr(), err)
			return
		}
	}
}

type registry struct {
	sync.Mutex // protects entries from concurrent mutation
	entries    map[string]func() interface{}
}

func (r *registry) register(key string, f func() interface{}) {
	r.Lock()
	defer r.Unlock()
	r.entries[key] = f
}

func (r *registry) value(key string) (func() interface{}, bool) {
	r.Lock()
	defer r.Unlock()
	f, ok := r.entries[key]
	return f, ok
}

func (r *registry) keys() (k []string) {
	r.Lock()
	defer r.Unlock()
	for e := range r.entries {
		k = append(k, e)
	}
	return
}
