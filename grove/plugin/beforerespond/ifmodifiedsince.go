package beforerespond

import (
	"net/http"

	"github.com/apache/incubator-trafficcontrol/grove/web"
)

const IfModifiedSinceName = "if_modified_since"

func init() {
	AddPlugin(5000, IfModifiedSinceName, ifModifiedSince, nil)
}

func ifModifiedSince(icfg interface{}, d Data) {
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
