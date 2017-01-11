package influxdb

import (
	"flag"
	"fmt"
	"strings"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

// Config holds the requisite information to connect to an influx db instance
// prefix is an optional variable that is set via the Flags method, so as to
// distinguish between flags for different influxdb instances
type Config struct {
	prefix   string
	User     string
	Password string
	URL      string
}

// Flags configures the cli flags for the config. If more than one config is
// present for a program, it should be differentiated by a different prefix, so
// that the flag names don't collide
func (c *Config) Flags(prefix string) {
	c.prefix = prefix
	flag.StringVar(&c.URL, flagName(c.prefix, "url"), "http://localhost:8086", "The influxdb url and port")
	flag.StringVar(&c.User, flagName(c.prefix, "user"), "", "The influxdb username to connect to the db")
	flag.StringVar(&c.Password, flagName(c.prefix, "password"), "", "The influxdb password to connect to the db")
}

// NewHTTPClient tries to use the given configuration to
func (c *Config) NewHTTPClient() (influx.Client, error) {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     c.URL,
		Username: c.User,
		Password: c.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error creating influx client")
	}
	_, _, err = client.Ping(10)
	return client, errors.Wrap(err, "Error creating influx client")
}

func flagName(prefix, name string) string {
	if prefix != "" {
		prefix = fmt.Sprintf("%s-", prefix)
	}
	return strings.ToLower(fmt.Sprintf("%s%s", prefix, name))
}
