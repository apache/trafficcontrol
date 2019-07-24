package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/ogier/pflag"
)

type Cfg struct {
	TOURL   *url.URL
	TOUser  string
	TOPass  string
	TempDir string
}

func GetCfg() (Cfg, error) {
	toURLPtr := flag.StringP("traffic-ops-url", "u", "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with the environment variable TO_URL.")
	toUserPtr := flag.StringP("traffic-ops-user", "U", "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER.")
	toPassPtr := flag.StringP("traffic-ops-password", "P", "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS.")
	noCachePtr := flag.BoolP("no-cache", "n", false, "Whether not to use existing cache files. Optional. Cache files will still be created, existing ones just won't be used.")
	flag.Parse()

	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	noCache := *noCachePtr

	urlSourceStr := "argument" // for error messages
	if toURL == "" {
		urlSourceStr = "environment variable"
		toURL = os.Getenv("TO_URL")
	}
	if toUser == "" {
		toUser = os.Getenv("TO_USER")
	}
	if toPass == "" {
		toPass = os.Getenv("TO_PASS")
	}

	if strings.TrimSpace(toURL) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-url or TO_URL environment variable. Usage: ./" + AppName + " --traffic-ops-url myurl --traffic-ops-user myuser --traffic-ops-password mypass")
	}
	if strings.TrimSpace(toUser) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-user or TO_USER environment variable. Usage: ./" + AppName + " --traffic-ops-url myurl --traffic-ops-user myuser --traffic-ops-password mypass")
	}
	if strings.TrimSpace(toPass) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-password or TO_PASS environment variable. Usage: ./" + AppName + " --traffic-ops-url myurl --traffic-ops-user myuser --traffic-ops-password mypass")
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err := ValidateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	tmpDir := os.TempDir()
	tmpDir = filepath.Join(tmpDir, TempSubdir)

	if noCache {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Fprintf(os.Stderr, "error deleting cache directory '"+tmpDir+"': "+err.Error()+"\n")
		}
	}

	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return Cfg{}, errors.New("creating temp directory '" + tmpDir + "': " + err.Error())
	}
	if err := ValidateDirWriteable(tmpDir); err != nil {
		return Cfg{}, errors.New("validating temp directory is writeable '" + tmpDir + "': " + err.Error())
	}

	return Cfg{
		TOURL:   toURLParsed,
		TOUser:  toUser,
		TOPass:  toPass,
		TempDir: tmpDir,
	}, nil
}

func ValidateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme expected 'http' or 'https', actual '" + u.Scheme + "'")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("no host")
	}
	return nil
}

func ValidateDirWriteable(dir string) error {
	testFileName := "testwrite.txt"
	testFilePath := filepath.Join(dir, testFileName)
	if err := os.RemoveAll(testFilePath); err != nil {
		// TODO don't log? This can be normal
		fmt.Fprintf(os.Stderr, "removing temp test file '"+testFilePath+"': "+err.Error()+"\n")
	}

	fl, err := os.Create(testFilePath)
	if err != nil {
		return errors.New("creating temp test file '" + testFilePath + "': " + err.Error())
	}
	defer fl.Close()

	if _, err := fl.WriteString("test"); err != nil {
		return errors.New("writing to temp test file '" + testFilePath + "': " + err.Error())
	}

	return nil
}
