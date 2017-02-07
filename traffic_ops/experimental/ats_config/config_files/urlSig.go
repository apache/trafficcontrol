package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"time"
)

// TODO pass riak creds to all funcs? Change all funcs to use lambdas?
func createUrlSigDotConfigFunc(riakUser string, riakPass string) ConfigFileCreatorFunc {
	return func(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {

		separator := " = "

		s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

		insecure := true // TODO fix riak certs, and change to false
		keys, err := RiakGetURLSigKeys(toClient, riakUser, riakPass, filename, insecure)
		if err != nil {
			return "", fmt.Errorf("error getting keys from Riak: %v", err)
		}

		for keyName, keyVal := range keys {
			s += fmt.Sprintf("%s%s%s\n", keyName, separator, keyVal)
		}

		return s, nil
	}
}
