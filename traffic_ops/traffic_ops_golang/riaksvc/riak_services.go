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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
)

// RiakPort is the port RIAK is listening on.
const RiakPort = 8087

// 5 second timeout
const timeOut = time.Second * 5

// MaxCommandExecutionAttempts ...
const MaxCommandExecutionAttempts = 5

type AuthOptions riak.AuthOptions

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

func GetRiakConfig(riakConfigFile string) (bool, *riak.AuthOptions, error) {
	riakConfBytes, err := ioutil.ReadFile(riakConfigFile)
	if err != nil {
		return false, nil, fmt.Errorf("reading riak conf '%v': %v", riakConfigFile, err)
	}

	rconf := &riak.AuthOptions{}
	rconf.TlsConfig = &tls.Config{}
	err = json.Unmarshal([]byte(riakConfBytes), &rconf)
	if err != nil {
		return false, nil, fmt.Errorf("Unmarshalling riak conf '%v': %v", riakConfigFile, err)
	}

	return true, rconf, nil
}

// deletes an object from riak storage
func DeleteObject(key string, bucket string, cluster StorageCluster) error {
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
func FetchObjectValues(key string, bucket string, cluster StorageCluster) ([]*riak.Object, error) {
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
func SaveObject(obj *riak.Object, bucket string, cluster StorageCluster) error {
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
func GetRiakCluster(db *sqlx.DB, authOptions *riak.AuthOptions) (StorageCluster, error) {
	riakServerQuery := `
		SELECT s.host_name, s.domain_name FROM server s
		INNER JOIN type t on s.type = t.id
		INNER JOIN status st on s.status = st.id
		WHERE t.name = 'RIAK' AND st.name = 'ONLINE'
		`

	if authOptions == nil {
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
			AuthOptions:   authOptions,
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
