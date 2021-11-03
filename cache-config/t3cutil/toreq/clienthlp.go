package toreq

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
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

// GetProfileByName returns the profile with the given name from Traffic Ops.
// It is a helper function equivalent to calling GetProfiles with RequestOptions with the Values (query string) with the key name set to the Profile name.
// If opts.Values[name] exists, it is overwritten with name.
func GetProfileByName(toClient *toclient.Session, name string, opts *toclient.RequestOptions) (tc.Profile, toclientlib.ReqInf, error) {
	if opts == nil {
		opts = &toclient.RequestOptions{}
	}
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", name)
	profiles, reqInf, err := toClient.GetProfiles(*opts)
	if err != nil {
		return tc.Profile{}, reqInf, err
	} else if reqInf.StatusCode == http.StatusNotModified {
		return tc.Profile{}, reqInf, nil
	} else if len(profiles.Response) == 0 {
		return tc.Profile{}, reqInf, fmt.Errorf("name '"+string(name)+" ' not found (no error, but len 0) reqInf %+v profiles %+v", reqInf, profiles)
	} else if len(profiles.Response) > 1 {
		return tc.Profile{}, reqInf, fmt.Errorf("expected 1 profile, got len %v val %+v", len(profiles.Response), profiles.Response)
	}
	return profiles.Response[0], reqInf, nil
}

// GetParametersByConfigFile returns the parameters with the given config file from Traffic Ops.
// It is a helper function equivalent to calling GetParameters with RequestOptions with the Values (query string) with the key configFile set to the config file.
// If opts.Values[configFile] exists, it is overwritten with name.
func GetParametersByConfigFile(toClient *toclient.Session, configFile string, opts *toclient.RequestOptions) ([]tc.Parameter, toclientlib.ReqInf, error) {
	if opts == nil {
		opts = &toclient.RequestOptions{}
	}
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("configFile", configFile)
	params, reqInf, err := toClient.GetParameters(*opts)
	return params.Response, reqInf, err
}

// GetParametersByName returns the parameters with the given name from Traffic Ops.
// It is a helper function equivalent to calling GetParameters with RequestOptions with the Values (query string) with the key name set to the name.
// If opts.Values[name] exists, it is overwritten with name.
func GetParametersByName(toClient *toclient.Session, name string, opts *toclient.RequestOptions) ([]tc.Parameter, toclientlib.ReqInf, error) {
	if opts == nil {
		opts = &toclient.RequestOptions{}
	}
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", name)
	params, reqInf, err := toClient.GetParameters(*opts)
	return params.Response, reqInf, err
}

func GetCDNByName(toClient *toclient.Session, name tc.CDNName, opts *toclient.RequestOptions) (tc.CDN, toclientlib.ReqInf, error) {
	if opts == nil {
		opts = &toclient.RequestOptions{}
	}
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("name", string(name))
	cdns, reqInf, err := toClient.GetCDNs(*opts)

	if err != nil {
		return tc.CDN{}, reqInf, err
	} else if reqInf.StatusCode == http.StatusNotModified {
		return tc.CDN{}, reqInf, nil
	} else if len(cdns.Response) == 0 {
		return tc.CDN{}, reqInf, fmt.Errorf("name '"+string(name)+" ' not found (no error, but len 0) reqInf %+v cdns %+v", reqInf, cdns)
	} else if len(cdns.Response) > 1 {
		return tc.CDN{}, reqInf, fmt.Errorf("expected 1, got len %v val %+v", len(cdns.Response), cdns.Response)
	}
	return cdns.Response[0], reqInf, nil
}

// GetDeliveryServiceURLSigKeys gets the URL Sig keys from Traffic Ops for the given delivery service.
// It is a helper function that calls traffic_ops/v4-client.Session.GetDeliveryServiceURLSignatureKeys
// to avoid confusion around the protocol named URL Sig.
func GetDeliveryServiceURLSigKeys(toClient *toclient.Session, dsName string, opts *toclient.RequestOptions) (tc.URLSigKeys, toclientlib.ReqInf, error) {
	if opts == nil {
		opts = &toclient.RequestOptions{}
	}
	resp, reqInf, err := toClient.GetDeliveryServiceURLSignatureKeys(dsName, *opts)
	return resp.Response, reqInf, err
}

// ReqOpts takes an http.Header and returns a traffic_ops/v4-client.RequestOptions with that header.
// This is a helper function, for brevity.
func ReqOpts(hdr http.Header) *toclient.RequestOptions {
	opts := toclient.NewRequestOptions()
	opts.Header = hdr
	return &opts
}
