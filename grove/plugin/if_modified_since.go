package plugin

import (
	"net/http"

	"github.com/apache/incubator-trafficcontrol/grove/web"
)

func init() {
	AddPlugin(5000, Funcs{beforeRespond: ifModifiedSince})
}

func ifModifiedSince(icfg interface{}, d BeforeRespondData) {
	if d.CacheObj == nil {
		return // if we don't have a cacheobj from the origin, there's no object to have been modified.
	}
	modifiedSince, ok := web.GetHTTPDate(d.Req.Header, "If-Modified-Since")
	if !ok {
		return
	}
	if d.CacheObj.LastModified.After(modifiedSince) {
		return
	}
	*d.Code, *d.Hdr, *d.Body = http.StatusNotModified, nil, nil
}
