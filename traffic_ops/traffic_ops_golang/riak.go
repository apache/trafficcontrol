package main

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/tostructs"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat/go-jwx/jwk"
	"io/ioutil"
	"net/http"
)

const RiakPort = 8087
const cdn_uri_keys_bucket = "cdn_uri_sig_keys" // riak namespace for cdn uri signing keys.

func getStringValue(resp *riak.FetchValueResponse) (string, error) {
	var obj *riak.Object

	if len(resp.Values) == 1 {
		obj = resp.Values[0]
	} else {
		return "", fmt.Errorf("no such object")
	}
	return string(obj.Value), nil
}

func assignDeliveryServiceUriKeysKeysHandler(db *sqlx.DB, cfg Config) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}
		defer r.Body.Close()

		xmlId := p["xml-id"]
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// validate that the received data is a valid jwk keyset
		var keys jwk.SymmetricKey
		if err := json.Unmarshal(data, &keys); err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}

		// create a storage object and store the data
		obj := &riak.Object{
			ContentType:     "text/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             xmlId,
			Value:           []byte(data),
		}

		err = saveObject(obj, cdn_uri_keys_bucket, db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", data)
	}
}

// saves an object to riak storage
func saveObject(obj *riak.Object, bucket string, db *sqlx.DB, cfg Config) error {
	// create and start a cluster
	cluster, err := getRiakCluster(db, cfg, 12)
	if err != nil {
		return err
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			log.Errorf("%v\n", err)
		}
	}()

	if err = cluster.Start(); err != nil {
		return err
	}

	// build store command and execute.
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(bucket).
		WithContent(obj).
		Build()
	if err != nil {
		return err
	}
	if err := cluster.Execute(cmd); err != nil {
		return err
	}

	return nil
}

// fetch an object from riak storage
func fetchObject(key string, bucket string, db *sqlx.DB, cfg Config) (*riak.FetchValueCommand, error) {
	// build the fetch command
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		return nil, err
	}
	// create and start a riak cluster
	cluster, err := getRiakCluster(db, cfg, 12)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cluster.Stop(); err != nil {
			log.Errorf("%v\n", err)
		}
	}()
	if err = cluster.Start(); err != nil {
		return nil, err
	}
	if err = cluster.Execute(cmd); err != nil {
		return nil, err
	}
	fvc := cmd.(*riak.FetchValueCommand)

	return fvc, err
}

// endpoint handler for fetching uri signing keys from riak
func urisignkeysHandler(db *sqlx.DB, cfg Config) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		xmlId := p["xml-id"]

		fvc, err := fetchObject(xmlId, cdn_uri_keys_bucket, db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		resp, err := getStringValue(fvc.Response)
		if err != nil {
			handleErr(err, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", resp)
	}
}

// returns a riak cluster of online riak nodes.
func getRiakCluster(db *sqlx.DB, cfg Config, maxNodes int) (*riak.Cluster, error) {
	riakServerQuery := fmt.Sprintf(`
		SELECT s.host_name, s.domain_name FROM server s 
		INNER JOIN type t on s.type = t.id 
		INNER JOIN status st on s.status = st.id 
		WHERE t.name = 'RIAK' AND st.name = 'ONLINE'
		`)

	var nodes []*riak.Node
	rows, err := db.Query(riakServerQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s tostructs.Server
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

	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}
	cluster, err := riak.NewCluster(opts)

	return cluster, err
}
