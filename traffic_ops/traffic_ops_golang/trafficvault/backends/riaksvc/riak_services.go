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
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/basho/riak-go-client"
)

const (
	// defaultRiakPort is the port RIAK is listening on, if no port is configured.
	defaultRiakPort                    = uint(8087)
	defaultTimeOut                     = time.Second * 5
	defaultHealthCheckInterval         = time.Second * 5
	defaultMaxCommandExecutionAttempts = 5
)

var (
	clusterServers []ServerAddr
	sharedCluster  *riak.Cluster
	clusterMutex   sync.Mutex

	healthCheckInterval time.Duration
)

type TOAuthOptions struct {
	riak.AuthOptions
	MaxTLSVersion *string
	Port          uint `json:"port"`
}

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

func setMaxTLSVersion(riakConfig *TOAuthOptions) error {
	if riakConfig.MaxTLSVersion == nil {
		if riakConfig.TlsConfig.MaxVersion == 0 {
			riakConfig.TlsConfig.MaxVersion = tls.VersionTLS11
		}
		return nil
	}
	tlsVersions := map[string]uint16{
		"1.0": tls.VersionTLS10,
		"1.1": tls.VersionTLS11,
		"1.2": tls.VersionTLS12,
		"1.3": tls.VersionTLS13,
	}
	var err error
	if version, exists := tlsVersions[*riakConfig.MaxTLSVersion]; exists {
		riakConfig.TlsConfig.MaxVersion = version
	} else {
		err = fmt.Errorf("%v is not a valid TLS version", riakConfig.MaxTLSVersion)
	}
	return err
}

type Config struct {
	riak.AuthOptions
	Port uint
}

func unmarshalRiakConfig(riakConfBytes json.RawMessage) (Config, error) {
	conf := Config{}
	rconf := &TOAuthOptions{}
	rconf.TlsConfig = &tls.Config{}
	err := json.Unmarshal(riakConfBytes, &rconf)
	if err != nil {
		return conf, err
	}
	if err := setMaxTLSVersion(rconf); err != nil {
		return conf, err
	}

	type config struct {
		Hci string `json:"HealthCheckInterval"`
	}

	var checkconfig config
	err = json.Unmarshal(riakConfBytes, &checkconfig)
	if err == nil {
		hci, _ := time.ParseDuration(checkconfig.Hci)
		if 0 < hci {
			healthCheckInterval = hci
		}
	} else {
		log.Infoln("Error unmarshalling riak config options: " + err.Error())
	}

	if healthCheckInterval <= 0 {
		healthCheckInterval = defaultHealthCheckInterval
		log.Infoln("HeathCheckInterval override")
	}

	log.Infoln("Riak health check interval set to:", healthCheckInterval)

	conf.AuthOptions = rconf.AuthOptions
	conf.Port = rconf.Port
	if conf.Port == 0 {
		conf.Port = defaultRiakPort
	}
	return conf, nil
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
		WithTimeout(defaultTimeOut).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}

// pingCluster pings the given Riak cluster, and returns nil on success, or any error
func pingCluster(cluster StorageCluster) error {
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
func fetchObjectValues(key string, bucket string, cluster StorageCluster) ([]*riak.Object, error) {
	if cluster == nil {
		return nil, errors.New("ERROR: No valid cluster on which to execute a command")
	}
	// build the fetch command
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		WithTimeout(defaultTimeOut).
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
		WithTimeout(defaultTimeOut).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}

type ServerAddr struct {
	FQDN string
	Port string
}

// getRiakServers returns the riak servers from the database. The riakPort may be nil, in which case the default port is returned.
func getRiakServers(tx *sql.Tx, riakPort *uint) ([]ServerAddr, error) {
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
	if riakPort == nil {
		riakPort = util.UIntPtr(defaultRiakPort)
	}
	portStr := strconv.Itoa(int(*riakPort))
	for rows.Next() {
		s := ServerAddr{Port: portStr}
		if err := rows.Scan(&s.FQDN); err != nil {
			return nil, errors.New("scanning riak servers: " + err.Error())
		}
		servers = append(servers, s)
	}

	return servers, nil
}

func getRiakCluster(servers []ServerAddr, authOptions *riak.AuthOptions) (*riak.Cluster, error) {
	if authOptions == nil {
		return nil, errors.New("ERROR: no riak auth information from riak.conf, cannot authenticate to any riak servers")
	}
	nodes := []*riak.Node{}
	for _, srv := range servers {
		nodeOpts := &riak.NodeOptions{
			RemoteAddress:       srv.FQDN + ":" + srv.Port,
			AuthOptions:         authOptions,
			HealthCheckInterval: healthCheckInterval,
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
		ExecutionAttempts: defaultMaxCommandExecutionAttempts,
	}
	cluster, err := riak.NewCluster(opts)
	if err != nil {
		return nil, errors.New("creating riak cluster: " + err.Error())
	}
	return cluster, err
}

func getRiakStorageCluster(servers []ServerAddr, authOptions *riak.AuthOptions) (StorageCluster, error) {
	cluster, err := getRiakCluster(servers, authOptions)
	if err != nil {
		return nil, err
	}
	return RiakStorageCluster{Cluster: cluster}, nil
}

func getPooledCluster(tx *sql.Tx, authOptions *riak.AuthOptions, riakPort *uint) (StorageCluster, error) {
	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	tryLoad := false

	// should we try to reload the cluster?
	newservers, err := getRiakServers(tx, riakPort)

	if err == nil {
		if 0 < len(newservers) {
			sort.Slice(newservers, func(ii, jj int) bool {
				return newservers[ii].FQDN < newservers[jj].FQDN ||
					(newservers[ii].FQDN == newservers[jj].FQDN && newservers[ii].Port < newservers[jj].Port)
			})
			if !reflect.DeepEqual(newservers, clusterServers) {
				tryLoad = true
				log.Infoln("Attempting to load a new set of riak servers")
				log.Infoln("new riak servers")
				for _, srv := range newservers {
					log.Infoln(" ", srv.FQDN+":"+srv.Port)
				}
			}
		}
	} else {
		log.Errorln("getting riak servers: " + err.Error())
	}

	if tryLoad {
		newcluster, err := getRiakCluster(newservers, authOptions)
		if err == nil {
			if err := newcluster.Start(); err == nil {
				log.Infof("New riak cluster started: %p\n", newcluster)

				if sharedCluster != nil {
					runtime.SetFinalizer(sharedCluster, func(c *riak.Cluster) {
						log.Infof("running finalizer for riak sharedcluster (%p)\n", c)
						if err := c.Stop(); err != nil {
							log.Errorf("in finalizer for riak sharedcluster (%p): stopping cluster: %s\n", c, err.Error())
						}
					})
				}

				sharedCluster = newcluster
				clusterServers = newservers
			} else {
				log.Errorln("starting riak cluster, reverting to previous: " + err.Error())
			}
		} else {
			log.Errorln("creating riak cluster, reverting to previous: " + err.Error())
		}
	}

	cluster := sharedCluster

	if cluster == nil {
		log.Errorln("getPooledCluster failed, returning nil cluster")
		return nil, errors.New("getPooledCluster unable to return cluster")
	}

	return RiakStorageCluster{Cluster: cluster}, nil
}

func withCluster(tx *sql.Tx, authOpts *riak.AuthOptions, riakPort *uint, f func(StorageCluster) error) error {
	cluster, err := getPooledCluster(tx, authOpts, riakPort)
	if err != nil {
		return errors.New("getting riak pooled cluster: " + err.Error())
	}
	return f(cluster)
}

// search searches Riak for the given query. Returns nil and a nil error if no object was found.
// If fields is empty, all fields will be returned.
func search(cluster StorageCluster, index string, query string, filterQuery string, numRows uint32, fields []string) ([]*riak.SearchDoc, error) {
	var searchDocs []*riak.SearchDoc
	for start := uint32(0); ; start += numRows {
		riakCmd := riak.NewSearchCommandBuilder().
			WithIndexName(index).
			WithQuery(query).
			WithNumRows(numRows).
			WithStart(start)
		if len(filterQuery) > 0 {
			riakCmd = riakCmd.WithFilterQuery(filterQuery)
		}
		if len(fields) > 0 {
			riakCmd = riakCmd.WithReturnFields(fields...)
		}
		iCmd, err := riakCmd.Build()

		if err != nil {
			return nil, errors.New("building Riak command: " + err.Error())
		}
		if err = cluster.Execute(iCmd); err != nil {
			return nil, errors.New("executing Riak command index '" + index + "' query '" + query + "': " + err.Error())
		}
		cmd, ok := iCmd.(*riak.SearchCommand)
		if !ok {
			return nil, fmt.Errorf("riak command unexpected type %T", iCmd)
		}
		if cmd.Response == nil {
			return nil, fmt.Errorf("riak received nil response")
		}
		if start == 0 {
			if cmd.Response.NumFound <= numRows {
				return cmd.Response.Docs, nil
			} else {
				searchDocs = make([]*riak.SearchDoc, cmd.Response.NumFound)
			}
		}

		// If the total number of docs is not evenly divisible by numRows
		if uint32(len(cmd.Response.Docs)) < numRows {
			numRows = uint32(len(cmd.Response.Docs))
		}

		for responseIndex := uint32(0); responseIndex < numRows; responseIndex += 1 {
			returnIndex := responseIndex + start
			searchDocs[returnIndex] = cmd.Response.Docs[responseIndex]
		}
		if cmd.Response.NumFound == numRows+start {
			return searchDocs, nil
		}
	}
}
