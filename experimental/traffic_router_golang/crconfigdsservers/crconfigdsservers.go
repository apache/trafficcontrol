package crconfigdsservers

type DSServers map[tc.DeliveryServiceName]DSServer

type DSServer struct {
	// TODO handle Steering
	Name    string
	DSFQDNs []string // TOOD determine if neccessary, how to index
}

type DSServer struct {
	DirectMatches                      map[string]tc.DeliveryServiceName
	DotStartSlashDotFooSlashDotDotStar map[string]tc.DeliveryServiceName
	RegexMatch                         map[*regexp.Regexp]tc.DeliveryServiceName
}
