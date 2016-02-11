InfluxDB_Tools: 

These tools are meant to be used to help a user quickly get new databases and continuous queries setup in influxdb for traffic stats.  
They are specific for traffic stats and are not meant to be generic to influxdb.
For more information see: http://traffic-control-cdn.net/docs/latest/admin/traffic_stats.html#influxdb-tools

Pre-Requisites: 
1. Go 1.5.x or later
2. Influxdb
3. configured $GOPATH (e.g. export GOPATH=~/go)

Using create_ts_databases.go
1. Install InfluxDb Client
	- go get github.com/influxdata/influxdb
	- cd $GOPATH/src/github.com/influxdata/influxdb
	- git checkout v0.9.6.1 (or whatever version of influxdb you are running)
	- go install

2. Build it
	- go build create_ts_databases.go

3. Run it 
	- ./create_ts_databases
	- optional flags:
		- influxUrl -  The influxdb url and port
		- replication -  The number of nodes in the cluster
	- example: ./create_ts_databases -influxUrl=localhost:8086 -replication=3


Using sync_ts_databases.go
1. Install InfluxDb Client (0.9.6.1 version)
	- go get github.com/influxdata/influxdb
	- cd $GOPATH/src/github.com/influxdata/influxdb
	- git checkout v0.9.6.1
	- go install

2. Build it
	- go build sync_ts_databases.go

3. Run it 
	- required flags:
		- sourceUrl - The URL of the source database 
		- targetUrl - The URL of the target database
	-optional flags:
		- database - The database to sync (default = sync all databases)
		- days - Days in the past to sync (default = sync all data)
	- example: ./sync_ts_databases -sourceUrl=http://influxdb-production-01.kabletown.net:8086 -targetUrl=http://influxdb-dev-01.kabletown.net:8086 -database=cache_stats -days=7

