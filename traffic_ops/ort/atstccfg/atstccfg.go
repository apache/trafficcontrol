package main

import (
	"fmt"
	"os"
	"time"
)

const AppName = "atstccfg"
const Version = "0.1"
const UserAgent = AppName + "/" + Version

const APIVersion = "1.2"
const TempSubdir = AppName + "_cache"
const TempCookieFileName = "cookies"
const TOCookieName = "mojolicious"

// TODO make the below configurable?
const TOInsecure = false
const TOTimeout = time.Second * 10
const CacheFileMaxAge = time.Minute

func main() {
	cfg, err := GetCfg()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config: "+err.Error()+"\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "DEBUG URL: '"+cfg.TOURL.String()+"' User: '"+cfg.TOUser+"' Pass: '"+cfg.TOPass+"'\n")
	fmt.Fprintf(os.Stderr, "DEBUG TempDir: '"+cfg.TempDir+"'\n")

	toFQDN := cfg.TOURL.Scheme + "://" + cfg.TOURL.Host
	fmt.Fprintf(os.Stderr, "DEBUG TO FQDN: '"+toFQDN+"'\n")
	fmt.Fprintf(os.Stderr, "DEBUG TO URL: '"+cfg.TOURL.String()+"'\n")

	toClient, err := GetClient(toFQDN, cfg.TOUser, cfg.TOPass, cfg.TempDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Logging in to Traffic Ops: "+err.Error()+"\n")
		os.Exit(1)
	}

	cfgFile, err := GetConfigFile(&toClient, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config file '"+cfg.TOURL.String()+"' from Traffic Ops: "+err.Error()+"\n")
		os.Exit(1)
	}
	fmt.Println(cfgFile)
	os.Exit(0)
}
