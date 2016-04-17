/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

var express = require('express');
var app = express();
var https = require('https');
var fs = require('fs');

/** HTTP server listen port */
const PORT = 8080;

/** Debug flag for development*/
const debugFlag = process.env.DEBUG;

var debug = function(debugString) {
  if (debugFlag) {
    console.log(debugString);
  }
};

// Generic entry point for all requests.
app.use(function(req, res, next) {
  debug("---New LiveStream request---");
  next();
});

/**
 * Verifies the existence of camera_id parameter.
 */
var cameraIdChecker = function(req, res, next) {
  debug("  Checking camera_id parameter...");
  if (req.query.camera_id) {
    next();
  }
  else
  {
    next({code: 400, message: "Missing required 'camera_id' parameter."});
  }
};

var returnURL = function(req, res, next) {
  debug("  Returning URL...");
  // TODO: need real IP/port
  var url = "192.168.0.1:32768/cgi-bin/mjpg/video.cgi?channel=0&subtype=1";
  res.status(200).send(url);
};

app.get("/LiveStream/v1",
	cameraIdChecker,
	returnURL);

/**
 * Error handler (has to be last).
 */
var errorHandler = function(err, req, res, next) {
  var code = 500;
  var message = "Unhandled error.";
  if (typeof err == "object") {
    if (err["code"]) {
      code = err["code"];
    }
    if (err["message"]) {
      message = err["message"];
    }
  }
  else {
    message = err;
  }
  res.status(code).send(message);
};

app.use(errorHandler);

const options = {
  key: fs.readFileSync('certs/key.pem'),
  cert: fs.readFileSync('certs/cert.pem')
};

https.createServer(options, app).listen(PORT);
