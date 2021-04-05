package riaksvc

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
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/basho/riak-go-client"
)

const DeliveryServiceSSLKeysBucket = "ssl"
const DNSSECKeysBucket = "dnssec"
const DSSSLKeyVersionLatest = "latest"
const DefaultDSSSLKeyVersion = DSSSLKeyVersionLatest
const URLSigKeysBucket = "url_sig_keys"

// CDNURIKeysBucket is the namespace or bucket used for CDN URI signing keys.
const CDNURIKeysBucket = "cdn_uri_sig_keys"

func MakeDSSSLKeyKey(dsName, version string) string {
	if version == "" {
		version = DefaultDSSSLKeyVersion
	}
	return dsName + "-" + version
}

func GetDeliveryServiceSSLKeysObjV15(xmlID string, version string, tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	key := tc.DeliveryServiceSSLKeysV15{}
	found := false
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		// get the deliveryservice ssl keys by xmlID and version
		ro, err := FetchObjectValues(MakeDSSSLKeyKey(xmlID, version), DeliveryServiceSSLKeysBucket, cluster)
		if err != nil {
			return err
		}
		if len(ro) == 0 {
			return nil // not found
		}
		if err := json.Unmarshal(ro[0].Value, &key); err != nil {
			log.Errorf("failed at unmarshaling sslkey response: %s\n", err)
			return errors.New("unmarshalling Riak result: " + err.Error())
		}
		found = true
		return nil
	})
	if err != nil {
		return key, false, err
	}
	return key, found, nil
}

func PutDeliveryServiceSSLKeysObj(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint) error {
	keyJSON, err := json.Marshal(&key)
	if err != nil {
		return errors.New("marshalling key: " + err.Error())
	}
	err = WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		obj := &riak.Object{
			ContentType:     "application/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             MakeDSSSLKeyKey(key.DeliveryService, key.Version.String()),
			Value:           []byte(keyJSON),
		}
		if err := SaveObject(obj, DeliveryServiceSSLKeysBucket, cluster); err != nil {
			return errors.New("saving Riak object: " + err.Error())
		}
		obj.Key = MakeDSSSLKeyKey(key.DeliveryService, DSSSLKeyVersionLatest)
		if err := SaveObject(obj, DeliveryServiceSSLKeysBucket, cluster); err != nil {
			return errors.New("saving Riak object: " + err.Error())
		}
		return nil
	})
	return err
}

func Ping(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint) (tc.RiakPingResp, error) {
	servers, err := GetRiakServers(tx, riakPort)
	if err != nil {
		return tc.RiakPingResp{}, errors.New("getting riak servers: " + err.Error())
	}
	for _, server := range servers {
		cluster, err := GetRiakStorageCluster([]ServerAddr{server}, authOpts)
		if err != nil {
			log.Errorf("RiakServersToCluster error for server %+v: %+v\n", server, err.Error())
			continue // try another server
		}
		if err = cluster.Start(); err != nil {
			log.Errorln("starting Riak cluster (for ping): " + err.Error())
			continue
		}
		if err := PingCluster(cluster); err != nil {
			if err := cluster.Stop(); err != nil {
				log.Errorln("stopping Riak cluster (after ping error): " + err.Error())
			}
			log.Errorf("Riak PingCluster error for server %+v: %+v\n", server, err.Error())
			continue
		}
		if err := cluster.Stop(); err != nil {
			log.Errorln("stopping Riak cluster (after ping success): " + err.Error())
		}
		return tc.RiakPingResp{Status: "OK", Server: server.FQDN + ":" + server.Port}, nil
	}
	return tc.RiakPingResp{}, errors.New("failed to ping any Riak server")
}

func GetDNSSECKeys(cdnName string, tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint) (tc.DNSSECKeysRiak, bool, error) {
	key := tc.DNSSECKeysRiak{}
	found := false
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		ro, err := FetchObjectValues(cdnName, DNSSECKeysBucket, cluster)
		if err != nil {
			return err
		}
		if len(ro) == 0 {
			return nil // not found
		}
		if err := json.Unmarshal(ro[0].Value, &key); err != nil {
			return errors.New("unmarshalling Riak dnssec response: " + err.Error())
		}
		found = true
		return nil
	})
	if err != nil {
		return key, false, err
	}
	return key, found, nil
}

func PutDNSSECKeys(keys tc.DNSSECKeysRiak, cdnName string, tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint) error {
	keyJSON, err := json.Marshal(&keys)
	if err != nil {
		return errors.New("marshalling keys: " + err.Error())
	}

	err = WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		obj := &riak.Object{
			ContentType:     "application/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             cdnName,
			Value:           []byte(keyJSON),
		}
		if err = SaveObject(obj, DNSSECKeysBucket, cluster); err != nil {
			return errors.New("saving Riak object: " + err.Error())
		}
		return nil
	})
	return err
}

func DeleteDNSSECKeys(cdnName string, tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint) error {
	cluster, err := GetPooledCluster(tx, authOpts, riakPort)
	if err != nil {
		return errors.New("getting riak cluster: " + err.Error())
	}
	if err := DeleteObject(cdnName, DNSSECKeysBucket, cluster); err != nil {
		return errors.New("deleting riak object: " + err.Error())
	}
	return nil
}

func GetBucketKey(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, bucket string, key string) ([]byte, bool, error) {
	val := []byte{}
	found := false
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		// get the deliveryservice ssl keys by xmlID and version
		ro, err := FetchObjectValues(key, bucket, cluster)
		if err != nil {
			return err
		}
		if len(ro) == 0 {
			return nil // not found
		}
		val = ro[0].Value
		found = true
		return nil
	})
	if err != nil {
		return val, false, err
	}
	return val, found, nil
}

func DeleteDSSSLKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, xmlID string, version string) error {
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		if err := DeleteObject(MakeDSSSLKeyKey(xmlID, version), DeliveryServiceSSLKeysBucket, cluster); err != nil {
			return errors.New("deleting SSL keys: " + err.Error())
		}
		return nil
	})
	return err
}

// DeleteDeliveryServicesSSLKey deletes a Delivery Service SSL key.
// This should almost never be used directly, prefer DeleteDSSSLKeys instead.
// This should only be used to delete keys, which may not conform to the MakeDSSSLKeyKey format. For example when deleting all keys on a delivery service, and some may have been created manually outside Traffic Ops, or are otherwise malformed.
func DeleteDeliveryServicesSSLKey(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, key string) error {
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		if err := DeleteObject(key, DeliveryServiceSSLKeysBucket, cluster); err != nil {
			return errors.New("deleting SSL keys: " + err.Error())
		}
		return nil
	})
	return err
}

func GetURISigningKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, xmlID string) ([]byte, bool, error) {
	cluster, err := GetPooledCluster(tx, authOpts, riakPort)
	if err != nil {
		return nil, false, errors.New("getting pooled Riak cluster: " + err.Error())
	}
	ro, err := FetchObjectValues(xmlID, CDNURIKeysBucket, cluster)
	if err != nil {
		return nil, false, errors.New("fetching riak objects: " + err.Error())
	}
	if len(ro) == 0 {
		bts, err := json.Marshal(tc.URISignerKeyset{})
		if err != nil {
			return nil, false, errors.New("marshalling empty URISignerKeyset: " + err.Error())
		}
		return bts, false, nil
	}
	if ro[0].Value == nil {
		return ro[0].Value, false, nil
	}
	return ro[0].Value, true, nil
}

func DeleteURISigningKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, xmlID string) error {
	cluster, err := GetPooledCluster(tx, authOpts, riakPort)
	if err != nil {
		return errors.New("getting pooled Riak cluster: " + err.Error())
	}
	if err := DeleteObject(xmlID, CDNURIKeysBucket, cluster); err != nil {
		return errors.New("deleting object: " + err.Error())
	}
	return nil
}

func PutURISigningKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, xmlID string, keysJson []byte) error {
	cluster, err := GetPooledCluster(tx, authOpts, riakPort)
	if err != nil {
		return errors.New("getting pooled Riak cluster: " + err.Error())
	}
	obj := &riak.Object{
		ContentType:     "text/json",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Key:             xmlID,
		Value:           keysJson,
	}
	if err = SaveObject(obj, CDNURIKeysBucket, cluster); err != nil {
		return errors.New("saving riak object: " + err.Error())
	}
	return nil
}

// GetURLSigConfigFileName returns the filename of the Apache Traffic Server URLSig config file
// TODO move to ats config directory/file
func GetURLSigConfigFileName(ds tc.DeliveryServiceName) string {
	return "url_sig_" + string(ds) + ".config"
}

func GetURLSigKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, ds tc.DeliveryServiceName) (tc.URLSigKeys, bool, error) {
	val := tc.URLSigKeys{}
	found := false
	key := GetURLSigConfigFileName(ds)
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		ro, err := FetchObjectValues(key, URLSigKeysBucket, cluster)
		if err != nil {
			return err
		}
		if len(ro) == 0 {
			return nil // not found
		}
		if err := json.Unmarshal(ro[0].Value, &val); err != nil {
			return errors.New("unmarshalling Riak response: " + err.Error())
		}
		found = true
		return nil
	})
	if err != nil {
		return val, false, err
	}
	return val, found, nil
}

func PutURLSigKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, ds tc.DeliveryServiceName, keys tc.URLSigKeys) error {
	keyJSON, err := json.Marshal(&keys)
	if err != nil {
		return errors.New("marshalling keys: " + err.Error())
	}
	err = WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		obj := &riak.Object{
			ContentType:     "application/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             GetURLSigConfigFileName(ds),
			Value:           []byte(keyJSON),
		}
		if err = SaveObject(obj, URLSigKeysBucket, cluster); err != nil {
			return errors.New("saving Riak object: " + err.Error())
		}
		return nil
	})
	return err
}

const SSLKeysIndex = "sslkeys"
const CDNSSLKeysLimit = 1000 // TODO: emulates Perl; reevaluate?

func GetCDNSSLKeysObj(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, cdnName string) ([]tc.CDNSSLKey, error) {
	keys := []tc.CDNSSLKey{}
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		// get the deliveryservice ssl keys by xmlID and version
		query := `cdn:` + cdnName
		filterQuery := `_yz_rk:*latest`
		fields := []string{"deliveryservice", "hostname", "certificate.crt", "certificate.key"}
		searchDocs, err := Search(cluster, SSLKeysIndex, query, filterQuery, CDNSSLKeysLimit, fields)
		if err != nil {
			return errors.New("riak search error: " + err.Error())
		}
		if len(searchDocs) == 0 {
			return nil // no error, and leave keys empty
		}
		keys = SearchDocsToCDNSSLKeys(searchDocs)
		return nil
	})
	if err != nil {
		return nil, errors.New("with cluster error: " + err.Error())
	}
	return keys, nil
}

// SearchDocsToCDNSSLKeys converts the SearchDoc array returned by Riak into a CDNSSLKey slice. If a SearchDoc doesn't contain expected fields, it creates the key with those fields defaulted to empty strings.
func SearchDocsToCDNSSLKeys(docs []*riak.SearchDoc) []tc.CDNSSLKey {
	keys := []tc.CDNSSLKey{}
	for _, doc := range docs {
		key := tc.CDNSSLKey{}
		if dss := doc.Fields["deliveryservice"]; len(dss) > 0 {
			key.DeliveryService = dss[0]
		}
		if hosts := doc.Fields["hostname"]; len(hosts) > 0 {
			key.HostName = hosts[0]
		}
		if crts := doc.Fields["certificate.crt"]; len(crts) > 0 {
			key.Certificate.Crt = crts[0]
		}
		if keys := doc.Fields["certificate.key"]; len(keys) > 0 {
			key.Certificate.Key = keys[0]
		}
		keys = append(keys, key)
	}
	return keys
}

// DeleteOldDeliveryServiceSSLKeys deletes all the SSL keys in Riak for delivery services in the given CDN that are not in the given existingXMLIDs.
func DeleteOldDeliveryServiceSSLKeys(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, cdn tc.CDNName, existingXMLIDs map[string]struct{}) error {
	dsVersions := map[string][]string{}
	err := WithCluster(tx, authOpts, riakPort, func(cluster StorageCluster) error {
		// get the deliveryservice ssl keys by xmlID and version
		query := `cdn:` + string(cdn)
		filterQuery := ""
		fields := []string{"_yz_rk", "deliveryservice"} // '_yz_rk' is the magic Riak field that populates the key. Without this, doc.Key would be empty.
		searchDocs, err := Search(cluster, SSLKeysIndex, query, filterQuery, CDNSSLKeysLimit, fields)
		if err != nil {
			return errors.New("riak search error: " + err.Error())
		}
		if len(searchDocs) == 0 {
			return nil // no error, and leave keys empty
		}

		for _, doc := range searchDocs {
			dses := doc.Fields["deliveryservice"]
			if len(dses) == 0 {
				log.Errorln("Riak had a CDN '" + string(cdn) + "' key with no delivery service '" + doc.Key + "' - ignoring!")
				continue
			}
			if len(dses) > 1 {
				log.Errorf("Riak had a CDN '"+string(cdn)+"' key with multiple delivery services '"+doc.Key+"' deliveryservices '%+v' - ignoring all but the first!\n", dses)
			}
			ds := dses[0]

			dsVersions[ds] = append(dsVersions[ds], doc.Key)
		}
		return nil
	})
	if err != nil {
		return errors.New("with cluster error: " + err.Error())
	}

	successes := []string{}
	failures := []string{}
	for ds, riakKeys := range dsVersions {
		if _, ok := existingXMLIDs[ds]; ok {
			continue
		}
		for _, riakKey := range riakKeys {
			err := DeleteDeliveryServicesSSLKey(tx, authOpts, riakPort, riakKey)
			if err != nil {
				log.Errorln("deleting Traffic Vault SSL keys for Delivery Service '" + ds + "' key '" + riakKey + "': " + err.Error())
				failures = append(failures, ds)
			} else {
				log.Infoln("Deleted Traffic Vault SSL keys for delivery service which has been deleted in the database '" + string(ds) + "' key '" + riakKey + "'")
				successes = append(successes, ds)
			}
		}
	}
	if len(failures) > 0 {
		return errors.New("successfully deleted Traffic Vault SSL keys for deleted dses [" + strings.Join(successes, ", ") + "], but failed to delete Traffic Vault SSL keys for [" + strings.Join(failures, ", ") + "]; see the error log for details")
	}
	return nil
}
