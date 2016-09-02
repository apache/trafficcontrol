package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	_ "github.com/Comcast/traffic_control/traffic_monitor/experimental/common/instrumentation"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/manager"
	_ "github.com/davecheney/gmx"
)

var GitRevision = "No Git Revision Specified. Please build with '-X main.GitRevision=${git rev-parse HEAD}'"
var BuildTimestamp = "No Build Timestamp Specified. Please build with '-X main.BuildTimestamp=`date +'%Y-%M-%dT%H:%M:%S'`"

// getStaticAppData returns app data available at start time.
// This should be called immediately, as it includes calculating when the app was started.
func getStaticAppData() (manager.StaticAppData, error) {
	var d manager.StaticAppData
	d.StartTime = time.Now()
	d.GitRevision = GitRevision
	d.FreeMemoryMB = math.MaxUint64 // TODO remove if/when nothing needs this
	d.Version = Version
	wd, err := os.Getwd()
	if err != nil {
		return manager.StaticAppData{}, err
	}
	d.WorkingDir = wd
	d.Name = os.Args[0]
	d.BuildTimestamp = BuildTimestamp
	return d, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	staticData, err := getStaticAppData()
	if err != nil {
		fmt.Printf("ERROR failed to get static app data: %v", err)
		return
	}

	opsConfigFile := flag.String("opsCfg", "", "The traffic ops config file")
	flag.Parse()

	if *opsConfigFile == "" {
		fmt.Println("The --opsCfg argument is required")
		os.Exit(1)
	}

	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// Start the Manager
	manager.Start(*opsConfigFile, staticData)
}
