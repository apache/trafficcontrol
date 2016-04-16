/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

var http = require("http");
var url = require("url");
const fs = require('fs');
var CameraFeedError = require('./CameraFeedError.js');

/** Debug flag for development*/
var debugFlag = true;

/** HTTP server listen port */
const PORT = 8080;

/** Name of the server for HTTP header */
const SERVER = "CameraFeed/0.1"

/*
 * API versions
 */
var CameraFeed_v1 = require('./CameraFeed_v1.js');
var apiMap = {
  'v1' : CameraFeed_v1
};

/*
 * Given the pathname of a URI (i.e., /v1) returns the API version.
 *
 * @param path URI pathname
 * @return API version string
 * @throws CameraFeedError on invalid pathname
 */
var getVersion = function(path) {
  var paths = path.split("/");
  if (paths.length != 2)
  {
    var returnJSON = CameraFeedError.buildJSON(
      CameraFeedError.InvalidPath,
      'Resource path did not contain correct number of parts. ' +
	'Must be /<version>');
    throw new CameraFeedError.CameraFeedError(
      "Invalid Resource", 400, returnJSON);
  }

  var apiVersion = paths[1];
  return apiVersion;
}

/**
 * Request event handler for HttpServer object.
 * See https://nodejs.org/api/http.html#http_event_request
 */
var requestHandler = function(request, response)
{
  var parsedURL = url.parse(request.url, true);

  if (debugFlag) {
    console.log("===Request===");
    console.log("Date: " + new Date().toJSON());
    console.log("URL: " + parsedURL.pathname);
    console.log("Method: " + request.method); // GET, POST, etc.
    console.log("Headers: " + JSON.stringify(request.headers));
    console.log("Query: " + JSON.stringify(parsedURL.query));
    console.log("===END===");
  }

  try
  {
    var apiVersion = getVersion(parsedURL.pathname);

    if (debugFlag) {
      console.log("--Client requested API version: " + apiVersion);
    }

    if (apiMap[apiVersion])
    {
      apiMap[apiVersion].setDebug(debugFlag);
      apiMap[apiVersion].processRequest(parsedURL.query, response);
    }
    else
    {
      var returnJSON = CameraFeedError.buildJSON(
	CameraFeedError.InvalidPath,
	'Invalid API version requested: ' + apiVersion + '. Must be one of: ' +
	  Object.keys(apiMap));
      throw new CameraFeedError.CameraFeedError(
	"Invalid API version", 400, returnJSON);
    }
  }
  catch (e)
  {
    // This helps in development when some other exception besides CameraFeedError
    // might be getting thrown.
    if (debugFlag) {
      console.log("------------------");
      console.log("Debug output for exception thrown during request");
      console.log("Name: " + e.name);
      console.log("Stack:");
      console.log(e.stack);
      console.log("------------------");
    }

    response.writeHead(e.getHttpCode(), {
      'Content-Type': 'application/JSON',
      'Server': SERVER
    });
    response.write(JSON.stringify(e.getJSON()));
  }
  response.end();
}

if (debugFlag) {
  console.log(new Date().toJSON());
  console.log("Starting CameraFeed microservice...");
}

http.createServer(requestHandler).listen(PORT);
