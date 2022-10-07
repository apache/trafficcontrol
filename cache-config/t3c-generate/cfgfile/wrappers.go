package cfgfile

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
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// This file has wrappers that turn lib/go-atscfg Make funcs into ConfigFileFunc types.
//
// We don't want to make lib/go-atscfg functions take a TOData, because then users wanting to generate a single file would have to fetch all kinds of data that file doesn't need, or else pass objects they know it doesn't currently need as nil and risk it crashing if that func is changed to use it in the future.
//
// But it's useful to map filenames to functions for dispatch. Hence these wrappers.
//
// The atsMajorVersion may be 0 to default to the Server Package Parameter.
//
// MakeConfigFilesList returns the list of config files, any warnings, and any error.
func MakeConfigFilesList(toData *t3cutil.ConfigData, dir string, atsMajorVersion uint) ([]atscfg.CfgMeta, []string, error) {
	configFiles, warnings, err := atscfg.MakeConfigFilesList(
		dir,
		toData.Server,
		toData.ServerParams,
		toData.DeliveryServices,
		toData.DeliveryServiceServers,
		toData.GlobalParams,
		toData.CacheGroups,
		toData.Topologies,
		&atscfg.ConfigFilesListOpts{
			ATSMajorVersion: atsMajorVersion,
		},
	)
	return configFiles, warnings, err
}

func Make12MFacts(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.Config12MFactsOpts{HdrComment: hdrCommentTxt}
	return atscfg.Make12MFacts(toData.Server, opts)
}

func MakeATSDotRules(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.ATSDotRulesOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeATSDotRules(toData.Server, toData.ServerParams, opts)
}

func MakeAstatsDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.AStatsDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeAStatsDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakeBGFetchDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.BGFetchDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeBGFetchDotConfig(toData.Server, opts)
}

func MakeCacheDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.CacheDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeCacheDotConfig(toData.Server, toData.Servers, toData.DeliveryServices, toData.DeliveryServiceServers, opts)
}

func MakeChkconfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeChkconfig(toData.ServerParams, nil)
}

func MakeDropQStringDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.DropQStringDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeDropQStringDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakeHostingDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.HostingDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeHostingDotConfig(toData.Server, toData.Servers, toData.ServerParams, toData.DeliveryServices, toData.DeliveryServiceServers, toData.Topologies, opts)
}

func MakeIPAllowDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.IPAllowDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeIPAllowDotConfig(
		toData.ServerParams,
		toData.Server,
		toData.Servers,
		toData.CacheGroups,
		toData.Topologies,
		opts,
	)
}

func MakeIPAllowDotYAML(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.IPAllowDotYAMLOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeIPAllowDotYAML(
		toData.ServerParams,
		toData.Server,
		toData.Servers,
		toData.CacheGroups,
		toData.Topologies,
		opts,
	)
}

func MakeLoggingDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.LoggingDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeLoggingDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakeLoggingDotYAML(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeLoggingDotYAML(
		toData.Server,
		toData.ServerParams,
		&atscfg.LoggingDotYAMLOpts{
			HdrComment:      hdrCommentTxt,
			ATSMajorVersion: cfg.ATSMajorVersion,
		},
	)
}

func MakeSSLServerNameYAML(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeSSLServerNameYAML(
		toData.Server,
		toData.DeliveryServices,
		toData.DeliveryServiceServers,
		toData.DeliveryServiceRegexes,
		toData.ParentConfigParams,
		toData.CDN,
		toData.Topologies,
		toData.CacheGroups,
		toData.ServerCapabilities,
		toData.DSRequiredCapabilities,
		&atscfg.SSLServerNameYAMLOpts{
			HdrComment:         hdrCommentTxt,
			VerboseComments:    true, // TODO add a CLI flag
			DefaultTLSVersions: cfg.DefaultTLSVersions,
			DefaultEnableH2:    cfg.DefaultEnableH2,
		},
	)
}

func MakeSNIDotYAML(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeSNIDotYAML(
		toData.Server,
		toData.DeliveryServices,
		toData.DeliveryServiceServers,
		toData.DeliveryServiceRegexes,
		toData.ParentConfigParams,
		toData.CDN,
		toData.Topologies,
		toData.CacheGroups,
		toData.ServerCapabilities,
		toData.DSRequiredCapabilities,
		&atscfg.SNIDotYAMLOpts{
			HdrComment:         hdrCommentTxt,
			VerboseComments:    true, // TODO add a CLI flag
			DefaultTLSVersions: cfg.DefaultTLSVersions,
			DefaultEnableH2:    cfg.DefaultEnableH2,
		},
	)
}

func MakeLogsXMLDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.LogsXMLDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeLogsXMLDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakePackages(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakePackages(toData.ServerParams, nil)
}

func MakeParentDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeParentDotConfig(
		toData.DeliveryServices,
		toData.Server,
		toData.Servers,
		toData.Topologies,
		toData.ServerParams,
		toData.ParentConfigParams,
		toData.ServerCapabilities,
		toData.DSRequiredCapabilities,
		toData.CacheGroups,
		toData.DeliveryServiceServers,
		toData.CDN,
		&atscfg.ParentConfigOpts{
			HdrComment:      hdrCommentTxt,
			AddComments:     cfg.ParentComments, // TODO add a CLI flag?
			ATSMajorVersion: cfg.ATSMajorVersion,
		},
	)
}

func MakePluginDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.PluginDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakePluginDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakeRecordsDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeRecordsDotConfig(
		toData.Server,
		toData.ServerParams,
		&atscfg.RecordsConfigOpts{
			ReleaseViaStr:           cfg.ViaRelease,
			DNSLocalBindServiceAddr: cfg.SetDNSLocalBind,
			HdrComment:              hdrCommentTxt,
			NoOutgoingIP:            cfg.NoOutgoingIP,
		},
	)
}

func MakeRegexRevalidateDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.RegexRevalidateDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeRegexRevalidateDotConfig(toData.Server, toData.DeliveryServices, toData.GlobalParams, toData.Jobs, opts)
}

func MakeRemapDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	remapAndCacheKeyParams := []tc.Parameter{}
	remapAndCacheKeyParams = append(remapAndCacheKeyParams, toData.RemapConfigParams...)
	remapAndCacheKeyParams = append(remapAndCacheKeyParams, toData.CacheKeyConfigParams...)
	return atscfg.MakeRemapDotConfig(
		toData.Server,
		toData.DeliveryServices,
		toData.DeliveryServiceServers,
		toData.DeliveryServiceRegexes,
		toData.ServerParams,
		toData.CDN,
		remapAndCacheKeyParams,
		toData.Topologies,
		toData.CacheGroups,
		toData.ServerCapabilities,
		toData.DSRequiredCapabilities,
		cfg.Dir,
		&atscfg.RemapDotConfigOpts{
			HdrComment:        hdrCommentTxt,
			VerboseComments:   true,
			UseStrategies:     cfg.UseStrategies == t3cutil.UseStrategiesFlagTrue || cfg.UseStrategies == t3cutil.UseStrategiesFlagCore,
			UseStrategiesCore: cfg.UseStrategies == t3cutil.UseStrategiesFlagCore,
			ATSMajorVersion:   cfg.ATSMajorVersion,
		},
	)
}

func MakeSSLMultiCertDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.SSLMultiCertDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeSSLMultiCertDotConfig(toData.Server, toData.DeliveryServices, opts)
}

func MakeStorageDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.StorageDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeStorageDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakeSysCtlDotConf(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.SysCtlDotConfOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeSysCtlDotConf(toData.Server, toData.ServerParams, opts)
}

func MakeVolumeDotConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.VolumeDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeVolumeDotConfig(toData.Server, toData.ServerParams, opts)
}

func MakeHeaderRewrite(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeHeaderRewriteDotConfig(
		fileName,
		toData.DeliveryServices,
		toData.DeliveryServiceServers,
		toData.Server,
		toData.Servers,
		toData.CacheGroups,
		toData.ServerParams,
		toData.ServerCapabilities,
		toData.DSRequiredCapabilities,
		toData.Topologies,
		&atscfg.HeaderRewriteDotConfigOpts{
			HdrComment:      hdrCommentTxt,
			ATSMajorVersion: cfg.ATSMajorVersion,
		},
	)
}

func MakeRegexRemap(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.RegexRemapDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeRegexRemapDotConfig(fileName, toData.Server, toData.DeliveryServices, opts)
}

func MakeSetDSCP(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.SetDSCPDotConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeSetDSCPDotConfig(fileName, toData.Server, opts)
}

func MakeURLSigConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.URLSigConfigOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeURLSigConfig(fileName, toData.Server, toData.ServerParams, toData.URLSigKeys, opts)
}

func MakeURISigningConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeURISigningConfig(fileName, toData.URISigningKeys, nil)
}

func MakeStrategiesDotYAML(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	return atscfg.MakeStrategiesDotYAML(
		toData.DeliveryServices,
		toData.Server,
		toData.Servers,
		toData.Topologies,
		toData.ServerParams,
		toData.ParentConfigParams,
		toData.ServerCapabilities,
		toData.DSRequiredCapabilities,
		toData.CacheGroups,
		toData.DeliveryServiceServers,
		toData.CDN,
		&atscfg.StrategiesYAMLOpts{
			HdrComment:      hdrCommentTxt,
			VerboseComments: cfg.ParentComments, // TODO add a CLI flag?
			ATSMajorVersion: cfg.ATSMajorVersion,
		},
	)
}

func MakeUnknownConfig(toData *t3cutil.ConfigData, fileName string, hdrCommentTxt string, cfg config.Cfg) (atscfg.Cfg, error) {
	opts := &atscfg.ServerUnknownOpts{HdrComment: hdrCommentTxt}
	return atscfg.MakeServerUnknown(fileName, toData.Server, toData.ServerParams, opts)
}
