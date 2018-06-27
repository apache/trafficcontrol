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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/basho/riak-go-client"
)

const DeliveryServiceSSLKeysBucket = "ssl"
const DNSSECKeysBucket = "dnssec"

func GetDeliveryServiceSSLKeysObj(xmlID string, version string, tx *sql.Tx, authOpts *riak.AuthOptions) (tc.DeliveryServiceSSLKeys, bool, error) {
	key := tc.DeliveryServiceSSLKeys{}
	if version == "" {
		xmlID += "-latest"
	} else {
		xmlID += "-" + version
	}
	found := false
	err := WithClusterTx(tx, authOpts, func(cluster StorageCluster) error {
		// get the deliveryservice ssl keys by xmlID and version
		ro, err := FetchObjectValues(xmlID, DeliveryServiceSSLKeysBucket, cluster)
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

func GetDeliveryServiceSSLKeysObjTx(xmlID string, version string, tx *sql.Tx, authOpts *riak.AuthOptions) (tc.DeliveryServiceSSLKeys, bool, error) {
	key := tc.DeliveryServiceSSLKeys{}
	if version == "" {
		xmlID += "-latest"
	} else {
		xmlID += "-" + version
	}
	found := false
	err := WithClusterTx(tx, authOpts, func(cluster StorageCluster) error {
		// get the deliveryservice ssl keys by xmlID and version
		ro, err := FetchObjectValues(xmlID, DeliveryServiceSSLKeysBucket, cluster)
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

func PutDeliveryServiceSSLKeysObj(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, authOpts *riak.AuthOptions) error {
	keyJSON, err := json.Marshal(&key)
	if err != nil {
		return errors.New("marshalling key: " + err.Error())
	}
	err = WithClusterTx(tx, authOpts, func(cluster StorageCluster) error {
		obj := &riak.Object{
			ContentType:     "text/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             key.DeliveryService,
			Value:           []byte(keyJSON),
		}
		if err = SaveObject(obj, DeliveryServiceSSLKeysBucket, cluster); err != nil {
			return errors.New("saving Riak object: " + err.Error())
		}
		return nil
	})
	return err
}

func PutDeliveryServiceSSLKeysObjTx(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, authOpts *riak.AuthOptions) error {
	keyJSON, err := json.Marshal(&key)
	if err != nil {
		return errors.New("marshalling key: " + err.Error())
	}
	err = WithClusterTx(tx, authOpts, func(cluster StorageCluster) error {
		obj := &riak.Object{
			ContentType:     "text/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             key.DeliveryService,
			Value:           []byte(keyJSON),
		}
		if err = SaveObject(obj, DeliveryServiceSSLKeysBucket, cluster); err != nil {
			return errors.New("saving Riak object: " + err.Error())
		}
		return nil
	})
	return err
}

func GetDNSSECKeys(cdnName string, tx *sql.Tx, authOpts *riak.AuthOptions) (tc.DNSSECKeys, bool, error) {
	key := tc.DNSSECKeys{}
	found := false
	log.Errorln("riaksvc.GetDNSSECKeys calling")
	err := WithClusterTx(tx, authOpts, func(cluster StorageCluster) error {
		log.Errorln("riaksvc.GetDNSSECKeys in WithClusterTx")
		ro, err := FetchObjectValues(cdnName, DNSSECKeysBucket, cluster)
		log.Errorln("riaksvc.GetDNSSECKeys fetched object values")
		if err != nil {
			log.Errorln("riaksvc.GetDNSSECKeys fetched object values returning err")
			return err
		}
		if len(ro) == 0 {
			log.Errorln("riaksvc.GetDNSSECKeys returning nil, len(ro) is 0")
			return nil // not found
		}
		log.Errorln("riaksvc.GetDNSSECKeys unmarshalling")
		if err := json.Unmarshal(ro[0].Value, &key); err != nil {
			log.Errorln("Unmarshaling Riak dnssec response: " + err.Error())
			return errors.New("unmarshalling Riak dnssec response: " + err.Error())
		}
		log.Errorln("riaksvc.GetDNSSECKeys unmarshalled, found true, returning nil err")
		found = true
		return nil
	})
	log.Errorln("riaksvc.GetDNSSECKeys out of WithCluster")
	if err != nil {
		log.Errorln("riaksvc.GetDNSSECKeys WithCluster err, returning err")
		return key, false, err
	}
	log.Errorln("riaksvc.GetDNSSECKeys returning success")
	return key, found, nil
}

func PutDNSSECKeys(keys tc.DNSSECKeys, cdnName string, tx *sql.Tx, authOpts *riak.AuthOptions) error {
	keyJSON, err := json.Marshal(&keys)
	if err != nil {
		return errors.New("marshalling keys: " + err.Error())
	}
	err = WithClusterTx(tx, authOpts, func(cluster StorageCluster) error {
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
