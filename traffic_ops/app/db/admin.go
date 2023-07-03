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

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gopkg.in/yaml.v2"
)

type dbConfig struct {
	Development envConfig `yaml:"development"`
	Test        envConfig `yaml:"test"`
	Integration envConfig `yaml:"integration"`
	Production  envConfig `yaml:"production"`
}

type envConfig struct {
	Driver string `yaml:"driver"`
	Open   string `yaml:"open"`
}

func (conf dbConfig) getEnvironmentConfig(env string) (envConfig, error) {
	switch env {
	case envDevelopment:
		return conf.Development, nil
	case envTest:
		return conf.Test, nil
	case envIntegration:
		return conf.Integration, nil
	case envProduction:
		return conf.Production, nil
	default:
		return envConfig{}, errors.New("invalid environment: " + env)
	}
}

// the possible environments to use.
const (
	envDevelopment = "development"
	envTest        = "test"
	envIntegration = "integration"
	envProduction  = "production"
)

// keys in the database config's "open" string value.
const (
	hostKey     = "host"
	portKey     = "port"
	userKey     = "user"
	passwordKey = "password"
	dbNameKey   = "dbname"
	sslModeKey  = "sslmode"
)

// available commands.
const (
	cmdCreateDB        = "createdb"
	cmdDropDB          = "dropdb"
	cmdCreateMigration = "create_migration"
	cmdCreateUser      = "create_user"
	cmdDropUser        = "drop_user"
	cmdShowUsers       = "show_users"
	cmdReset           = "reset"
	cmdUpgrade         = "upgrade"
	cmdMigrate         = "migrate"
	cmdUp              = "up"
	cmdDown            = "down"
	cmdRedo            = "redo"
	// Deprecated: Migrate only tracks migration timestamp and dirty status, not a status for each migration.
	// Use CmdDBVersion to check the migration timestamp and dirty status.
	cmdStatus     = "status"
	cmdDBVersion  = "dbversion"
	cmdSeed       = "seed"
	cmdLoadSchema = "load_schema"
	cmdPatch      = "patch"
)

// Default file system paths for TODB files.
const (
	defaultDBDir            = "db/"
	defaultDBConfigPath     = defaultDBDir + "dbconf.yml"
	defaultDBMigrationsPath = defaultDBDir + "migrations"
	defaultDBSeedsPath      = defaultDBDir + "seeds.sql"
	defaultDBSchemaPath     = defaultDBDir + "create_tables.sql"
	defaultDBPatchesPath    = defaultDBDir + "patches.sql"
)

// Default file system paths for TV files.
const (
	defaultTrafficVaultDir            = defaultDBDir + "trafficvault/"
	defaultTrafficVaultDBConfigPath   = defaultTrafficVaultDir + "dbconf.yml"
	defaultTrafficVaultMigrationsPath = defaultTrafficVaultDir + "migrations"
	defaultTrafficVaultSchemaPath     = defaultTrafficVaultDir + "create_tables.sql"
)

// Default connection information.
const (
	defaultEnvironment = envDevelopment
	defaultDBSuperUser = "postgres"
)

const (
	// 2021012200000000_max_request_header_bytes_default_zero.sql.
	lastSquashedMigrationTimestamp uint = 2021012200000000
	// 2021012700000000_update_interfaces_multiple_routers.up.sql.
	firstMigrationTimestamp uint = 2021012700000000
)

// globals that are passed in via CLI flags and used in commands.
var (
	environment    string
	trafficVault   bool
	dbVersion      uint
	dbVersionDirty bool
)

// globals that are parsed out of DBConfigFile and used in commands.
var (
	connectionString string
	dbDriver         string
	dbName           string
	dbSuperUser      = defaultDBSuperUser
	dbUser           string
	dbPassword       string
	hostIP           string
	hostPort         string
	sslMode          string
	migrateInstance  *migrate.Migrate
	migrationName    string
)

// Actual TODB file paths.
var (
	dbConfigPath    = defaultDBConfigPath
	dbMigrationsDir = defaultDBMigrationsPath
	dbSeedsPath     = defaultDBSeedsPath
	dbSchemaPath    = defaultDBSchemaPath
	dbPatchesPath   = defaultDBPatchesPath
)

// Actual TV file paths.
var (
	trafficVaultDBConfigPath   = defaultTrafficVaultDBConfigPath
	trafficVaultMigrationsPath = defaultTrafficVaultMigrationsPath
	trafficVaultSchemaPath     = defaultTrafficVaultSchemaPath
)

func parseDBConfig() error {
	var cfgPath string
	if trafficVault {
		cfgPath = trafficVaultDBConfigPath
	} else {
		cfgPath = dbConfigPath
	}
	confBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("reading DB conf '%s': %w", cfgPath, err)
	}

	dbConfig := dbConfig{}
	err = yaml.Unmarshal(confBytes, &dbConfig)
	if err != nil {
		return errors.New("unmarshalling DB conf yaml: " + err.Error())
	}

	envConfig, err := dbConfig.getEnvironmentConfig(environment)
	if err != nil {
		return errors.New("getting environment config: " + err.Error())
	}

	dbDriver = envConfig.Driver
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
	stringPointers := []*string{&hostIP, &hostPort, &dbUser, &dbPassword, &dbName, &sslMode}
	keys := []string{hostKey, portKey, userKey, passwordKey, dbNameKey, sslModeKey}
	for index, stringPointer := range stringPointers {
		key := keys[index]
		*stringPointer, ok = open[key]
		if !ok {
			return errors.New("unable to get '" + key + "' for environment '" + environment + "'")
		}
	}

	return nil
}

func createDB() {
	dbExistsCmd := exec.Command("psql", "-h", hostIP, "-U", dbSuperUser, "-p", hostPort, "-tAc", "SELECT 1 FROM pg_database WHERE datname='"+dbName+"'")
	stderr := bytes.Buffer{}
	dbExistsCmd.Stderr = &stderr
	out, err := dbExistsCmd.Output()
	// An error is returned if the database could not be found, which is to be expected. Don't exit on this error.
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to check if DB already exists: "+err.Error()+", stderr: "+stderr.String())
	}
	if len(out) > 0 {
		fmt.Println("Database " + dbName + " already exists")
		return
	}
	createDBCmd := exec.Command("createdb", "-h", hostIP, "-p", hostPort, "-U", dbSuperUser, "-e", "--owner", dbUser, dbName)
	out, err = createDBCmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't create db " + dbName + ": " + err.Error())
	}
}

func dropDB() {
	fmt.Println("Dropping database: " + dbName)
	cmd := exec.Command("dropdb", "-h", hostIP, "-p", hostPort, "-U", dbSuperUser, "-e", "--if-exists", dbName)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't drop db " + dbName + ": " + err.Error())
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
	if err = os.MkdirAll(dbMigrationsDir, os.ModePerm); err != nil {
		die("Creating migrations directory " + dbMigrationsDir + ": " + err.Error())
	}
	migrationTime := time.Now()
	formattedMigrationTime := migrationTime.Format("20060102150405") + fmt.Sprintf("%02d", migrationTime.Nanosecond()%100)
	for _, direction := range []string{"up", "down"} {
		var migrationFile *os.File
		basename := fmt.Sprintf("%s_%s.%s.sql", formattedMigrationTime, migrationName, direction)
		filename := filepath.Join(dbMigrationsDir, basename)
		if migrationFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644); err != nil {
			die("Creating migration " + filename + ": " + err.Error())
		}
		defer log.Close(migrationFile, "closing migration "+filename)
		if _, err = migrationFile.Write([]byte(apacheLicense2)); err != nil {
			die("Writing content to migration " + filename + ": " + err.Error())
		}
		fmt.Println("Created migration ", filename)
	}
}

func createUser() {
	fmt.Println("Creating user: " + dbUser)
	userExistsCmd := exec.Command("psql", "-h", hostIP, "-U", dbSuperUser, "-p", hostPort, "-tAc", "SELECT 1 FROM pg_roles WHERE rolname='"+dbUser+"'")
	stderr := bytes.Buffer{}
	userExistsCmd.Stderr = &stderr
	out, err := userExistsCmd.Output()
	// An error is returned if the user could not be found, which is to be expected. Don't exit on this error.
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to check if user already exists: "+err.Error()+", stderr: "+stderr.String())
	}
	if len(out) > 0 {
		fmt.Println("User " + dbUser + " already exists")
		return
	}
	createUserCmd := exec.Command("psql", "-h", hostIP, "-p", hostPort, "-U", dbSuperUser, "-etAc", "CREATE USER "+dbUser+" WITH LOGIN ENCRYPTED PASSWORD '"+dbPassword+"'")
	out, err = createUserCmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't create user " + dbUser)
	}
}

func dropUser() {
	cmd := exec.Command("dropuser", "-h", hostIP, "-p", hostPort, "-U", dbSuperUser, "-i", "-e", dbUser)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't drop user " + dbUser)
	}
}

func showUsers() {
	cmd := exec.Command("psql", "-h", hostIP, "-p", hostPort, "-U", dbSuperUser, "-ec", `\du`)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
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
	if !trafficVault {
		patch()
	}
}

func maybeMigrateFromGoose() bool {
	var versionErr error
	dbVersion, dbVersionDirty, versionErr = migrateInstance.Version()
	if versionErr == nil {
		return false
	}
	if !errors.Is(versionErr, migrate.ErrNilVersion) {
		die("Error running migrate version: " + versionErr.Error())
	}
	if err := migrateInstance.Steps(1); err != nil {
		die("Error migrating to Migrate from Goose: " + err.Error())
	}
	dbVersion, dbVersionDirty, _ = migrateInstance.Version()
	return true
}

// runFirstMigration is essentially Migrate.Migrate(FirstMigrationTimestamp) but without the obligatory Migrate.versionExists() call.
// If calling Migrate.versionExists() is made optional, runFirstMigration() can be replaced.
func runFirstMigration() error {
	sourceDriver, sourceDriverErr := source.Open("file:" + dbMigrationsDir)
	if sourceDriverErr != nil {
		return fmt.Errorf("opening the migration source driver: " + sourceDriverErr.Error())
	}
	dbDriver, dbDriverErr := database.Open(connectionString)
	if dbDriverErr != nil {
		return fmt.Errorf("opening the dbdriver: " + dbDriverErr.Error())
	}
	firstMigration, firstMigrationName, migrationReadErr := sourceDriver.ReadUp(firstMigrationTimestamp)
	if migrationReadErr != nil {
		return fmt.Errorf("reading migration %s: %w", firstMigrationName, migrationReadErr)
	}
	if setDirtyVersionErr := dbDriver.SetVersion(int(firstMigrationTimestamp), true); setDirtyVersionErr != nil {
		return fmt.Errorf("setting the dirty version: %w", setDirtyVersionErr)
	}
	if migrateErr := dbDriver.Run(firstMigration); migrateErr != nil {
		return fmt.Errorf("running the migration: %w", migrateErr)
	}
	if setVersionErr := dbDriver.SetVersion(int(firstMigrationTimestamp), false); setVersionErr != nil {
		return fmt.Errorf("setting the version after successfully running the migration: %w", setVersionErr)
	}
	return nil
}

func runMigrations() {
	migratedFromGoose, migrationsShouldRun := initMigrate()
	if !migrationsShouldRun {
		return
	}
	if !trafficVault && dbVersion == lastSquashedMigrationTimestamp && !dbVersionDirty {
		if migrateErr := runFirstMigration(); migrateErr != nil {
			die(fmt.Sprintf("Error migrating from DB version %d to %d: %s", lastSquashedMigrationTimestamp, firstMigrationTimestamp, migrateErr.Error()))
		}
	}
	if upErr := migrateInstance.Up(); errors.Is(upErr, migrate.ErrNoChange) {
		if !migratedFromGoose {
			fmt.Fprintln(os.Stderr, upErr.Error())
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
	if err := migrateInstance.Steps(-1); err != nil {
		die("Error running migrate down: " + err.Error())
	}
}

func redo() {
	initMigrate()
	if downErr := migrateInstance.Steps(-1); downErr != nil {
		die("Error running migrate down 1 in 'redo': " + downErr.Error())
	}
	if upErr := migrateInstance.Steps(1); upErr != nil {
		die("Error running migrate up 1 in 'redo': " + upErr.Error())
	}
}

// Deprecated: Migrate does not track migration status of past migrations. Use dbversion() to check the migration timestamp and dirty status.
func status() {
	getDBVersion()
}

func getDBVersion() {
	initMigrate()
	fmt.Printf("dbversion %d", dbVersion)
	if dbVersionDirty {
		fmt.Printf(" (dirty)")
	}
	fmt.Println()
}

func seed() {
	if trafficVault {
		die("seed not supported for trafficvault environment")
	}
	fmt.Println("Seeding database with required data.")
	seedsBytes, err := ioutil.ReadFile(dbSeedsPath)
	if err != nil {
		die("unable to read '" + dbSeedsPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", hostIP, "-p", hostPort, "-d", dbName, "-U", dbUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(seedsBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPassword)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't patch database w/ required data")
	}
}

func loadSchema() {
	fmt.Println("Creating database tables.")
	schemaPath := dbSchemaPath
	if trafficVault {
		schemaPath = trafficVaultSchemaPath
	}
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		die("unable to read '" + schemaPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", hostIP, "-p", hostPort, "-d", dbName, "-U", dbUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(schemaBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPassword)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't create database tables")
	}
}

func patch() {
	if trafficVault {
		die("patch not supported for trafficvault environment")
	}
	fmt.Println("Patching database with required data fixes.")
	patchesBytes, err := ioutil.ReadFile(dbPatchesPath)
	if err != nil {
		die("unable to read '" + dbPatchesPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", hostIP, "-p", hostPort, "-d", dbName, "-U", dbUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(patchesBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPassword)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		die("Can't patch database with required data")
	}
}

func die(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func usage() string {
	programName := os.Args[0]
	var buff strings.Builder
	buff.WriteString("Usage: ")
	buff.WriteString(programName)
	buff.WriteString(` [OPTION] OPERATION [ARGUMENT(S)]

-c, --config CFG         Provide a path to a database configuration file,
                         instead of using the default (./db/dbconf.yml or
                         ./db/trafficvault/dbconf.yml for Traffic Vault)
-e, --env ENV            Use configuration for environment ENV (defined in
                         the database configuration file)
-h, --help               Show usage information and exit
-m, --migrations-dir DIR Use DIR as the migrations directory, instead of the
                         default (./db/migrations/ or
                         ./db/trafficvault/migrations for Traffic Vault)
-p, --patches PATCH      Provide a path to a set of database patch statements,
                         instead of using the default (./db/patches.sql)
-s, --schema SCHEMA      Provide a path to a schema file, instead of using the
                         default (./db/create_tables.sql or
                         ./db/trafficvault/create_tables.sql for Traffic Vault)
-S, --seeds SEEDS        Provide a path to a seeds statements file, instead of
                         using the default (./db/seeds.sql)
-v, --trafficvault       Perform operations for Traffic Vault instead of the
                         Traffic Ops database

OPERATION      The operation to perform; one of the following:

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
    dbversion   - Prints the current migration timestamp
    drop_user   - Execute 'drop_user' the user for the current environment
                  (traffic_ops).
    load_schema - Loads the database schema.
    patch       - Execute sql from db/patches.sql for loading post-migration data
                  patches (NOTE: not supported with --trafficvault option).
    redo        - Roll back the most recently applied migration, then run it again.
    reset       - Execute db 'dropdb', 'createdb', load_schema, migrate on the
                  database for the current environment.
    seed        - Execute sql from db/seeds.sql for loading static data (NOTE: not
                  supported with --trafficvault option). This MUST ONLY be done
				  after the schema is loaded and migrations have been run.
    show_users  - Execute sql to show all of the user for the current environment.
    status      - Prints the current migration timestamp (Deprecated, status is now an
                  alias for dbversion and will be removed in a future Traffic
                  Control release).
    upgrade     - Execute migrate and patch on the database for the current
                  environment.`,
	)
	return buff.String()
}

// collapses two options for 'name', using a default if given, stored into 'dest'.
// if the two option values conflict, the whole program dies and an error is printed to
// stderr.
func collapse(o1, o2, name, def string, dest *string) {
	if o1 == "" {
		if o2 == "" {
			*dest = def
			return
		}
		*dest = o2
		return
	}
	if o2 == "" {
		*dest = o1
		return
	}
	if o1 != o2 {
		die("conflicting definitions of '" + name + "' - must be specified only once\n" + usage())
	}
	*dest = o1
}

func main() {
	flag.Usage = func() { fmt.Fprintln(os.Stderr, usage()) }

	var shortCfg string
	var longCfg string
	flag.StringVar(&shortCfg, "c", "", "Provide a path to a database configuration file, instead of using the default (./db/dbconf.yml or ./db/trafficvault/dbconf.yml for Traffic Vault)")
	flag.StringVar(&longCfg, "config", "", "Provide a path to a database configuration file, instead of using the default (./db/dbconf.yml or ./db/trafficvault/dbconf.yml for Traffic Vault)")

	var shortEnv string
	var longEnv string
	flag.StringVar(&shortEnv, "e", "", "Use configuration for environment ENV (defined in the database configuration file)")
	flag.StringVar(&longEnv, "env", "", "Use configuration for environment ENV (defined in the database configuration file)")

	var shortMigrations string
	var longMigrations string
	flag.StringVar(&shortMigrations, "m", "", "Use DIR as the migrations directory, instead of the default (./db/migrations/ or ./db/trafficvault/migrations for Traffic Vault)")
	flag.StringVar(&longMigrations, "migrations-dir", "", "Use DIR as the migrations directory, instead of the default (./db/migrations/ or ./db/trafficvault/migrations for Traffic Vault)")

	var shortPatches string
	var longPatches string
	flag.StringVar(&shortPatches, "p", "", "Provide a path to a set of database patch statements, instead of using the default (./db/patches.sql)")
	flag.StringVar(&longPatches, "patches", "", "Provide a path to a set of database patch statements, instead of using the default (./db/patches.sql)")

	var shortSchema string
	var longSchema string
	flag.StringVar(&shortSchema, "s", "", "Provide a path to a schema file, instead of using the default (./db/create_tables.sql or ./db/trafficvault/create_tables.sql for Traffic Vault)")
	flag.StringVar(&longSchema, "schema", "", "Provide a path to a schema file, instead of using the default (./db/create_tables.sql or ./db/trafficvault/create_tables.sql for Traffic Vault)")

	var shortSeeds string
	var longSeeds string
	flag.StringVar(&shortSeeds, "S", "", "Provide a path to a seeds statements file, instead of using the default (./db/seeds.sql)")
	flag.StringVar(&longSeeds, "seeds", "", "Provide a path to a seeds statements file, instead of using the default (./db/seeds.sql)")

	flag.BoolVar(&trafficVault, "v", false, "Perform operations for Traffic Vault instead of the Traffic Ops database")
	flag.BoolVar(&trafficVault, "trafficvault", false, "Perform operations for Traffic Vault instead of the Traffic Ops database")
	flag.Parse()

	if trafficVault {
		collapse(shortCfg, longCfg, "config", defaultTrafficVaultDBConfigPath, &trafficVaultDBConfigPath)
		collapse(shortMigrations, longMigrations, "migrations-dir", defaultTrafficVaultMigrationsPath, &trafficVaultMigrationsPath)
		collapse(shortSchema, longSchema, "schema", defaultTrafficVaultSchemaPath, &trafficVaultSchemaPath)
	} else {
		collapse(shortCfg, longCfg, "config", defaultDBConfigPath, &dbConfigPath)
		collapse(shortMigrations, longMigrations, "migrations-dir", defaultDBMigrationsPath, &dbMigrationsDir)
		collapse(shortSchema, longSchema, "schema", defaultDBSchemaPath, &dbSchemaPath)
	}
	collapse(shortEnv, longEnv, "environment", defaultEnvironment, &environment)
	collapse(shortPatches, longPatches, "patches", defaultDBPatchesPath, &dbPatchesPath)
	collapse(shortSeeds, longSeeds, "seeds", defaultDBSeedsPath, &dbSeedsPath)

	if flag.Arg(0) == cmdCreateMigration {
		if len(flag.Args()) != 2 {
			die(usage())
		}
		migrationName = flag.Arg(1)
	} else if len(flag.Args()) != 1 || flag.Arg(0) == "" {
		die(usage())
	}
	if environment == "" {
		die(usage())
	}
	if err := parseDBConfig(); err != nil {
		die(err.Error())
	}
	commands := make(map[string]func())

	commands[cmdCreateDB] = createDB
	commands[cmdDropDB] = dropDB
	commands[cmdCreateMigration] = createMigration
	commands[cmdCreateUser] = createUser
	commands[cmdDropUser] = dropUser
	commands[cmdShowUsers] = showUsers
	commands[cmdReset] = reset
	commands[cmdUpgrade] = upgrade
	commands[cmdMigrate] = runMigrations
	commands[cmdUp] = runUp
	commands[cmdDown] = down
	commands[cmdRedo] = redo
	commands[cmdStatus] = status
	commands[cmdDBVersion] = getDBVersion
	commands[cmdSeed] = seed
	commands[cmdLoadSchema] = loadSchema
	commands[cmdPatch] = patch

	userCmd := flag.Arg(0)
	if cmd, ok := commands[userCmd]; ok {
		cmd()
	} else {
		die("invalid command: " + userCmd + "\n" + usage())
	}
}

// initMigrate initializes Migrate and returns whether a migration from Goose to
// Migrate was performed during said initialization and whether or not there are
// any actual migrations to run.
//
// If the second return value is false, then the global "Migrate" variable IS
// NOT SET so the caller must skip all migrations.
func initMigrate() (bool, bool) {
	var err error
	connectionString = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", dbDriver, dbUser, dbPassword, hostIP, hostPort, dbName, sslMode)
	if trafficVault {
		migrateInstance, err = migrate.New("file:"+trafficVaultMigrationsPath, connectionString)
	} else {
		migrateInstance, err = migrate.New("file:"+dbMigrationsDir, connectionString)
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, false
	}
	if err != nil {
		die("Starting Migrate: " + err.Error())
	}
	migrateInstance.Log = &Log{}
	return maybeMigrateFromGoose(), true
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
