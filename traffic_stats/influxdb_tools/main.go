/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "influx-tools"
	app.Version = "0.1.0"
	app.Usage = "influx-tools provides cli methods for creating and syncing the requisite influxdb databases"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "create",
			Usage:  "create the influxDB tables",
			Action: create,
			Flags:  createFlags(),
		},
		cli.Command{
			Name:   "sync",
			Usage:  "sync the influxDB tables",
			Action: sync,
			Flags:  syncFlags(),
		},
	}
	app.Run(os.Args)
}
