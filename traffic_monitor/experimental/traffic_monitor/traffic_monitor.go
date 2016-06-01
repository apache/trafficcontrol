package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	_ "github.com/Comcast/traffic_control/traffic_monitor/experimental/common/instrumentation"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/manager"
	_ "github.com/davecheney/gmx"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	opsConfigFile := flag.String("opsCfg", "", "The traffic ops config file")
	flag.Parse()

	if *opsConfigFile == "" {
		fmt.Println("The --opsCfg argument is required")
		os.Exit(1)
	}

	// Start the Manager
	manager.Start(*opsConfigFile)
}
