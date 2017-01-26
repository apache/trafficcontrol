package gmx

// pkg/runtime instrumentation

import "runtime"

var memstats runtime.MemStats

func init() {
	Publish("runtime.gomaxprocs", runtimeGOMAXPROCS)
	Publish("runtime.numcgocall", runtimeNumCgoCall)
	Publish("runtime.numcpu", runtimeNumCPU)
	Publish("runtime.numgoroutine", runtimeNumGoroutine)
	Publish("runtime.version", runtimeVersion)

	Publish("runtime.memstats", runtimeMemStats)
}

func runtimeGOMAXPROCS() interface{} {
	return runtime.GOMAXPROCS(0)
}

func runtimeNumCgoCall() interface{} {
	return runtime.NumCgoCall()
}

func runtimeNumCPU() interface{} {
	return runtime.NumCPU()
}

func runtimeNumGoroutine() interface{} {
	return runtime.NumGoroutine()
}

func runtimeVersion() interface{} {
	return runtime.Version()
}

func runtimeMemStats() interface{} {
	runtime.ReadMemStats(&memstats)
	return memstats
}
