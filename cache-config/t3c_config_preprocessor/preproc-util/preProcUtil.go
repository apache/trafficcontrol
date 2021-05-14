package preproc_util

import (
	"encoding/json"
	"errors"
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"io"
	"sort"
)

func WriteConfigs(configs []t3cutil.ATSConfigFile, output io.Writer) error {
	sort.Sort(config.ATSConfigFiles(configs))
	if err := json.NewEncoder(output).Encode(configs); err != nil {
		return errors.New("encoding and writing configs: " + err.Error())
	}
	return nil
}
