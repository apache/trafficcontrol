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
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Development EnvConfig `yaml:"development"`
	Test        EnvConfig `yaml:"test"`
	Integration EnvConfig `yaml:"integration"`
	Production  EnvConfig `yaml:"production"`
}

type EnvConfig struct {
	Driver string `yaml:"driver"`
	Open   string `yaml:"open"`
}

func (conf DBConfig) getEnvironmentConfig(env string) (EnvConfig, error) {
	switch env {
	case EnvDevelopment:
		return conf.Development, nil
	case EnvTest:
		return conf.Test, nil
	case EnvIntegration:
		return conf.Integration, nil
	case EnvProduction:
		return conf.Production, nil
	default:
		return EnvConfig{}, errors.New("invalid environment: " + env)
	}
}

const (
	// the possible environments to use
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvIntegration = "integration"
	EnvProduction  = "production"

	// keys in the database config's "open" string value
	HostKey     = "host"
	PortKey     = "port"
	UserKey     = "user"
	PasswordKey = "password"
	DBNameKey   = "dbname"
	SSLModeKey  = "sslmode"

	// available commands
	CmdCreateDB        = "createdb"
	CmdDropDB          = "dropdb"
	CmdCreateMigration = "create_migration"
	CmdCreateUser      = "create_user"
	CmdDropUser        = "drop_user"
	CmdShowUsers       = "show_users"
	CmdReset           = "reset"
	CmdUpgrade         = "upgrade"
	CmdMigrate         = "migrate"
	CmdUp              = "up"
	CmdDown            = "down"
	CmdRedo            = "redo"
	// Deprecated: Migrate only tracks migration version and dirty status, not a status for each migration.
	// Use CmdDBVersion to check the migration version and dirty status.
	CmdStatus     = "status"
	CmdDBVersion  = "dbversion"
	CmdSeed       = "seed"
	CmdLoadSchema = "load_schema"
	CmdPatch      = "patch"

	dbDir              = "db/"
	DBConfigPath       = dbDir + "dbconf.yml"
	DBMigrationsPath   = dbDir + "migrations"
	DBMigrationsSource = "file:" + DBMigrationsPath
	DBSeedsPath        = dbDir + "seeds.sql"
	DBSchemaPath       = dbDir + "create_tables.sql"
	DBPatchesPath      = dbDir + "patches.sql"
	DefaultEnvironment = EnvDevelopment
	DefaultDBSuperUser = "postgres"

	TrafficVaultDBConfigPath     = TrafficVaultDir + "dbconf.yml"
	TrafficVaultMigrationsSource = "file:" + TrafficVaultDir + "migrations"
	TrafficVaultDir              = dbDir + "trafficvault/"
	TrafficVaultSchemaPath       = TrafficVaultDir + "create_tables.sql"
)

var (
	// globals that are passed in via CLI flags and used in commands
	Environment    string
	TrafficVault   bool
	DBVersion      uint
	DBVersionDirty bool

	// globals that are parsed out of DBConfigFile and used in commands
	DBDriver      string
	DBName        string
	DBSuperUser   = DefaultDBSuperUser
	DBUser        string
	DBPassword    string
	HostIP        string
	HostPort      string
	SSLMode       string
	Migrate       *migrate.Migrate
	MigrationName string
)

func parseDBConfig() error {
	dbConfigPath := DBConfigPath
	if TrafficVault {
		dbConfigPath = TrafficVaultDBConfigPath
	}
	confBytes, err := ioutil.ReadFile(dbConfigPath)
	if err != nil {
		return errors.New("reading DB conf '" + dbConfigPath + "': " + err.Error())
	}

	dbConfig := DBConfig{}
	err = yaml.Unmarshal(confBytes, &dbConfig)
	if err != nil {
		return errors.New("unmarshalling DB conf yaml: " + err.Error())
	}

	envConfig, err := dbConfig.getEnvironmentConfig(Environment)
	if err != nil {
		return errors.New("getting environment config: " + err.Error())
	}

	DBDriver = envConfig.Driver
	// parse the 'open' string into a map
	open := make(map[string]string)
	pairs := strings.Split(envConfig.Open, " ")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.Split(pair, "=")
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			continue
		}
		open[kv[0]] = kv[1]
	}

	ok := false
	stringPointers := []*string{&HostIP, &HostPort, &DBUser, &DBPassword, &DBName, &SSLMode}
	keys := []string{HostKey, PortKey, UserKey, PasswordKey, DBNameKey, SSLModeKey}
	for index, stringPointer := range stringPointers {
		key := keys[index]
		*stringPointer, ok = open[key]
		if !ok {
			return errors.New("unable to get '" + key + "' for environment '" + Environment + "'")
		}
	}

	return nil
}

func createDB() {
	dbExistsCmd := exec.Command("psql", "-h", HostIP, "-U", DBSuperUser, "-p", HostPort, "-tAc", "SELECT 1 FROM pg_database WHERE datname='"+DBName+"'")
	stderr := bytes.Buffer{}
	dbExistsCmd.Stderr = &stderr
	out, err := dbExistsCmd.Output()
	// An error is returned if the database could not be found, which is to be expected. Don't exit on this error.
	if err != nil {
		fmt.Println("unable to check if DB already exists: " + err.Error() + ", stderr: " + stderr.String())
	}
	if len(out) > 0 {
		fmt.Println("Database " + DBName + " already exists")
		return
	}
	createDBCmd := exec.Command("createdb", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-e", "--owner", DBUser, DBName)
	out, err = createDBCmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't create db " + DBName)
	}
}

func dropDB() {
	fmt.Println("Dropping database: " + DBName)
	cmd := exec.Command("dropdb", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-e", "--if-exists", DBName)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't drop db " + DBName)
	}
}

func createMigration() {
	const apacheLicense2 = `/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership.  The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

`
	var err error
	if err = os.MkdirAll(DBMigrationsPath, os.ModePerm); err != nil {
		die("Creating migrations directory " + DBMigrationsPath + ": " + err.Error())
	}
	migrationTime := time.Now()
	formattedMigrationTime := migrationTime.Format("20060102150405") + fmt.Sprintf("%02d", migrationTime.Nanosecond()%100)
	for _, direction := range []string{"up", "down"} {
		var migrationFile *os.File
		basename := fmt.Sprintf("%s_%s.%s.sql", formattedMigrationTime, MigrationName, direction)
		filename := filepath.Join(DBMigrationsPath, basename)
		if migrationFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644); err != nil {
			die("Creating migration " + filename + ": " + err.Error())
		}
		defer log.Close(migrationFile, "Closing migration "+filename)
		if migrationFile.Write([]byte(apacheLicense2)); err != nil {
			die("Writing content to migration " + filename + ": " + err.Error())
		}
		fmt.Printf("Created migration %s\n", filename)
	}
}

func createUser() {
	fmt.Println("Creating user: " + DBUser)
	userExistsCmd := exec.Command("psql", "-h", HostIP, "-U", DBSuperUser, "-p", HostPort, "-tAc", "SELECT 1 FROM pg_roles WHERE rolname='"+DBUser+"'")
	stderr := bytes.Buffer{}
	userExistsCmd.Stderr = &stderr
	out, err := userExistsCmd.Output()
	// An error is returned if the user could not be found, which is to be expected. Don't exit on this error.
	if err != nil {
		fmt.Println("unable to check if user already exists: " + err.Error() + ", stderr: " + stderr.String())
	}
	if len(out) > 0 {
		fmt.Println("User " + DBUser + " already exists")
		return
	}
	createUserCmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-etAc", "CREATE USER "+DBUser+" WITH LOGIN ENCRYPTED PASSWORD '"+DBPassword+"'")
	out, err = createUserCmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't create user " + DBUser)
	}
}

func dropUser() {
	cmd := exec.Command("dropuser", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-i", "-e", DBUser)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't drop user " + DBUser)
	}
}

func showUsers() {
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-ec", `\du`)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't show users")
	}
}

func reset() {
	createUser()
	dropDB()
	createDB()
	loadSchema()
	runMigrations()
}

func upgrade() {
	runMigrations()
	if !TrafficVault {
		seed()
		patch()
	}
}

func maybeMigrateFromGoose() bool {
	var versionErr error
	DBVersion, DBVersionDirty, versionErr = Migrate.Version()
	if versionErr == nil {
		return false
	}
	if versionErr != migrate.ErrNilVersion {
		die("Error running migrate version: " + versionErr.Error())
	}
	if err := Migrate.Steps(1); err != nil {
		die("Error migrating to Migrate from Goose: " + err.Error())
	}
	return true
}

func runMigrations() {
	migratedFromGoose := initMigrate()
	upErr := Migrate.Up()
	if upErr == migrate.ErrNoChange {
		if !migratedFromGoose {
			println(upErr.Error())
		}
	} else if upErr != nil {
		die("Error running migrate up: " + upErr.Error())
	}

}

func runUp() {
	runMigrations()
}

func down() {
	initMigrate()
	if err := Migrate.Steps(-1); err != nil {
		die("Error running migrate down: " + err.Error())
	}
}

func redo() {
	initMigrate()
	if downErr := Migrate.Steps(-1); downErr != nil {
		die("Error running migrate down 1 in 'redo': " + downErr.Error())
	}
	if upErr := Migrate.Steps(1); upErr != nil {
		die("Error running migrate up 1 in 'redo': " + upErr.Error())
	}
}

// Deprecated: Migrate does not track migration status of past migrations. Use dbversion() to check the migration version and dirty status.
func status() {
	dbVersion()
}

func dbVersion() {
	initMigrate()
	fmt.Printf("dbversion %d", DBVersion)
	if DBVersionDirty {
		fmt.Printf(" (dirty)")
	}
	println()
}

func seed() {
	if TrafficVault {
		die("seed not supported for trafficvault environment")
	}
	fmt.Println("Seeding database w/ required data.")
	seedsBytes, err := ioutil.ReadFile(DBSeedsPath)
	if err != nil {
		die("unable to read '" + DBSeedsPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-d", DBName, "-U", DBUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(seedsBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+DBPassword)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't patch database w/ required data")
	}
}

func loadSchema() {
	fmt.Println("Creating database tables.")
	schemaPath := DBSchemaPath
	if TrafficVault {
		schemaPath = TrafficVaultSchemaPath
	}
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		die("unable to read '" + DBSchemaPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-d", DBName, "-U", DBUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(schemaBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+DBPassword)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't create database tables")
	}
}

func patch() {
	if TrafficVault {
		die("patch not supported for trafficvault environment")
	}
	fmt.Println("Patching database with required data fixes.")
	patchesBytes, err := ioutil.ReadFile(DBPatchesPath)
	if err != nil {
		die("unable to read '" + DBPatchesPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-d", DBName, "-U", DBUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(patchesBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+DBPassword)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't patch database w/ required data")
	}
}

func die(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func usage() string {
	programName := os.Args[0]
	home := "$HOME"
	home = os.Getenv("HOME")
	return `
Usage:  ` + programName + ` [--trafficvault] [--env (development|test|production|integration)] [arguments]

Example:  ` + programName + ` --env=test reset

Purpose:  This script is used to manage the Traffic Ops database and Traffic Vault PostgreSQL backend database.
          The Traffic Ops environments and database names are defined in the dbconf.yml, and for Traffic Vault
          they are defined in trafficvault/dbconf.yml. In order to execute commands against the Traffic Vault
          database, the the --trafficvault option.

NOTE:
Postgres Superuser: The 'postgres' superuser needs to be created to run ` + programName + ` and setup databases.
If the 'postgres' superuser has not been created or password has not been set then run the following commands accordingly.

Create the 'postgres' user as a super user (if not created):

     $ createuser postgres --superuser --createrole --createdb --login --pwprompt

Modify your ` + home + `/.pgpass file which allows for easy command line access by defaulting the user and password for the database
without prompts.

 Postgres .pgpass file format:
 hostname:port:database:username:password

 ----------------------
 Example Contents
 ----------------------
 *:*:*:postgres:your-postgres-password
 *:*:*:traffic_ops:the-password-in-dbconf.yml
 *:*:*:traffic_vault:the-password-in-trafficvault-dbconf.yml
 ----------------------

 Save the following example into this file ` + home + `/.pgpass with the permissions of this file
 so only your user can read and write.

     $ chmod 0600 ` + home + `/.pgpass

===================================================================================================================
` + programName + ` arguments:

migrate     - Execute migrate (without seeds or patches) on the database for the
              current environment.
up          - Alias for 'migrate'
down        - Roll back a single migration from the current version.
createdb    - Execute db 'createdb' the database for the current environment.
dropdb      - Execute db 'dropdb' on the database for the current environment.
create_migration NAME
            - Creates a pair of timestamped up/down migrations titled NAME.
create_user - Execute 'create_user' the user for the current environment
              (traffic_ops).
dbversion   - Prints the current migration version
drop_user   - Execute 'drop_user' the user for the current environment
              (traffic_ops).
patch       - Execute sql from db/patches.sql for loading post-migration data
              patches (NOTE: not supported with --trafficvault option).
redo        - Roll back the most recently applied migration, then run it again.
reset       - Execute db 'dropdb', 'createdb', load_schema, migrate on the
              database for the current environment.
seed        - Execute sql from db/seeds.sql for loading static data (NOTE: not
              supported with --trafficvault option).
show_users  - Execute sql to show all of the user for the current environment.
status      - Prints the current migration version (Deprecated, status is now an
              alias for dbversion and will be removed in a future Traffic
              Control release).
upgrade     - Execute migrate, seed, and patches on the database for the current
              environment.
`
}

func main() {
	flag.StringVar(&Environment, "env", DefaultEnvironment, "The environment to use (defined in "+DBConfigPath+").")
	flag.BoolVar(&TrafficVault, "trafficvault", false, "Run this for the Traffic Vault database")
	flag.Parse()
	if flag.Arg(0) == CmdCreateMigration {
		if len(flag.Args()) != 2 {
			die(usage())
		}
		MigrationName = flag.Arg(1)
	} else if len(flag.Args()) != 1 || flag.Arg(0) == "" {
		die(usage())
	}
	if Environment == "" {
		die(usage())
	}
	if err := parseDBConfig(); err != nil {
		die(err.Error())
	}
	commands := make(map[string]func())

	commands[CmdCreateDB] = createDB
	commands[CmdDropDB] = dropDB
	commands[CmdCreateMigration] = createMigration
	commands[CmdCreateUser] = createUser
	commands[CmdDropUser] = dropUser
	commands[CmdShowUsers] = showUsers
	commands[CmdReset] = reset
	commands[CmdUpgrade] = upgrade
	commands[CmdMigrate] = runMigrations
	commands[CmdUp] = runUp
	commands[CmdDown] = down
	commands[CmdRedo] = redo
	commands[CmdStatus] = status
	commands[CmdDBVersion] = dbVersion
	commands[CmdSeed] = seed
	commands[CmdLoadSchema] = loadSchema
	commands[CmdPatch] = patch

	userCmd := flag.Arg(0)
	if cmd, ok := commands[userCmd]; ok {
		cmd()
	} else {
		fmt.Println(usage())
		die("invalid command: " + userCmd)
	}
}

func initMigrate() bool {
	var err error
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", DBDriver, DBUser, DBPassword, HostIP, HostPort, DBName, SSLMode)
	if TrafficVault {
		Migrate, err = migrate.New(TrafficVaultMigrationsSource, connectionString)
	} else {
		Migrate, err = migrate.New(DBMigrationsSource, connectionString)
	}
	if err != nil {
		die("Starting Migrate: " + err.Error())
	}
	Migrate.Log = &Log{}
	return maybeMigrateFromGoose()
}

// Log represents the logger
// https://github.com/golang-migrate/migrate/blob/v4.14.1/internal/cli/log.go#L9-L12
type Log struct {
	verbose bool
}

// Printf prints out formatted string into a log
// https://github.com/golang-migrate/migrate/blob/v4.14.1/internal/cli/log.go#L14-L21
func (l *Log) Printf(format string, v ...interface{}) {
	if l.verbose {
		fmt.Printf(format, v...)
	} else {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}

// Println prints out args into a log
// https://github.com/golang-migrate/migrate/blob/v4.14.1/internal/cli/log.go#L23-L30
func (l *Log) Println(args ...interface{}) {
	if l.verbose {
		fmt.Println(args...)
	} else {
		fmt.Fprintln(os.Stderr, args...)
	}
}

// Verbose shows if verbose print enabled
// https://github.com/golang-migrate/migrate/blob/v4.14.1/internal/cli/log.go#L32-L35
func (l *Log) Verbose() bool {
	return l.verbose
}
