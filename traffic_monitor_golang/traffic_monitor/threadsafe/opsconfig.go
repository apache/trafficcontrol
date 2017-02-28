package threadsafe

import (
	"sync"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/handler"
)

// OpsConfig provides safe access for multiple reader goroutines and a single writer to a stored OpsConfig object.
// This could be made lock-free, if the performance was necessary
type OpsConfig struct {
	opsConfig *handler.OpsConfig
	m         *sync.RWMutex
}

// NewOpsConfig returns a new single-writer-multiple-reader OpsConfig
func NewOpsConfig() OpsConfig {
	return OpsConfig{m: &sync.RWMutex{}, opsConfig: &handler.OpsConfig{}}
}

// Get gets the internal OpsConfig object. This MUST NOT be modified. If modification is necessary, copy the object.
func (o *OpsConfig) Get() handler.OpsConfig {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.opsConfig
}

// Set sets the internal OpsConfig object. This MUST NOT be called from multiple goroutines.
func (o *OpsConfig) Set(newOpsConfig handler.OpsConfig) {
	o.m.Lock()
	*o.opsConfig = newOpsConfig
	o.m.Unlock()
}
