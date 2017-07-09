// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

const DefaultConfigPath = "/etc/goto/config.json"

func main() {
	configFileName := flag.String("cfg", "", "The config file path")
	flag.Parse()
	if *configFileName == "" {
		*configFileName = DefaultConfigPath
	}

	cfg, err := LoadConfig(*configFileName)
	if err != nil {
		fmt.Println("Error loading config '" + *configFileName + "': " + err.Error())
		return
	}

	sslStr := "require"
	if !cfg.DBSSL {
		sslStr = "disable"
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.DBUser, cfg.DBPass, cfg.DBServer, cfg.DBDB, sslStr))
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	RegisterRoutes(ServerData{DB: db, Config: cfg})
	fmt.Println("Listening on " + cfg.HTTPPort)
	if err := http.ListenAndServeTLS(":"+cfg.HTTPPort, cfg.CertPath, cfg.KeyPath, nil); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
		return
	}
}
