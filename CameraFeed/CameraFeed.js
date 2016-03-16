/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

var https = require("https");
var url = require("url");
const fs = require('fs');
var CameraFeedError = require('./CameraFeedError.js');

// TODO: remember to allow security exception until better certs are had

// Model URI:
// https://localhost:8080/CameraFeed/v1?<some parameters>

/** Debug flag for development*/
var debugFlag = true;

/** HTTP server listen port */
const PORT = 8080;

/** REST resources */
const CAMERA_RESOURCE = "CameraFeed"

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
 * Given the pathname of a URI (i.e., /CameraFeed/v1) verifies that
 * the resource is valid and returns the API version.
 *
 * @param path URI pathname
 * @return API version string
 * @throws CameraFeedError on invalid pathname or resource
 */
var verifyResourceVersion = function(path) {
  var paths = path.split("/");
  if (paths.length != 3)
  {
    var returnJSON = CameraFeedError.buildJSON(
      CameraFeedError.InvalidPath,
      'Resource path did not contain correct number of parts. Must be /' +
	CAMERA_RESOURCE + '/<version>');
    throw new CameraFeedError.CameraFeedError(
      "Invalid Resource", 400, returnJSON);
  }

  var resource = paths[1];
  var apiVersion = paths[2];

  if (CAMERA_RESOURCE.localeCompare(resource) != 0)
  {
    var returnJSON = CameraFeedError.buildJSON(
      CameraFeedError.InvalidPath,
      'Invalid resource. Must be ' + CAMERA_RESOURCE);
    throw new CameraFeedError.CameraFeedError(
      "Invalid Resource", 400, returnJSON);
  }

  return apiVersion;
}

/**
 * Request event handler for HttpsServer object.
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
    var apiVersion = verifyResourceVersion(parsedURL.pathname);

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
      console.log(e.name);
      console.log(e.stack);
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

const options = {
  key: fs.readFileSync('keys/key.pem'),
  cert: fs.readFileSync('keys/cert.pem')
};

https.createServer(options, requestHandler).listen(PORT);

// TODO: next need to register with API gateway
