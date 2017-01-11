package influxdb

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

// QueryDB takes an influx client interface, cmd, and db strings. It tries to
// execute the query on the given client, and returns a slice of influx.Result,
// and an error, if any were present
func QueryDB(client influx.Client, cmd, db string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: db,
	}
	response, err := client.Query(q)
	if err != nil {
		return res, errors.Wrapf(err, "failed to execute cmd: %s on db %s", cmd, db)
	}
	if response.Error() != nil {
		return res, errors.Wrapf(response.Error(), "got error in query response for cmd: %s on db %s", cmd, db)
	}
	res = response.Results

	return res, err
}

// Create takes an influx.Client, and a create cmd in order to create objects
// in the influx db (it doesn't require a db argument like QueryDB)
func Create(client influx.Client, cmd string) (err error) {
	_, err = QueryDB(client, cmd, "")
	return err
}
