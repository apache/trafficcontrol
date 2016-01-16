##### Intro
This is a prototype for Traffic Ops 2.0 server.  See also https://github.com/Comcast/traffic_control/wiki/Traffic_Ops_20

##### One time generation of CRUD files 
  ```
  [jvd@laika tools (master *=)]$ go run gen_goto2.go root ******* to_development
  [asn cachegroup cachegroup_parameter cdn deliveryservice deliveryservice_regex deliveryservice_server deliveryservice_tmuser division federation federation_deliveryservice federation_federation_resolver federation_resolver federation_tmuser goose_db_version hwinfo job job_agent job_result job_status log parameter phys_location profile profile_parameter regex region role server servercheck staticdnsentry stats_summary status tm_user to_extension type]
  asn: Ok 4809
  cachegroup: Ok 6064
  cachegroup_parameter: Ok 5644
  cdn: Ok 4882
  deliveryservice: Ok 11863
  deliveryservice_regex: Ok 5768
  deliveryservice_server: Ok 5787
  deliveryservice_tmuser: Ok 5811
  division: Ok 4905
  federation: Ok 5295
  federation_deliveryservice: Ok 5968
  federation_federation_resolver: Ok 6150
  federation_resolver: Ok 5535
  federation_tmuser: Ok 5592
  goose_db_version: Ok 5387
  hwinfo: Ok 5092
  job: Ok 6364
  job_agent: Ok 5235
  job_result: Ok 5382
  job_status: Ok 5163
  log: Ok 5110
  parameter: Ok 5218
  phys_location: Ok 6328
  profile: Ok 5049
  profile_parameter: Ok 5461
  regex: Ok 4890
  region: Ok 4944
  role: Ok 4909
  server: Ok 9382
  servercheck: Ok 8004
  staticdnsentry: Ok 5880
  stats_summary: Ok 5804
  status: Ok 5002
  tm_user: Ok 7678
  to_extension: Ok 6608
  type: Ok 5071
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

2. Generate the docs.go file and copy it to where the swagger UI files live:
  ```
  jvd@laika playground (master *%=)]$ swagger -apiPackage github.com/Comcast/traffic_control/traffic_ops/experimental/server/api -mainApiFile github.com/Comcast/traffic_control/traffic_ops/experimental/server/api/action.go -format go
  2016/01/14 07:47:27 Start parsing
  2016/01/14 07:47:29 Finish parsing
  2016/01/14 07:47:29 Doc file generated
  [jvd@laika playground (master *%=)]$ ls -ltr
  total 696
  -rw-r--r--   1 jvd  staff    1899 Jan 12 12:04 web.go
  -rw-r--r--   1 jvd  staff     991 Jan 12 12:04 notes.txt
  drwxr-xr-x   3 jvd  staff     102 Jan 12 12:04 db
  drwxr-xr-x   4 jvd  staff     136 Jan 12 12:04 client
  drwxr-xr-x   3 jvd  staff     102 Jan 13 09:05 routes
  drwxr-xr-x   3 jvd  staff     102 Jan 13 09:05 output_format
  drwxr-xr-x   3 jvd  staff     102 Jan 13 09:05 csconfig
  drwxr-xr-x   3 jvd  staff     102 Jan 13 09:05 crconfig
  drwxr-xr-x   3 jvd  staff     102 Jan 13 09:05 auth
  drwxr-xr-x  41 jvd  staff    1394 Jan 13 09:05 api
  -rw-r--r--   1 jvd  staff    3267 Jan 13 09:05 main.go
  drwxr-xr-x   5 jvd  staff     170 Jan 13 15:39 tools
  drwxr-xr-x   5 jvd  staff     170 Jan 13 17:41 conf
  -rw-r--r--   1 jvd  staff    7663 Jan 14 07:46 README.md
  -rw-r--r--   1 jvd  staff  335376 Jan 14 07:47 docs.go
  [jvd@laika playground (master *%=)]$ cp docs.go ../../swagger-api/swagger-ui/
  ```

3. Run the web.go app:
  ```
  [jvd@laika swagger-ui (master *%=)]$ go run web.go docs.go -port 8081 -api http://localhost:8080 -staticPath ./dist
  2016/01/14 07:53:20 62372 0.0.0.0:8081
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
