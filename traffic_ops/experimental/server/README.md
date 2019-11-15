<!--
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
-->

##### Intro
This is a prototype for Traffic Ops 2.0 server.  See also https://github.com/Comcast/traffic_control/wiki/Traffic_Ops_20

##### One time generation of CRUD files 
  ```
  [jvd@laika tools (master=)]$ go run gen_goto2.go postgres to_user **** to_development localhost 5432
  [goose_db_version federation_deliveryservice hwinfo job_result log deliveryservice_regex deliveryservice_server deliveryservice_tmuser federation_federation_resolver federation_resolver federation_tmuser federation job job_agent job_status deliveryservice parameter division cdn region profile_parameter role servercheck status stats_summary to_extension tm_user type asn cachegroup_parameter cachegroup profile regex phys_location server staticdnsentry]
  goose_db_version: Ok 5032
  federation_deliveryservice: Ok 5312
  hwinfo: Ok 4588
  job_result: Ok 4854
  log: Ok 4630
  deliveryservice_regex: Ok 5152
  deliveryservice_server: Ok 5163
  deliveryservice_tmuser: Ok 5187
  federation_federation_resolver: Ok 5470
  federation_resolver: Ok 4935
  federation_tmuser: Ok 5008
  federation: Ok 4759
  job: Ok 5884
  job_agent: Ok 4715
  job_status: Ok 4635
  deliveryservice: Ok 11287
  parameter: Ok 4690
  division: Ok 4385
  cdn: Ok 4402
  region: Ok 4440
  profile_parameter: Ok 4877
  role: Ok 4421
  servercheck: Ok 7460
  status: Ok 4498
  stats_summary: Ok 5252
  to_extension: Ok 6064
  tm_user: Ok 7174
  type: Ok 4583
  asn: Ok 4329
  cachegroup_parameter: Ok 5036
  cachegroup: Ok 5528
  profile: Ok 4537
  regex: Ok 4394
  phys_location: Ok 5776
  server: Ok 8878
  staticdnsentry: Ok 5312
  [jvd@laika tools (master *=)]$
  ```

##### Using swagger and the go swagger tools
We're using https://github.com/yvasiyarov/swagger to generate the swagger testing files. To get the swagger pages up:
Note: for now, we are using the web.go method to get the swagger pages up, later we'll move that to a hosted index.json. To start, do:

1. Get the latest swagger ui, and mod it to point to the right place, and support jwt:
  ```
  [jvd@laika swagger-ui (master *=)]$ git diff
  diff --git a/dist/index.html b/dist/index.html
  index 4531c44..d6d4399 100644
  --- a/dist/index.html
  +++ b/dist/index.html
  @@ -35,7 +35,7 @@
         if (url && url.length > 1) {
           url = decodeURIComponent(url[1]);
         } else {
  -        url = "http://petstore.swagger.io/v2/swagger.json";
  +        url = "http://localhost:8081/api/2.0";
         }

         // Pre load translate...
  @@ -81,9 +81,12 @@
         function addApiKeyAuthorization(){
           var key = encodeURIComponent($('#input_apiKey')[0].value);
           if(key && key.trim() != "") {
  -            var apiKeyAuth = new SwaggerClient.ApiKeyAuthorization("api_key", key, "query");
  -            window.swaggerUi.api.clientAuthorizations.add("api_key", apiKeyAuth);
  -            log("added key " + key);
  +            //var apiKeyAuth = new SwaggerClient.ApiKeyAuthorization("api_key", key, "query");
  +            //window.swaggerUi.api.clientAuthorizations.add("api_key", apiKeyAuth);
  +            var apiKeyAuth = new SwaggerClient.ApiKeyAuthorization( "Authorization", "Bearer " + key, "header" );
  +            window.swaggerUi.api.clientAuthorizations.add( "bearer", apiKeyAuth );
  +            log( "Set bearer token: " + key );
  +            //log("added key " + key);
           }
         }

  [jvd@laika swagger-ui (master *=)]$ git remote -v
  origin  https://github.com/swagger-api/swagger-ui.git (fetch)
  origin  https://github.com/swagger-api/swagger-ui.git (push)
  [jvd@laika swagger-ui (master *=)]$
  ```

2. Generate the docs.go file:
  ```
  [jvd@laika swagger-api]$ pwd
  /Users/jvd/work/gh/swagger-api
  [jvd@laika swagger-api]$ swagger -apiPackage github.com/apache/trafficcontrol/traffic_ops/experimental/server/api -mainApiFile github.com/apache/trafficcontrol/traffic_ops/experimental/server/api/api.go -format go
  2016/01/16 09:57:34 Start parsing
  2016/01/16 09:57:36 Finish parsing
  2016/01/16 09:57:36 Doc file generated
  [jvd@laika swagger-api]$ ls -l *.go
  -rw-r--r--  1 jvd  staff  343786 Jan 16 09:57 docs.go
  -rw-r--r--  1 jvd  staff    2212 Jan 14 08:59 web.go
  ```

3. Run the web.go app:
  ```
   [jvd@laika swagger-api]$ pwd
  /Users/jvd/work/gh/swagger-api
  [jvd@laika swagger-api]$ ls -l *.go
  -rw-r--r--  1 jvd  staff  343786 Jan 16 09:57 docs.go
  -rw-r--r--  1 jvd  staff    2212 Jan 14 08:59 web.go
  [jvd@laika swagger-api]$ go run web.go docs.go -port 8081 -api http://localhost:8080 -staticPath ./swagger-ui/dist/
  Hi
  Handle!! /api/2.0/ 0x2200
  2016/01/16 09:57:47 7807 0.0.0.0:8081
  ```

4. Get a token (make sure you have the go api app running on port 8080):
   ```
   [jvd@laika ~]$ curl --header "Content-Type:application/json" -XPOST http://localhost:8080/login -d'{"u":"jvd", "p":"******"}'
  {"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTMwNDU3NjMsInJvbGUiOjYsInVzZXJpZCI6OTB9.YKWiI3_rWyy3iD7xLJ2VOunb7xnKNNakTSt8KoZ5S1k"}[jvd@laika ~]$
  [jvd@laika ~]$
   ```

5. Use this token string in the swagger UI that you get to by pointing your browser at http://localhost:8081

Note that there are still some conversions that need to be made manually to have everything work (will fix this later), null.*->*, int64->int, float64->float.

##### Converting your existing database from MySQL to PostgreSQL
We're using PostgreSQL as a database. 
0. Prepare your tools, patch FromMySqlToPostgreSql. On a MAC:
   ```
   [jvd@laika FromMySqlToPostgreSql (master *%=)] brew install php56
   [jvd@laika FromMySqlToPostgreSql (master *%=)] PATH="/usr/local/sbin:$PATH"
   [jvd@laika FromMySqlToPostgreSql (master *%=)]$ pwd
   /Users/jvd/work/gh/AntolyUss/FromMySqlToPostgreSql
   [jvd@laika FromMySqlToPostgreSql (master *%=)]$ git diff
   diff --git a/migration/FromMySqlToPostgreSql/FromMySqlToPostgreSql.php b/migration/FromMySqlToPostgreSql/FromMySqlToPostgreSql.php
   index fee695d..a7063b7 100755
   --- a/migration/FromMySqlToPostgreSql/FromMySqlToPostgreSql.php
   +++ b/migration/FromMySqlToPostgreSql/FromMySqlToPostgreSql.php
   @@ -342,7 +342,7 @@ class FromMySqlToPostgreSql
                   . "\t-------------------------------------------------------"
                   . PHP_EOL . PHP_EOL;

   -        $this->log($strError, true);
   +        $this->log($strError . PHP_EOL);

         if (!empty($this->strWriteErrorLogTo)) {
             if (is_resource($this->resourceErrorLog)) {
   ```  

1. Set up pg environment

  ```
  to_integration=# \q
  jvd@pixel:~/work/gh/AnatolyUss/FromMySqlToPostgreSql$ psql --user postgres postgres
  psql (9.4.5)
  Type "help" for help.

  postgres=# create user to_user password '*******';
  CREATE ROLE
  postgres=# create schema to_development;
  CREATE SCHEMA
  postgres=# grant all on all tables in schema to_development to to_user;
  GRANT
  postgres=#
  ```
2. Run FromMySqlToPostgreSql. See https://github.com/AnatolyUss/FromMySqlToPostgreSql for prereqs. The config file that worked for me:

  ```
  {
      "source_description" : [
          "Connection string to your MySql database",
          "Please ensure, that you have defined your connection string properly.",
          "Ensure, that details like 'charset=UTF8' are included in your connection string (if necessary)."
      ],
      "source" : "mysql:host=localhost;port=3306;charset=UTF8;dbname=to_development,root,your_mysql_passwd",
      
      "target_description" : [
          "Connection string to your PostgreSql database",
          "Please ensure, that you have defined your connection string properly.",
          "Ensure, that details like options='[double dash]client_encoding=UTF8' are included in your connection string (if necessary)."
      ],
      "target" : "pgsql:host=localhost;port=5432;dbname=to_development;options=--client_encoding=UTF8,postgres,your_postgres_passwd",
      
      "encoding_description" : [
          "PHP encoding type.",
          "If not supplied, then UTF-8 will be used as a default."
      ],
      "encoding" : "UTF-8",
      
      "schema_description" : [
          "schema - a name of the schema, that will contain all migrated tables.",
          "If not supplied, then a new schema will be created automatically."
      ],
      "schema" : "public",
      
      "data_chunk_size_description" : [
          "During migration each table's data will be split into chunks of data_chunk_size (in MB).",
          "If not supplied, then 10 MB will be used as a default."
      ],
      "data_chunk_size" : 10
  }
  ```
3. Fix goose table
  ``` 
  to_development=# alter table goose_db_version add column is_applied_bool bool;
  ALTER TABLE
  to_development=# \d goose_db_version;
                                          Table "public.goose_db_version"
       Column      |            Type             |                           Modifiers                           
  -----------------+-----------------------------+---------------------------------------------------------------
   id              | numeric                     | not null default nextval('goose_db_version_id_seq'::regclass)
   version_id      | bigint                      | not null
   is_applied      | smallint                    | not null
   tstamp          | timestamp without time zone | default now()
   is_applied_bool | boolean                     | 
  Indexes:
      "goose_db_version_pkey" PRIMARY KEY, btree (id)
      "public_goose_db_version_id1_idx" UNIQUE, btree (id)

  to_development=# update goose_db_version set is_applied_bool=true where is_applied=1;
  UPDATE 46
  to_development=# alter table goose_db_version drop column is_applied;
  ALTER TABLE
  to_development=# alter table goose_db_version rename column is_applied_bool to is_applied;
  ALTER TABLE
  to_development=# \d goose_db_version;
                                       Table "public.goose_db_version"
     Column   |            Type             |                           Modifiers                           
  ------------+-----------------------------+---------------------------------------------------------------
   id         | numeric                     | not null default nextval('goose_db_version_id_seq'::regclass)
   version_id | bigint                      | not null
   tstamp     | timestamp without time zone | default now()
   is_applied | boolean                     | 
  Indexes:
      "goose_db_version_pkey" PRIMARY KEY, btree (id)
      "public_goose_db_version_id1_idx" UNIQUE, btree (id)

  to_development=#
  ```


Note that migrating views DOES NOT work with this tool, there's a small syntax error (too many levels of ()) in the conversion.

##### Migrating to the latest version with goose
TBD.
