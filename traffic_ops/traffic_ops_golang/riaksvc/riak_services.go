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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"

	"github.com/basho/riak-go-client"
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

// PingCluster pings the given Riak cluster, and returns nil on success, or any error
func PingCluster(cluster StorageCluster) error {
	if cluster == nil {
		return errors.New("ERROR: No valid cluster on which to execute a command")
	}
	pingCommandBuilder := riak.PingCommandBuilder{}
	iCmd, err := pingCommandBuilder.Build()
	if err != nil {
		return errors.New("building riak ping command: " + err.Error())
	}
	if err := cluster.Execute(iCmd); err != nil {
		return errors.New("executing riak ping command: " + err.Error())
	}
	cmd, ok := iCmd.(*riak.PingCommand)
	if !ok {
		return fmt.Errorf("unexpected riak command type: %T", iCmd)
	}
	if err := cmd.Error(); err != nil {
		return errors.New("riak ping command returned error: " + err.Error())
	}
	if !cmd.Success() {
		return errors.New("riak ping command returned failure, but no error")
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

type ServerAddr struct {
	FQDN string
	Port string
}

func GetRiakServers(tx *sql.Tx) ([]ServerAddr, error) {
	rows, err := tx.Query(`
SELECT CONCAT(s.host_name, '.', s.domain_name) FROM server s
JOIN type t ON s.type = t.id
JOIN status st ON s.status = st.id
WHERE t.name = 'RIAK' AND st.name = 'ONLINE'
`)
	if err != nil {
		return nil, errors.New("querying riak servers: " + err.Error())
	}
	defer rows.Close()
	servers := []ServerAddr{}
	portStr := strconv.Itoa(RiakPort)
	for rows.Next() {
		s := ServerAddr{Port: portStr}
		if err := rows.Scan(&s.FQDN); err != nil {
			return nil, errors.New("scanning riak servers: " + err.Error())
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func RiakServersToCluster(servers []ServerAddr, authOptions *riak.AuthOptions) (StorageCluster, error) {
	if authOptions == nil {
		return nil, errors.New("ERROR: no riak auth information from riak.conf, cannot authenticate to any riak servers")
	}
	nodes := []*riak.Node{}
	for _, srv := range servers {
		nodeOpts := &riak.NodeOptions{
			RemoteAddress: srv.FQDN + ":" + srv.Port,
			AuthOptions:   authOptions,
		}
		nodeOpts.AuthOptions.TlsConfig.ServerName = srv.FQDN
		node, err := riak.NewNode(nodeOpts)
		if err != nil {
			return nil, errors.New("creating riak node: " + err.Error())
		}
		nodes = append(nodes, node)
	}
	if len(nodes) == 0 {
		return nil, errors.New("ERROR: no available riak servers")
	}
	opts := &riak.ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: MaxCommandExecutionAttempts,
	}
	cluster, err := riak.NewCluster(opts)
	if err != nil {
		return nil, errors.New("creating riak cluster: " + err.Error())
	}
	return RiakStorageCluster{Cluster: cluster}, nil
}

func GetRiakClusterTx(tx *sql.Tx, authOptions *riak.AuthOptions) (StorageCluster, error) {
	servers, err := GetRiakServers(tx)
	if err != nil {
		return nil, errors.New("getting riak servers: " + err.Error())
	}
	cluster, err := RiakServersToCluster(servers, authOptions)
	if err != nil {
		return nil, errors.New("creating riak cluster from servers: " + err.Error())
	}
	return cluster, nil
}

func WithClusterTx(tx *sql.Tx, authOpts *riak.AuthOptions, f func(StorageCluster) error) error {
	cluster, err := GetRiakClusterTx(tx, authOpts)
	if err != nil {
		return errors.New("getting riak cluster: " + err.Error())
	}
	if err = cluster.Start(); err != nil {
		return errors.New("starting riak cluster: " + err.Error())
	}
	defer func() {
		if err := cluster.Stop(); err != nil {
			log.Errorln("error stopping Riak cluster: " + err.Error())
		}
	}()
	return f(cluster)
}

// StartCluster gets and starts a riak cluster, returning an error if either getting or starting fails.
func StartCluster(tx *sql.Tx, authOptions *riak.AuthOptions) (StorageCluster, error) {
	cluster, err := GetRiakClusterTx(tx, authOptions)
	if err != nil {
		return nil, errors.New("getting cluster: " + err.Error())
	}
	if err = cluster.Start(); err != nil {
		return nil, errors.New("starting cluster: " + err.Error())
	}
	return cluster, nil
}

// StopCluster stops the cluster, logging any error rather than returning it. This is designed to be called in a defer.
func StopCluster(c StorageCluster) {
	if err := c.Stop(); err != nil {
		log.Errorln("stopping riak cluster: " + err.Error())
	}
}

// Search searches Riak for the given query. Returns nil and a nil error if no object was found.
func Search(cluster StorageCluster, index string, query string, filterQuery string, numRows int) ([]*riak.SearchDoc, error) {
	iCmd, err := riak.NewSearchCommandBuilder().
		WithIndexName(index).
		WithQuery(query).
		WithFilterQuery(filterQuery).
		WithNumRows(uint32(numRows)).
		Build()
	if err != nil {
		return nil, errors.New("building Riak command: " + err.Error())
	}
	if err = cluster.Execute(iCmd); err != nil {
		return nil, errors.New("executing Riak command index '" + index + "' query '" + query + "': " + err.Error())
	}
	cmd, ok := iCmd.(*riak.SearchCommand)
	if !ok {
		return nil, fmt.Errorf("Riak command unexpected type %T", iCmd)
	}
	if cmd.Response == nil || cmd.Response.NumFound == 0 {
		return nil, nil
	}
	return cmd.Response.Docs, nil
}
