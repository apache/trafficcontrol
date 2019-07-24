package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

var scopeConfigFileFuncs = map[string]func(toClient **toclient.Session, cfg Cfg, resource string, fileName string) (string, error){
	"cdns":     GetConfigFileCDN,
	"servers":  GetConfigFileServer,
	"profiles": GetConfigFileProfile,
}

func GetConfigFile(toClient **toclient.Session, cfg Cfg) (string, error) {

	pathParts := strings.Split(cfg.TOURL.Path, "/")

	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFile pathParts %++v\n", pathParts)

	if len(pathParts) < 8 {
		fmt.Fprintf(os.Stderr, "DEBUG GetConfigFile pathParts < 7, calling TO\n")
		return GetConfigFileFromTrafficOps(toClient, cfg)
	}
	scope := pathParts[3]
	resource := pathParts[4]
	fileName := pathParts[7]

	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFile scope '%v' resource '%v' fileName '%v'\n", scope, resource, fileName)

	if scopeConfigFileFunc, ok := scopeConfigFileFuncs[scope]; ok {
		return scopeConfigFileFunc(toClient, cfg, resource, fileName)
	}

	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFile unknown scope, calling TO\n")
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

func GetConfigFileCDN(toClient **toclient.Session, cfg Cfg, cdnNameOrID string, fileName string) (string, error) {
	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFileCDN cdn '"+cdnNameOrID+"' fileName '"+fileName+"'\n")
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

func GetConfigFileProfile(toClient **toclient.Session, cfg Cfg, profileNameOrID string, fileName string) (string, error) {
	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFileProfile profile '"+profileNameOrID+"' fileName '"+fileName+"'\n")
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

var serverConfigFileFuncs = map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error){
	"parent.config": GetConfigFileServerParentDotConfig,
}

func GetConfigFileServer(toClient **toclient.Session, cfg Cfg, serverNameOrID string, fileName string) (string, error) {

	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFileServer server '"+serverNameOrID+"' fileName '"+fileName+"'\n")
	if getCfgFunc, ok := serverConfigFileFuncs[fileName]; ok {
		return getCfgFunc(toClient, cfg, serverNameOrID)
	}
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

func GetConfigFileFromTrafficOps(toClient **toclient.Session, cfg Cfg) (string, error) {
	path := cfg.TOURL.Path
	if cfg.TOURL.RawQuery != "" {
		path += "?" + cfg.TOURL.RawQuery
	}
	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFile path '"+path+"'\n")
	fmt.Fprintf(os.Stderr, "DEBUG GetConfigFile url '"+cfg.TOURL.String()+"'\n")

	body, err := TrafficOpsRequest(toClient, cfg, http.MethodGet, cfg.TOURL.String(), nil)
	if err != nil {
		return "", errors.New("Requesting path '" + path + "': " + err.Error())
	}

	WriteCookiesToFile(CookiesToString((*toClient).Client.Jar.Cookies(cfg.TOURL)), cfg.TempDir)

	return string(body), nil
}
