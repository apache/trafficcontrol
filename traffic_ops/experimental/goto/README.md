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

#GoTO (Golang Traffic Ops)
##A web API for SQL databases

GoTO is a server/some other stuff written in Go that allows for RESTful interaction with SQL databases through an Angular web API.

This is written for the Comcast [Traffic Ops](http://trafficcontrol.apache.org/docs/latest/development/traffic_ops.html) database, but I'm pretty sure it should probably work for all databases.

## Install

1. First, fork a copy of this sick repo. "GoTO" a directory of your choice and type in

```
git clone https://github.com/comcast/traffic_control.git

```
2. Then, make a `.dbInfo` file that follows this syntax, 
  replacing the content in brackets with your own data:
  ```
  USERNAME="[databaseUsername]"
  PASSWORD="[databasePassword]"
  DATABASE="[databaseName]"
  ```
  For example, if you want to work with the `foo` database with username `johndoe` and password `password`, 
  your `.dbInfo file should look like this:
  ```
  USERNAME="johndoe"
  PASSWORD="password"
  DATABASE="foo"
  ```
  3. Now, you can run the server by typing this into your terminal:
  ```
  ./run
  ```
  You should get a `Starting server.` message. If all goes well, the server should start on port 8080! You can change this in `server.go`.

  Alternatively, if the Bash script isn't working out for you (i.e. you're using Windows or something weird like that), you can start the server with:
  ```
  go run server.go $USERNAME $PASSWORD $DATABASE
  ```
(Filling in $USERNAME, $PASSWORD, and $DATABASE with your own creds, of course.)

  Then, start up the Angular front-end by running
  ```
  python -m SimpleHTTPServer
  ```

  Should be up and running on :8000! Make sure ./run is still going concurrently.

## Debugging
  If you're getting errors in the Install process or you happen to be Mark, make sure you can answer "yes" to
  the following questions. If you're still having issues, that really sucks.
  * Do you have the most recent version of Go [installed](https://golang.org/doc/install)? Try uninstalling/reinstalling.
  * Did you make a `.dbInfo` file? (See step two of the [Install](http://github.com/cjqian/GoTO#install) notes.)
  * Are you running `./run` from your `GoTO/` folder and not a subfolder?

  See `./run` for execution examples. Also, are your database credentials correct?
  * Is your `mysql` up and running? Type `mysql` into your terminal to verify.
  * Do you have the latest version of this code? Run `git pull` to get an update. 
  * Also, make sure you've checked out `master` branch and not a development branch.

## Syntax 
###GET 
  ```
  curl http://127.0.0.1:8080/tableName?columnAValue=1&columnBValue<50/id
  ```
  You can get a row of the `tableName` by `id`, or by other column values. 
  Column value parameters follow a `?` after the `tableName`, and are separated by `&`. 

  Examples:
  ```
  table?value<100
  ```

  ```
  table?value<100&value2>=100
  ```
###POST
####API
Currently, the API only supports adding one view or one row to a table at a time. If you wish to add multiple, you'll need to pass a JSON array of the many views or many queries to a curl POST. The file needs to be of the following form: 
 ```
  [{"name":"viewName", "query":"select foo.id, bar.name from foo join bar"}]
  ```

  You can see examples of POST files in `testFiles`; specifically, `newView` and `newViews` for new views, or `newAsn` or `newAsns` for new rows.

###PUT
  Put follows the same syntax as POST (but with PUT). On the SQL end, "UPDATES." Can be done with views, too.

###DELETE

  In this example, all rows from database `foo` with `swag < 100` are deleted. 
  ```
  curl -X DELETE http//127.0.0.1:8080/foo?swag<100/
  ```

  In this example,  row with id `1` from database `foo` is deleted. 
  ```
  curl -X DELETE http//127.0.0.1:8080/foo/1
  ```

  What happens if you do something like this?
  ```
  curl -X DELETE http//127.0.0.1:8080/foo
  ```

* If `foo` is a table, all the rows are deleted from the table. You cannot drop the table, just like you cannot drop the bass. (Ooh, sick burn.)
* If `foo` is a view, it is dropped.

  This was an arbitrary decision on my part. Let me know if you have more elegant solutions.

##Packages
###Local
  * sqlParser processes all interactions with the database. It contains `sqlParser.go`, which contains most of the CRUD methods, and `sqlTypeMap`, which has functions mapping values of type interface{} to string and vice-versa.
  * urlParser parses the url into a Request.
  * outputFormatter wraps the query into an encodable struct.

  There are more details in the comments of each of these packages.
###Other
  * I'm also using AngularJS, jQuery, Bootstrap.
  * `jmoiron/sqlx` has been super useful. Thanks!
