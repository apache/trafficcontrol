package grvssl

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"strings"
	"github.com/apache/trafficcontrol/lib/go-log"
	"io/ioutil"
	"net"
	"net/http"
	"fmt"
	"time"

	"github.com/spacemonkeygo/openssl"
	"golang.org/x/net/http2"

	"github.com/apache/trafficcontrol/grove/remapdata"
	"github.com/apache/trafficcontrol/grove/web"
)

type TlsCtxStore struct {
	CtxMap map[string]*openssl.Ctx
}

var CtxStore TlsCtxStore

func CreateTlsCtxStore(rules []remapdata.RemapRule) error {

	CtxStore = TlsCtxStore{
		CtxMap: make(map[string]*openssl.Ctx),
	}

	for _, rule := range rules {
		if rule.CertificateFile == "" && rule.CertificateKeyFile == "" {
			continue
		}
		if rule.CertificateFile == "" {
			log.Errorln("rule " + rule.Name + ": has a certificate but no key, using default certificate\n")
			continue
		}
		if rule.CertificateKeyFile == "" {
			log.Errorln("rule " + rule.Name + ": has a key but no certificate, using default certificate\n")
			continue
		}
		ctx, cname, err := NewCtxFromFiles(rule.CertificateFile, rule.CertificateKeyFile)
		if err != nil {
			log.Errorln("rule " + rule.Name + ": unable to get a context for this certificate.")
			continue
		}
		ctx.SetTLSExtServernameCallback(SniCallback)
		CtxStore.CtxMap[cname] = ctx
		fmt.Println("loaded context for rule " + rule.Name + " with cname: " + cname)
	}
	return nil
}

func InterceptListenTLS(network string, laddr string, defaultContext *openssl.Ctx) (net.Listener, *web.ConnMap, func(net.Conn, http.ConnState), error) {
	l, err := net.Listen(network, laddr)

	if err != nil {
		return l, nil, nil, err
	}


	connMap := web.NewConnMap()
	interceptListener := openssl.NewListener(l, defaultContext)
	return interceptListener, connMap, web.GetConnStateCallback(connMap), nil
}

func NewCtxFromFiles(cert_file string, key_file string) (*openssl.Ctx, string, error) {
	ctx, err := openssl.NewCtx()
	if err != nil {
		return nil, "", err
	}

	cert_bytes, err := ioutil.ReadFile(cert_file)
	if err != nil {
		return nil, "", err
	}

	certs := openssl.SplitPEM(cert_bytes)
	if len(certs) == 0 {
		return nil, "", fmt.Errorf("No PEM certificate found in '%s'", cert_file)
	}
	first, certs := certs[0], certs[1:]
	cert, err := openssl.LoadCertificateFromPEM(first)
	if err != nil {
		return nil, "", err
	}

	name, err := cert.GetSubjectName()
	if err != nil {
		return nil, "", fmt.Errorf("Unable to parse the subject name from this certificate.")
	}

	cname, ok := name.GetEntry(openssl.NID_commonName)
	if !ok {
		return nil, "", fmt.Errorf("Unable to parse the common name from this certificate.")
	}

	err = ctx.UseCertificate(cert)
	if err != nil {
		return nil, "", err
	}

	for _, pem := range certs {
		cert, err := openssl.LoadCertificateFromPEM(pem)
		if err != nil {
			return nil, "", err
		}
		err = ctx.AddChainCertificate(cert)
		if err != nil {
			return nil, "", err
		}
	}

	key_bytes, err := ioutil.ReadFile(key_file)
	if err != nil {
		return nil, "", err
	}

	key, err := openssl.LoadPrivateKeyFromPEM(key_bytes)
	if err != nil {
		return nil, "", err
	}

	err = ctx.UsePrivateKey(key)
	if err != nil {
		return nil, "", err
	}

	return ctx, cname, nil
}

func SniCallback(ssl *openssl.SSL) openssl.SSLTLSExtErr {
	var ctx *openssl.Ctx
	var i int = 0

	hostname := ssl.GetServername()
	cnames := make([]string, 3, 3)
	dparts := strings.Split(hostname, ".")
	domain := strings.Join(dparts[1:], ".")

	cnames[0] = hostname
	cnames[1] = domain
	cnames[2] = ("*." + domain)

	for _, name := range cnames {
		ctx = CtxStore.CtxMap[name]
		if ctx != nil {
			break
		}
		i++
	}

	if ctx != nil {
		ssl.SetSSLCtx(ctx)
		log.Errorln("found a context for: " + cnames[i] + ", setting the ssl context")
	} else {
		for _, name := range cnames {
			log.Errorln("could not find a context for: " + name)
		}
	}

	return 0
}

func StartServer(handler http.Handler, listener net.Listener, connState func(net.Conn, http.ConnState), port int, idleTimeout time.Duration, readTimeout time.Duration, writeTimeout time.Duration, protocol string) *http.Server {
	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%d", port),
		ConnState:    connState,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// TODO configurable H2 timeouts and buffer sizes
	h2Conf := &http2.Server{
		IdleTimeout: idleTimeout,
	}
	if err := http2.ConfigureServer(server, h2Conf); err != nil {
		log.Errorln(" server configuring HTTP/2: " + err.Error())
	}

	go func() {
		log.Infof("listening on %s://%d\n", protocol, port)
		if err := server.Serve(listener); err != nil {
			log.Errorf("serving %s port %v: %v\n", strings.ToUpper(protocol), port, err)
		}
	}()
	return server
}