package main

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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat/go-jwx/jwk"
)

// RiakPort is the port RIAK is listening on.
const RiakPort = 8087

// CDNURIKeysBucket is the namespace or bucket used for CDN URI signing keys.
const CDNURIKeysBucket = "cdn_uri_sig_keys"

// SSLKeysBucket ...
const SSLKeysBucket = "ssl"

// 5 second timeout
const timeOut = time.Second * 5

// MaxCommandExecutionAttempts ...
const MaxCommandExecutionAttempts = 5

// StorageCluster ...
type StorageCluster interface {
	Start() error
	Stop() error
	Execute(riak.Command) error
}

// RiakStorageCluster ...
type RiakStorageCluster struct {
	Cluster *riak.Cluster
}

// Stop ...
func (ri RiakStorageCluster) Stop() error {
	return ri.Cluster.Stop()
}

// Start ...
func (ri RiakStorageCluster) Start() error {
	return ri.Cluster.Start()
}

// Execute ...
func (ri RiakStorageCluster) Execute(command riak.Command) error {
	return ri.Cluster.Execute(command)
}

// URISignerKeyset is the container for the CDN URI signing keys
type URISignerKeyset struct {
	RenewalKid *string               `json:"renewal_kid"`
	Keys       []jwk.EssentialHeader `json:"keys"`
}

// deletes an object from riak storage
func deleteObject(key string, bucket string, cluster StorageCluster) error {
	if cluster == nil {
		return errors.New("ERROR: No valid cluster on which to execute a command")
	}

	// build store command and execute.
	cmd, err := riak.NewDeleteValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		WithTimeout(timeOut).
		Build()
	if err != nil {
		return err
	}

	err = cluster.Execute(cmd)

	if err != nil {
		return err
	}

	return nil
}

// fetch an object from riak storage
func fetchObjectValues(key string, bucket string, cluster StorageCluster) ([]*riak.Object, error) {
	if cluster == nil {
		return nil, errors.New("ERROR: No valid cluster on which to execute a command")
	}
	// build the fetch command
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		WithTimeout(timeOut).
		Build()
	if err != nil {
		return nil, err
	}

	if err = cluster.Execute(cmd); err != nil {
		return nil, err
	}
	fvc := cmd.(*riak.FetchValueCommand)

	// no object found with given key
	if fvc.Response == nil || fvc.Response.IsNotFound {
		return nil, nil
	}
	return fvc.Response.Values, nil
}

// saves an object to riak storage
func saveObject(obj *riak.Object, bucket string, cluster StorageCluster) error {
	if cluster == nil {
		return errors.New("ERROR: No valid cluster on which to execute a command")
	}
	if obj == nil {
		return errors.New("ERROR: cannot save a nil object")
	}
	// build store command and execute.
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(bucket).
		WithContent(obj).
		WithTimeout(timeOut).
		Build()
	if err != nil {
		return err
	}
	err = cluster.Execute(cmd)
	if err != nil {
		return err
	}

	return nil
}

// returns a riak cluster of online riak nodes.
func getRiakCluster(db *sqlx.DB, cfg Config) (StorageCluster, error) {
	riakServerQuery := `
		SELECT s.host_name, s.domain_name FROM server s 
		INNER JOIN type t on s.type = t.id 
		INNER JOIN status st on s.status = st.id 
		WHERE t.name = 'RIAK' AND st.name = 'ONLINE'
		`

	if cfg.RiakAuthOptions == nil {
		return nil, errors.New("ERROR: no riak auth information from riak.conf, cannot authenticate to any riak servers")
	}

	var nodes []*riak.Node
	rows, err := db.Query(riakServerQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s tc.Server
		var n *riak.Node
		if err := rows.Scan(&s.HostName, &s.DomainName); err != nil {
			return nil, err
		}
		addr := fmt.Sprintf("%s.%s:%d", s.HostName, s.DomainName, RiakPort)
		nodeOpts := &riak.NodeOptions{
			RemoteAddress: addr,
			AuthOptions:   cfg.RiakAuthOptions,
		}
		nodeOpts.AuthOptions.TlsConfig.ServerName = fmt.Sprintf("%s.%s", s.HostName, s.DomainName)
		n, err := riak.NewNode(nodeOpts)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	if len(nodes) == 0 {
		return nil, errors.New("ERROR: no available riak servers")
	}

	opts := &riak.ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: MaxCommandExecutionAttempts,
	}

	cluster, err := riak.NewCluster(opts)

	return RiakStorageCluster{Cluster: cluster}, err
}

// validates URISigingKeyset json.
func validateURIKeyset(msg map[string]URISignerKeyset) error {
	var renewalKidFound int
	var renewalKidMatched = false

	for key, value := range msg {
		issuer := key
		renewalKid := value.RenewalKid
		if issuer == "" {
			return errors.New("JSON Keyset has no issuer")
		}

		if renewalKid != nil {
			renewalKidFound++
		}

		for _, skey := range value.Keys {
			if skey.Algorithm == "" {
				return errors.New("A Key has no algorithm, alg, specified")
			}
			if skey.KeyID == "" {
				return errors.New("A Key has no key id, kid, specified")
			}
			if renewalKid != nil && strings.Compare(*renewalKid, skey.KeyID) == 0 {
				renewalKidMatched = true
			}
		}
	}

	// should only have one renewal_kid
	switch renewalKidFound {
	case 0:
		return errors.New("No renewal_kid was found in any keyset")
	case 1: // okay, this is what we want
		break
	default:
		return errors.New("More than one renewal_kid was found in the keysets")
	}

	// the renewal_kid should match the kid of one key
	if !renewalKidMatched {
		return errors.New("No key was found with a kid that matches the renewal kid")
	}

	return nil
}
