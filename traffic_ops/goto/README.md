#GoTO (Golang Traffic Ops)
##A web API for SQL databases

GoTO is a server/some other stuff written in Go that allows for RESTful interaction with SQL databases through an Angular web API.

This is written for the Comcast [Traffic Ops](http://traffic-control-cdn.net/docs/latest/development/traffic_ops.html) database, but I'm pretty sure it should probably work for all databases.

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
####Posting to a table
  ```
  curl -X POST --data "filename=$YOURFILEHERE" http://127.0.0.1:8080/tableName
  ```
For now, if you want to post a new row to a table, you need to have everything in a JSON file (`$YOURFILEHERE`)
  in your GoTO directory.

  Eventually, information added through the front-end will be passed as JSON data. 

  You can see examples of POST files in `testFiles`; specifically, `newAsn` and `newAsns`.

####Custom views
  ```
  curl -X POST --data "filename=$YOURFILEHERE" http://127.0.0.1:8080/
  ```

  Say you only want certain columns from a table, or you want a complex SQL query like a join.
  For now, you make a JSON file ($YOURFILEHERE) that follows the following format:
  ```
  [{"name":"viewName", "query":"select foo.id, bar.name from foo join bar"}]
  ```
  You can add multiple views if you'd like. Then, you can interact with the view like a table.

  You can see examples of POST files in `testFiles`; specifically, `newView` and `newViewss`.

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
  * `jmoiron/sqlx` has been super useful. Thanks!
