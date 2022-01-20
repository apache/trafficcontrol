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
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	riak "github.com/basho/riak-go-client"
)

var vault_ip string
var vault_port uint
var vault_user string
var vault_pass string
var vault_action string
var dry_run bool
var insecure bool

func connectToRiak(vault_ip string, vault_port uint, insecure bool) *riak.Cluster {

	tlsConfig := tls.Config{
		ServerName:         vault_ip,
		InsecureSkipVerify: insecure,
	}

	authOptions := riak.AuthOptions{
		User:      vault_user,
		Password:  vault_pass,
		TlsConfig: &tlsConfig,
	}

	vaultAddr := fmt.Sprintf("%s:%d", vault_ip, vault_port)
	nodeOpts := &riak.NodeOptions{
		RemoteAddress: vaultAddr,
		AuthOptions:   &authOptions,
	}

	log.Printf("Connecting to %s", vaultAddr)
	var node *riak.Node
	var err error
	if node, err = riak.NewNode(nodeOpts); err != nil {
		log.Fatal(err.Error())
	}

	nodes := []*riak.Node{node}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}

	cluster, err := riak.NewCluster(opts)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := cluster.Start(); err != nil {
		log.Fatal(err.Error())
	}

	return cluster
}

func listBuckets(cluster *riak.Cluster) []string {
	log.Print("Listing Riak buckets")
	cmd, err := riak.NewListBucketsCommandBuilder().
		WithAllowListing().
		Build()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := cluster.Execute(cmd); err != nil {
		log.Fatal(err)
	}

	lbc := cmd.(*riak.ListBucketsCommand)
	rsp := lbc.Response
	return rsp.Buckets
}

func listKeys(cluster *riak.Cluster, bucket string) []string {
	cmd, err := riak.NewListKeysCommandBuilder().
		WithAllowListing().
		WithBucket(bucket).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	if err := cluster.Execute(cmd); err != nil {
		log.Fatal(err)
	}

	lkc := cmd.(*riak.ListKeysCommand)
	rsp := lkc.Response
	return rsp.Keys
}

func getValue(cluster *riak.Cluster, bucket string, key string) ([]byte, bool) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	if err := cluster.Execute(cmd); err != nil {
		log.Fatal(err)
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response

	if rsp.IsNotFound {
		log.Print("ERROR Key not found: ", bucket, ":", key)
		return nil, false
	}

	return rsp.Values[0].Value, true

}

type SSLRecord struct {
	Country string `json:"country"`
	Cdn     string `json:"cdn"`
	XmlId   string `json:"deliveryservice"`
	Org     string `json:"org"`
}

// Converts Riak keys of form "ds_<id#>-<version>" to new form "<xmlid>-<version>"
func convertSSLToXmlID(cluster *riak.Cluster) {
	bucket := "ssl"
	keys := listKeys(cluster, bucket)

	stagedRecords := make([]riak.Object, 0)

	allFound := true
	for _, key := range keys {
		value, found := getValue(cluster, bucket, key)
		if found {
			log.Printf("[Read Record] Bucket: %s, Key: %s, Value: %s...", bucket, key, value[0:75])
			if !strings.HasPrefix(key, "ds") {
				continue
			}

			splitKey := strings.Split(key, "-")
			version := splitKey[len(splitKey)-1]

			sslRecord := SSLRecord{}
			json.Unmarshal(value, &sslRecord)

			newKey := fmt.Sprintf("%s-%s", sslRecord.XmlId, version)

			newObj := riak.Object{Bucket: bucket,
				Key:         newKey,
				ContentType: "application/json",
				Value:       value}

			stagedRecords = append(stagedRecords, newObj)

		} else {
			allFound = false
		}
	}

	if !allFound {
		log.Print("Some keys are missing, please correct and retry. Exiting!")
		os.Exit(1)
	}

	log.Print("Inserting new renamed records")
	for _, record := range stagedRecords {
		cmd, err := riak.NewStoreValueCommandBuilder().
			WithContent(&record).
			Build()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("[Write Record] Bucket: %s, Key: %s, Value: %s...",
			record.Bucket, record.Key, string(record.Value)[0:75])

		if dry_run {
			continue
		}

		if err := cluster.Execute(cmd); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func init() {
	flag.StringVar(&vault_ip, "vault_ip", "", "IP/Hostname of Vault")
	flag.UintVar(&vault_port, "vault_port", 8087, "Protobuffers port of Vault")
	flag.StringVar(&vault_user, "vault_user", "", "Riak Username")
	flag.StringVar(&vault_pass, "vault_password", "", "Riak Password")
	flag.StringVar(&vault_action, "vault_action", "", "Action: list_buckets|list_keys|list_values|convert_ssl_to_xmlid")
	flag.BoolVar(&dry_run, "dry_run", false, "Do not perform writes")
	flag.BoolVar(&insecure, "insecure", false, "Disable TLS certificate checks when connecting to cluster. Defaults to false")
}

func main() {
	log.Print("Traffic Control Traffic Vault Util")
	flag.Parse()

	if dry_run {
		log.Print("---- DRY RUN --- ")
	}

	if vault_ip == "" {
		log.Fatal("Must provide Traffic Vault IP or host")
	}

	cluster := connectToRiak(vault_ip, vault_port, insecure)
	defer func() {
		if err := cluster.Stop(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	switch vault_action {
	case "list_buckets":
		buckets := listBuckets(cluster)
		log.Print("Buckets: ", buckets)

	case "list_keys":
		buckets := listBuckets(cluster)
		for _, bucket := range buckets {
			keys := listKeys(cluster, bucket)
			for _, key := range keys {
				log.Printf("Bucket: %s, Key: %s", bucket, key)
			}
		}

	case "list_values":
		buckets := listBuckets(cluster)
		for _, bucket := range buckets {
			keys := listKeys(cluster, bucket)
			for _, key := range keys {
				value, found := getValue(cluster, bucket, key)
				if found {
					log.Printf("Bucket: %s, Key: %s, Value: %s", bucket, key, string(value))
				} else {
					log.Printf("Bucket: %s, Key: %s, NOT FOUND", bucket, key)
				}
			}
		}

	case "convert_ssl_to_xmlid":
		convertSSLToXmlID(cluster)

	default:
		log.Print("Unknown vault_action: ", vault_action)
		log.Print("Allowed actions: list_buckets|list_keys|list_values|convert_ssl_to_xmlid")
		os.Exit(1)

	}
}
