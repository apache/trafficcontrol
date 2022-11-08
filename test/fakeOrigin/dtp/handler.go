package dtp

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const XDtpHdrStr = `X-Dtp`
const XDtpCcHdrStr = `X-Dtp-Cc`

type DTPHandler struct{}

func NewDTPHandler() DTPHandler {
	return DTPHandler{}
}

// This is for content generation
type Generator interface {
	ContentType() string
	io.ReadSeeker
}

type NewGeneratorFunc func(
	map[string]string,
	int64,
) Generator

var GlobalGeneratorFuncs = map[string]NewGeneratorFunc{}

// No content generation
type HandlerFunc func(
	http.ResponseWriter,
	*http.Request,
	map[string]string,
)

var GlobalHandlerFuncs = map[string]HandlerFunc{}

type NewForwardFunc func(
	http.ResponseWriter,
	*http.Request,
	Generator,
	map[string]string,
	int64,
) Generator

var GlobalForwarderFuncs = map[string]NewForwardFunc{}

/*
func timeoutAndDrop(w http.ResponseWriter, r *http.Request, durstr string) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		DebugLogf("Can't get a hijack\n")
		return
	}

	conn, buf, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		DebugLogf("hijack failed: %s\n", err)
		return
	}

	delaydur, err := time.ParseDuration(durstr)
	if err == nil {
		if 0 < delaydur {
			time.Sleep(delaydur)
		}
	}

	conn.Close()
}
*/

func (h DTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cmds []string

	/*
		if GlobalConfig.Drop {
			DebugLogf("Dropping request for '%s'\n", r.URL.String())
			timeoutAndDrop(w, r, "60s")
			return
		}
	*/
	if time.Duration(0) != GlobalConfig.StallDuration {
		DebugLogf("Stalling request for '%s'\n", r.URL.String())
		time.Sleep(GlobalConfig.StallDuration)
	}

	DebugLogf("Serving request for '%s'\n", r.URL.String())
	cmdFromStr := func(s string) string {
		DebugLogf("Parsing `%s`\n", s)
		begin, end := -1, -1
		for i, c := range s {
			if c == '~' {
				if begin == -1 {
					begin = i + 1
				} else {
					end = i
					break
				}
			}
		}
		if begin != -1 && end != -1 {
			return s[begin:end]
		}
		if begin != -1 {
			return s[begin:]
		}
		return ""
	}

	/* Check the URL */
	DebugLog(`Parsing path: ` + r.URL.EscapedPath())
	for _, part := range strings.Split(r.URL.EscapedPath(), `/`) {
		DebugLog("Parsing part: " + part)
		s, err := url.PathUnescape(part)
		if err != nil {
			continue
		}

		cmd := cmdFromStr(s)
		if cmd != `` {
			cmds = append(cmds, cmd)
		}
	}

	// Check the query
	for _, part := range strings.Split(r.URL.RawQuery, `&`) {
		s, err := url.PathUnescape(part)
		if err != nil {
			continue
		}

		cmd := cmdFromStr(s)
		if cmd != `` {
			cmds = append(cmds, cmd)
		}
	}

	// Check the headers
	for _, part := range r.Header[XDtpHdrStr] {
		for _, hdr := range strings.Split(part, `,`) {
			cmd := cmdFromStr(hdr)
			if cmd != `` {
				cmds = append(cmds, cmd)
			}
		}
	}

	// requests to map
	reqdat := make(map[string]string)

	// Check for special cache-control header
	for _, part := range r.Header[XDtpCcHdrStr] {
		reqdat[XDtpCcHdrStr] = part
	}

	for _, cmd := range cmds {
		var key string
		var val string
		dotind := strings.IndexByte(cmd, '.')
		if dotind <= 0 {
			key = cmd
		} else {
			key = cmd[:dotind]
			val = cmd[dotind+1:]
		}
		DebugLogf("Setting '%s' to '%s'\n", key, val)
		reqdat[key] = val
	}

	// check for connection handler, direct code or hijack
	if hcode, ok := reqdat[`h`]; ok {
		if hfunc, ok := GlobalHandlerFuncs[hcode]; ok {
			hfunc(w, r, reqdat)
			return
		}
	}

	// check for content generator, failover to text
	pcode := reqdat[`p`]
	genfunc, genok := GlobalGeneratorFuncs[pcode]
	if !genok {
		genfunc = GlobalGeneratorFuncs[`txt`]
	}

	// process headers
	if !ProcessHeaders(w, r, reqdat) {
		return
	}

	lastmod, _ := strconv.ParseInt(reqdat[`lm`], 10, 64)

	generator := genfunc(reqdat, lastmod)
	w.Header()[`Content-Type`] = []string{generator.ContentType()}

	// look for events, byte position events, handle extra args
	var forwarder NewForwardFunc = nil
	if fwdarg, ok := reqdat[`f`]; ok {
		fcode := fwdarg
		if find := strings.Index(fwdarg, `.`); 0 < find {
			fcode = fwdarg[0:find]
			args := fwdarg[find+1:]
			reqdat[fcode] = args
			DebugLogf("Adding args for forwarder %s: %s\n", fcode, args)
		}

		var ok = false
		if forwarder, ok = GlobalForwarderFuncs[fcode]; ok {
			DebugLogf("Processing a forwarder %s\n", fcode)
		}
	}

	if nil != forwarder {
		generator = forwarder(w, r, generator, reqdat, lastmod)
		if nil == generator {
			return
		}
	}

	http.ServeContent(w, r, "", time.Unix(lastmod, 0), generator)
}

func ProcessHeaders(w http.ResponseWriter, r *http.Request, reqdat map[string]string) bool {

	if _, ok := reqdat["cksum_req"]; ok {
		reqhdr := r.Header
		hdrstr := fmt.Sprintf("%v", reqhdr)
		hash := md5.Sum([]byte(hdrstr))
		hashstr := hex.EncodeToString(hash[:])
		w.Header().Set("X-Request-Header-Cksum", hashstr)
		DebugLogf("cksum header: '%s'", hashstr)
	}

	var lastmod int64 = 0
	var maxage int64 = 0

	now := time.Now().Unix()

	// handle last modified and cache control
	// note that 'lm' will always be available downstream
	if lmstr, ok := reqdat[`lm`]; ok {
		lmsec, err := strconv.ParseInt(lmstr, 10, 64)
		if nil == err && 0 < lmsec {
			lastmod = lmsec
		}
	} else if uistr, ok := reqdat[`ui`]; ok {
		uisec, err := strconv.ParseInt(uistr, 10, 64)
		if nil == err && 0 < uisec {
			lastmod = (now / uisec) * uisec
			maxage = now - lastmod // expire at quanta
			w.Header().Set(`Cache-Control`, fmt.Sprintf("max-age=%d", maxage))
			reqdat[`lm`] = strconv.FormatInt(lastmod, 10)
			reqdat[`ma`] = strconv.FormatInt(maxage, 10)
		}
	} else {
		reqdat[`lm`] = `0`
	}

	seed, _ := strconv.ParseInt(reqdat[`rnd`], 10, 64)

	// Evaluate the size and apply random factors
	sz := EvalNumber(reqdat[`s`], seed^lastmod)
	reqdat[`sz`] = strconv.FormatInt(sz, 10)

	if etagstr, ok := reqdat[`etag`]; ok {
		if len(etagstr) == 0 {
			etagnum := lastmod ^ sz
			etagstr = fmt.Sprintf("\"%d\"", etagnum)
		}
		w.Header().Set(`Etag`, etagstr)
	}

	// Special header to hard set cache control
	if ccstr, ok := reqdat[XDtpCcHdrStr]; ok {
		if 0 < len(ccstr) {
			w.Header().Set(`Cache-Control`, ccstr)
		}
	}

	// Initial delay, hard set
	delaydur, err := time.ParseDuration(reqdat["idelay"])
	if err == nil {
		if 0 < delaydur {
			time.Sleep(delaydur)
		}
	} else { // initial delay with rand
		delay := EvalNumber(reqdat[`dly`], seed^now)
		if 0 < delay {
			time.Sleep(time.Duration(delay))
		}
	}

	// These should override any of the previous
	if hdr, ok := reqdat[`hdr`]; ok {
		parts := strings.Split(hdr, `.`)
		DebugLogf("Reading hdr: %v", parts)
		for ind := 1; ind < len(parts); ind += 2 {
			w.Header().Set(parts[ind-1], parts[ind])
		}
	}

	if hdr, ok := reqdat[`hdr64`]; ok {
		parts := strings.Split(hdr, `.`)
		DebugLogf("Reading hdr64: %v", parts)
		for ind := 1; ind < len(parts); ind += 2 {
			val, err := base64.URLEncoding.DecodeString(parts[ind])
			if err != nil {
				DebugLogf("Failed to decode hdr64 '%s': %s.\n", parts[ind], err.Error())
				continue
			}
			w.Header().Set(parts[ind-1], string(val))
		}
	}

	sc := EvalNumber(reqdat[`sc`], seed^now)
	if sc != 200 && sc != 0 {
		w.WriteHeader(int(sc))
		return false
	}

	// remove request headers
	if rmhdrs, ok := reqdat[`rmhdrs`]; ok {
		hdrs := strings.Split(rmhdrs, `.`)
		for _, hdr := range hdrs {
			r.Header.Del(hdr)
		}
	}

	return true
}
