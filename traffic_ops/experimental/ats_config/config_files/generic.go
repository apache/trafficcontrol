package config_files

import (
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"regexp"
	"time"
)

// my $separator ||= {
// 	"records.config"  => " ",
// 	"plugin.config"   => " ",
// 	"sysctl.conf"     => " = ",
// 	"url_sig_.config" => " = ",
// 	"astats.config"   => "=",
// };

func createGenericDotConfigFunc(separator string) ConfigFileCreatorFunc {
	underscoreDigitSuffixRegex := regexp.MustCompile("__[0-9]+$")

	return func(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
		s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

		paramMap := createParamsMap(params)
		fileParams := paramMap[filename]
		for name, val := range fileParams {
			name := underscoreDigitSuffixRegex.ReplaceAllString(name, "")
			s += name + separator + val + "\n"
		}
		return s, nil
	}
}
