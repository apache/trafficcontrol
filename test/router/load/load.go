package load

import (
	"strings"
	"sync"

	"fmt"
	"github.com/Comcast/traffic_control/test/router/client"
	"github.com/Comcast/traffic_control/test/router/data"
)

type LoadTest struct {
	CaFile                string   `json:"caFile"`
	Cdn                   string   `json:"cdn"`
	TxCount               int      `json:"txCount"`
	Connections           int      `json:"connections"`
	HttpDeliveryServices  []string `json:"httpDeliveryServices"`
	HttpsDeliveryServices []string `json:"httpsDeliveryServices"`
}

func DoLoadTest(loadtest LoadTest, done chan struct{}) chan data.HttpResult {
	resultsChan := make(chan data.HttpResult)

	go func() {
		fmt.Println("Starting load test", loadtest)
		defer close(done)
		var waitGroup sync.WaitGroup
		for _, deliveryService := range loadtest.HttpDeliveryServices {
			waitGroup.Add(1)
			go func(ds string) {
				defer waitGroup.Done()
				host := strings.Join([]string{ds, loadtest.Cdn}, ".")
				tlsConfig := client.MustGetTlsConfiguration(host, loadtest.CaFile)
				client.ExerciseDeliveryService(false, tlsConfig, host, loadtest.TxCount, loadtest.Connections, resultsChan)
			}(deliveryService)
		}

		for _, deliveryService := range loadtest.HttpsDeliveryServices {
			waitGroup.Add(1)
			go func(ds string) {
				defer waitGroup.Done()
				host := strings.Join([]string{ds, loadtest.Cdn}, ".")
				tlsConfig := client.MustGetTlsConfiguration(host, loadtest.CaFile)
				client.ExerciseDeliveryService(true, tlsConfig, host, loadtest.TxCount, loadtest.Connections, resultsChan)
			}(deliveryService)
		}

		waitGroup.Wait()
	}()

	return resultsChan
}
