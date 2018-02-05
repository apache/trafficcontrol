package startup

import (
	"github.com/apache/incubator-trafficcontrol/grove/config"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, hello)
}

func hello(cfg config.Config) {
	log.Errorf("Hello World! I'm a startup plugin! We're starting with %v bytes!\n", cfg.CacheSizeBytes)
}
