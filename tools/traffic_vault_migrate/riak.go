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
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/basho/riak-go-client"
	"log"
	"reflect"
	"time"
)

const (
	BUCKET_SSL     = "ssl"
	BUCKET_DNSSEC  = "dnssec"
	BUCKET_URL_SIG = "url_sig_keys"
	BUCKET_URI_SIG = "cdn_uri_sig_keys"

	INDEX_SSL = "sslkeys"

	SCHEMA_RIAK_KEY    = "_yz_rk"
	SCHEMA_RIAK_BUCKET = "_yz_rb"
)

var (
	SCHEMA_SSL_FIELDS = [...]string{SCHEMA_RIAK_KEY, SCHEMA_RIAK_BUCKET}
)

type RiakConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	VerifyTLS  bool   `json:"tls"`
	TLSVersion string `json:"tlsVersion"`
}
type RiakBackend struct {
	sslKeys     RiakSSLKeyTable
	dnssecKeys  RiakDNSSecKeyTable
	uriSignKeys RiakURISignKeyTable
	urlSigKeys  RiakURLSigKeyTable
	cfg         RiakConfig
	cluster     *riak.Cluster
}

func (rb *RiakBackend) String() string {
	data := fmt.Sprintf("Riak server %v@%v:%v\n", rb.cfg.User, rb.cfg.Host, rb.cfg.Port)
	data += fmt.Sprintf("\tSSL Keys: %v\n", len(rb.sslKeys.Records))
	data += fmt.Sprintf("\tDNSSec Keys: %v\n", len(rb.dnssecKeys.Records))
	data += fmt.Sprintf("\tURI Keys: %v\n", len(rb.uriSignKeys.Records))
	data += fmt.Sprintf("\tURL Keys: %v\n", len(rb.urlSigKeys.Records))
	return data
}
func (rb *RiakBackend) Name() string {
	return "Riak"
}
func (rb *RiakBackend) ReadConfig(s string) error {
	cfgGeneric, err := UnmarshalConfig(s, reflect.TypeOf(rb.cfg))
	if err != nil {
		return err
	}

	rb.cfg = *cfgGeneric.Interface().(*RiakConfig)
	return nil
}
func (rb *RiakBackend) Insert() error {
	if err := rb.sslKeys.insertKeys(rb.cluster); err != nil {
		return err
	}
	if err := rb.dnssecKeys.insertKeys(rb.cluster); err != nil {
		return err
	}
	if err := rb.urlSigKeys.insertKeys(rb.cluster); err != nil {
		return err
	}
	if err := rb.uriSignKeys.insertKeys(rb.cluster); err != nil {
		return err
	}
	return nil
}
func (rb *RiakBackend) ValidateKey() []string {
	errors := []string{}
	if errs := rb.sslKeys.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	if errs := rb.dnssecKeys.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	if errs := rb.uriSignKeys.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	if errs := rb.urlSigKeys.validate(); errs != nil {
		errors = append(errors, errs...)
	}

	return errors
}
func (rb *RiakBackend) SetSSLKeys(keys []SSLKey) error {
	rb.sslKeys.fromGeneric(keys)
	return nil
}
func (rb *RiakBackend) SetDNSSecKeys(keys []DNSSecKey) error {
	rb.dnssecKeys.fromGeneric(keys)
	return nil
}
func (rb *RiakBackend) SetURISignKeys(keys []URISignKey) error {
	rb.uriSignKeys.fromGeneric(keys)
	return nil
}
func (rb *RiakBackend) SetURLSigKeys(keys []URLSigKey) error {
	rb.urlSigKeys.fromGeneric(keys)
	return nil
}
func (rb *RiakBackend) Start() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !rb.cfg.VerifyTLS,
		MaxVersion:         tls.VersionTLS11,
	}
	auth := &riak.AuthOptions{
		User:      rb.cfg.User,
		Password:  rb.cfg.Password,
		TlsConfig: tlsConfig,
	}

	cluster, err := getRiakCluster(rb.cfg, auth)
	if err != nil {
		return err
	}
	if err := cluster.Start(); err != nil {
		return err
	}

	rb.cluster = cluster
	rb.sslKeys = RiakSSLKeyTable{}
	rb.dnssecKeys = RiakDNSSecKeyTable{}
	rb.urlSigKeys = RiakURLSigKeyTable{}
	rb.uriSignKeys = RiakURISignKeyTable{}
	return nil
}
func (rb *RiakBackend) Stop() error {
	if err := rb.cluster.Stop(); err != nil {
		return err
	}
	return nil
}
func (rb *RiakBackend) Ping() error {
	return ping(rb.cluster)
}
func (rb *RiakBackend) GetSSLKeys() ([]SSLKey, error) {
	return rb.sslKeys.toGeneric(), nil
}
func (rb *RiakBackend) GetDNSSecKeys() ([]DNSSecKey, error) {
	return rb.dnssecKeys.toGeneric(), nil
}
func (rb *RiakBackend) GetURISignKeys() ([]URISignKey, error) {
	return rb.uriSignKeys.toGeneric(), nil
}
func (rb *RiakBackend) GetURLSigKeys() ([]URLSigKey, error) {
	return rb.urlSigKeys.toGeneric(), nil
}
func (rb *RiakBackend) Fetch() error {
	if err := rb.sslKeys.gatherKeys(rb.cluster); err != nil {
		return err
	}
	if err := rb.dnssecKeys.gatherKeys(rb.cluster); err != nil {
		return err
	}
	if err := rb.urlSigKeys.gatherKeys(rb.cluster); err != nil {
		return err
	}
	if err := rb.uriSignKeys.gatherKeys(rb.cluster); err != nil {
		return err
	}

	return nil
}

type RiakSSLKeyRecord struct {
	tc.DeliveryServiceSSLKeys
	CommonRecord
}
type RiakSSLKeyTable struct {
	Records []RiakSSLKeyRecord
}

func (tbl *RiakSSLKeyTable) gatherKeys(cluster *riak.Cluster) error {
	searchDocs, err := search(cluster, INDEX_SSL, "cdn:*", "", 1000, SCHEMA_SSL_FIELDS[:])
	if err != nil {
		return err
	}
	if len(searchDocs) == 0 {
		return errors.New("No ssl keys")
	}

	tbl.Records = make([]RiakSSLKeyRecord, len(searchDocs))
	for i, doc := range searchDocs {
		objs, err := getObject(cluster, doc.Bucket, doc.Key)
		if err != nil {
			return err
		}
		if len(objs) < 1 {
			return errors.New(fmt.Sprintf("Unable to find any objects with key %v and bucket %v, but search results were returned!", doc.Key, doc.Bucket))
		}
		if len(objs) > 1 {
			return errors.New(fmt.Sprintf("More than 1 ssl key record found %v\n", len(objs)))
		}
		var obj tc.DeliveryServiceSSLKeys
		if err = json.Unmarshal(objs[0].Value, &obj); err != nil {
			return err
		}
		tbl.Records[i] = RiakSSLKeyRecord{
			DeliveryServiceSSLKeys: obj,
			CommonRecord:           CommonRecord{},
		}
	}
	return nil
}
func (tbl *RiakSSLKeyTable) toGeneric() []SSLKey {
	keys := make([]SSLKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = SSLKey{
			DeliveryServiceSSLKeys: record.DeliveryServiceSSLKeys,
			CommonRecord:           record.CommonRecord,
		}
	}

	return keys
}
func (tbl *RiakSSLKeyTable) fromGeneric(keys []SSLKey) {
	tbl.Records = make([]RiakSSLKeyRecord, len(keys))

	for i, record := range keys {
		tbl.Records[i] = RiakSSLKeyRecord{
			DeliveryServiceSSLKeys: record.DeliveryServiceSSLKeys,
			CommonRecord:           record.CommonRecord,
		}
	}
}
func (tbl *RiakSSLKeyTable) insertKeys(cluster *riak.Cluster) error {
	for _, record := range tbl.Records {
		objBytes, err := json.Marshal(record.DeliveryServiceSSLKeys)
		if err != nil {
			return err
		}
		err = setObject(cluster, makeRiakObject(objBytes, record.DeliveryService+"-"+record.Version.String()), BUCKET_SSL)
		if err != nil {
			return err
		}
	}
	return nil
}
func (tbl *RiakSSLKeyTable) validate() []string {
	errs := []string{}
	for i, record := range tbl.Records {
		if record.DeliveryService == "" {
			errs = append(errs, fmt.Sprintf("SSL Key #%v: Delivery Service is blank!", i))
		}
		if record.CDN == "" {
			errs = append(errs, fmt.Sprintf("SSL Key #%v: CDN is blank!", i))
		}
	}
	return errs
}

type RiakDNSSecKeyRecord struct {
	CDN string
	Key tc.DNSSECKeysRiak
	CommonRecord
}
type RiakDNSSecKeyTable struct {
	Records []RiakDNSSecKeyRecord
}

func (tbl *RiakDNSSecKeyTable) gatherKeys(cluster *riak.Cluster) error {
	tbl.Records = []RiakDNSSecKeyRecord{}
	objs, err := getObjects(cluster, BUCKET_DNSSEC)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		key := tc.DNSSECKeysRiak{}
		if err := json.Unmarshal(obj.Value, &key); err != nil {
			return err
		}
		tbl.Records = append(tbl.Records, RiakDNSSecKeyRecord{
			CDN:          obj.Key,
			CommonRecord: CommonRecord{},
			Key:          key,
		})
	}
	return nil
}
func (tbl *RiakDNSSecKeyTable) toGeneric() []DNSSecKey {
	keys := make([]DNSSecKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = DNSSecKey{
			CDN:                    record.CDN,
			DNSSECKeysTrafficVault: tc.DNSSECKeysTrafficVault(record.Key),
			CommonRecord:           CommonRecord{},
		}
	}

	return keys
}
func (tbl *RiakDNSSecKeyTable) fromGeneric(keys []DNSSecKey) {
	tbl.Records = make([]RiakDNSSecKeyRecord, len(keys))

	for i, record := range keys {
		tbl.Records[i] = RiakDNSSecKeyRecord{
			CDN:          record.CDN,
			CommonRecord: record.CommonRecord,
			Key:          tc.DNSSECKeysRiak(record.DNSSECKeysTrafficVault),
		}
	}
}
func (tbl *RiakDNSSecKeyTable) insertKeys(cluster *riak.Cluster) error {
	for _, record := range tbl.Records {
		objBytes, err := json.Marshal(record.Key)
		if err != nil {
			return err
		}

		err = setObject(cluster, makeRiakObject(objBytes, record.CDN), BUCKET_DNSSEC)
		if err != nil {
			return err
		}
	}
	return nil
}
func (tbl *RiakDNSSecKeyTable) validate() []string {
	errs := []string{}
	for i, record := range tbl.Records {
		if record.CDN == "" {
			errs = append(errs, fmt.Sprintf("DNSSec Key #%v: CDN is blank!", i))
		}
	}
	return errs
}

type RiakURLSigKeyRecord struct {
	Key             tc.URLSigKeys
	DeliveryService string
	CommonRecord
}
type RiakURLSigKeyTable struct {
	Records []RiakURLSigKeyRecord
}

func (tbl *RiakURLSigKeyTable) gatherKeys(cluster *riak.Cluster) error {
	tbl.Records = []RiakURLSigKeyRecord{}
	objs, err := getObjects(cluster, BUCKET_URL_SIG)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		key := tc.URLSigKeys{}
		if err := json.Unmarshal(obj.Value, &key); err != nil {
			return err
		}
		tbl.Records = append(tbl.Records, RiakURLSigKeyRecord{
			DeliveryService: obj.Key,
			CommonRecord:    CommonRecord{},
			Key:             key,
		})
	}
	return nil
}
func (tbl *RiakURLSigKeyTable) toGeneric() []URLSigKey {
	keys := make([]URLSigKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = URLSigKey{
			URLSigKeys:      record.Key,
			DeliveryService: record.DeliveryService,
			CommonRecord:    record.CommonRecord,
		}
	}

	return keys
}
func (tbl *RiakURLSigKeyTable) fromGeneric(keys []URLSigKey) {
	tbl.Records = make([]RiakURLSigKeyRecord, len(keys))
	for i, key := range keys {
		tbl.Records[i] = RiakURLSigKeyRecord{
			Key:             key.URLSigKeys,
			DeliveryService: key.DeliveryService,
			CommonRecord:    key.CommonRecord,
		}
	}
}
func (tbl *RiakURLSigKeyTable) insertKeys(cluster *riak.Cluster) error {
	for _, record := range tbl.Records {
		objBytes, err := json.Marshal(record.Key)
		if err != nil {
			return err
		}

		err = setObject(cluster, makeRiakObject(objBytes, "url_sig_"+record.DeliveryService+".config"), BUCKET_URL_SIG)
		if err != nil {
			return err
		}
	}
	return nil
}
func (tbl *RiakURLSigKeyTable) validate() []string {
	errs := []string{}
	for i, record := range tbl.Records {
		if record.DeliveryService == "" {
			errs = append(errs, fmt.Sprintf("URL Key #%v: Delivery Service is blank!", i))
		}
	}
	return errs
}

type RiakURISignKeyRecord struct {
	Key             map[string]tc.URISignerKeyset
	DeliveryService string
	CommonRecord
}
type RiakURISignKeyTable struct {
	Records []RiakURISignKeyRecord
}

func (tbl *RiakURISignKeyTable) gatherKeys(cluster *riak.Cluster) error {
	tbl.Records = []RiakURISignKeyRecord{}
	objs, err := getObjects(cluster, BUCKET_URI_SIG)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		key := map[string]tc.URISignerKeyset{}
		if err := json.Unmarshal(obj.Value, &key); err != nil {
			return err
		}

		tbl.Records = append(tbl.Records, RiakURISignKeyRecord{
			DeliveryService: obj.Key,
			CommonRecord:    CommonRecord{},
			Key:             key,
		})
	}
	return nil
}
func (tbl *RiakURISignKeyTable) toGeneric() []URISignKey {
	keys := make([]URISignKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = URISignKey{
			DeliveryService: record.DeliveryService,
			Keys:            record.Key,
			CommonRecord:    record.CommonRecord,
		}
	}

	return keys
}
func (tbl *RiakURISignKeyTable) fromGeneric(keys []URISignKey) {
	tbl.Records = make([]RiakURISignKeyRecord, len(keys))

	for i, record := range keys {
		tbl.Records[i] = RiakURISignKeyRecord{
			Key:             record.Keys,
			DeliveryService: record.DeliveryService,
			CommonRecord:    record.CommonRecord,
		}
	}
}
func (tbl *RiakURISignKeyTable) insertKeys(cluster *riak.Cluster) error {
	for _, record := range tbl.Records {
		objBytes, err := json.Marshal(record.Key)
		if err != nil {
			return err
		}

		err = setObject(cluster, makeRiakObject(objBytes, record.DeliveryService), BUCKET_URI_SIG)
		if err != nil {
			return err
		}
	}
	return nil
}
func (tbl *RiakURISignKeyTable) validate() []string {
	errs := []string{}
	for i, record := range tbl.Records {
		if record.DeliveryService == "" {
			errs = append(errs, fmt.Sprintf("URI Key #%v: Delivery Service is blank!", i))
		}
	}
	return errs
}

// Riak functions
func makeRiakObject(data []byte, key string) *riak.Object {
	return &riak.Object{
		ContentType:     rfc.ApplicationJSON,
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Key:             key,
		Value:           data,
	}
}
func getObjects(cluster *riak.Cluster, bucket string) ([]*riak.Object, error) {
	objs := []*riak.Object{}
	keys, err := listKeys(cluster, bucket)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		objects, err := getObject(cluster, bucket, key)
		if err != nil {
			return nil, err
		}
		if len(objects) > 1 {
			return nil, errors.New(fmt.Sprintf("Unexpected number of objects %v, ignoring\n", len(objects)))
		}

		objs = append(objs, objects[0])
	}

	return objs, nil
}
func getObject(cluster *riak.Cluster, bucket string, key string) ([]*riak.Object, error) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		WithTimeout(time.Second * 60).
		Build()
	if err != nil {
		return nil, err
	}

	if err := cluster.Execute(cmd); err != nil {
		return nil, err
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response

	if rsp.IsNotFound {
		return nil, errors.New(fmt.Sprintf("ERROR Key not found: %v:%v", bucket, key))
	}

	return rsp.Values, nil
}
func setObject(cluster *riak.Cluster, obj *riak.Object, bucket string) error {
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(bucket).
		WithContent(obj).
		WithTimeout(time.Second * 5).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}
func search(cluster *riak.Cluster, index string, query string, filterQuery string, numRows uint32, fields []string) ([]*riak.SearchDoc, error) {
	var searchDocs []*riak.SearchDoc
	start := uint32(0)
	for i := uint32(0); ; i += 1 {
		riakCmd := riak.NewSearchCommandBuilder().
			WithQuery(query).
			WithNumRows(numRows).
			WithStart(start)
		if len(index) > 0 {
			riakCmd = riakCmd.WithIndexName(index)
		}
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
		if start == 0 {
			if cmd.Response == nil || cmd.Response.NumFound == 0 {
				return nil, nil
			}
			if cmd.Response.NumFound <= numRows {
				return cmd.Response.Docs, nil
			} else if numRows < cmd.Response.NumFound {
				searchDocs = make([]*riak.SearchDoc, cmd.Response.NumFound)
			}
			if cmd.Response.NumFound > 10000 {
				fmt.Printf("WARNING: found %v rows, press enter to continue", cmd.Response.NumFound)
				_, _ = fmt.Scanln()
			}
		}

		// If the total number of docs is not evenly divisible by 1000
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
		start += numRows
	}
}
func listKeys(cluster *riak.Cluster, bucket string) ([]string, error) {
	cmd, err := riak.NewListKeysCommandBuilder().
		WithBucket(bucket).
		WithTimeout(time.Second * 60).
		WithAllowListing().
		Build()
	if err != nil {
		return nil, err
	}

	if err := cluster.Execute(cmd); err != nil {
		return nil, err
	}

	lkc := cmd.(*riak.ListKeysCommand)
	rsp := lkc.Response
	return rsp.Keys, nil
}
func ping(cluster *riak.Cluster) error {
	ping := riak.PingCommandBuilder{}
	cmd, err := ping.Build()
	if err != nil {
		log.Fatal(err)
	}

	if err = cluster.Execute(cmd); err != nil {
		log.Fatal(err)
	}
	response, ok := cmd.(*riak.PingCommand)
	if !ok {
		log.Fatalf("Unexpected riak command type: %v", cmd)
	}

	if response.Error() != nil {
		return response.Error()
	}

	if !response.Success() {
		return errors.New("unable to ping riak")
	}

	return nil
}
func getRiakCluster(srv RiakConfig, authOptions *riak.AuthOptions) (*riak.Cluster, error) {
	if authOptions == nil {
		return nil, errors.New("ERROR: no riak auth information from riak.conf, cannot authenticate to any riak servers")
	}
	nodes := []*riak.Node{}
	nodeOpts := &riak.NodeOptions{
		RemoteAddress:       srv.Host + ":" + srv.Port,
		AuthOptions:         authOptions,
		HealthCheckInterval: time.Second * 5,
	}
	if nodeOpts.AuthOptions.TlsConfig != nil {
		nodeOpts.AuthOptions.TlsConfig.ServerName = srv.Host
	}
	node, err := riak.NewNode(nodeOpts)
	if err != nil {
		return nil, errors.New("creating riak node: " + err.Error())
	}
	nodes = append(nodes, node)
	if len(nodes) == 0 {
		return nil, errors.New("ERROR: no available riak servers")
	}
	opts := &riak.ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 2,
	}
	cluster, err := riak.NewCluster(opts)
	if err != nil {
		return nil, errors.New("creating riak cluster: " + err.Error())
	}
	return cluster, err
}
