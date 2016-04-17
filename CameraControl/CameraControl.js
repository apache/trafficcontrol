/*
 * @author Michael Albers
 * For CSCI 5799
 */

'use strict';

var express = require('express');
var app = express();
var http = require('http');
var https = require('https');
var fs = require('fs');

const PORT = 8080;

const debugFlag = process.env.DEBUG;

var debug = function(debugString) {
  if (debugFlag) {
    console.log(debugString);
  }
};

// Generic entry point for all requests.
app.use(function(req, res, next) {
  debug("---New PTZ request---");
  next();
});

/**
 * Checks for the existence and valid values of the 'action' parameter.
 */
var actionChecker = function(req, res, next) {
  debug("  Checking action...");
  var action = req.query.action;
  if (action) {
    if (action.toUpperCase() == "START" ||
	action.toUpperCase() == "STOP") {
      // Camera only accepts lower case
      req.query.action = action.toLowerCase();
      next();
    }
    else {
      next({code: 400, message: "Invalid 'action' parameter, must be " +
	    "either 'start' or 'stop'."});
    }
  }
  else {
    next({code: 400, message: "Missing required 'action' parameter."});
  }
};

/**
 * Checks for the existence and valid values of the 'direction' parameter.
 */
var directionChecker = function(req, res, next) {
  debug("  Checking direction...");
  var direction = req.query.direction;
  if (direction) {
    if (direction.toUpperCase() == "UP" ||
	direction.toUpperCase() == "DOWN" ||
	direction.toUpperCase() == "LEFT" ||
	direction.toUpperCase() == "RIGHT") {
      // Camera only accepts direction with upper-cased first letter (like
      // proper noun)
      req.query.direction = direction.charAt(0).toUpperCase() +
	direction.slice(1);
      next();
    }
    else {
      next({code: 400,
	    message: "Invalid 'direction' parameter, must be one of: " +
	    "'up', 'down', 'left' or 'right'."});
    }
  }
  else {
    next({code: 400, message: "Missing required 'direction' parameter."});
  }
};

/**
 * Checks for the existence and valid values of the 'velocity' parameter.
 */
var velocityChecker = function(req, res, next) {
  debug("  Checking velocity...");
  var velocity = req.query.velocity;
  if (velocity) {
    var integerTest = /^[1-8]$/;
    if (integerTest.test(velocity)) {
      next();
    }
    else {
      next({code: 400, message: "Invalid 'velocity' parameter, " +
	    "must be an integer between 1 and 8, inclusive."});
    }
  }
  else {
    next({code: 400, message: "Missing required 'velocity' parameter."});
  }
};

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

/**
 * Builds and sends final control command to the camers.
 */
var sendCommand = function(req, res, next) {
  debug("  Building command for camera...");
  var hostName = '192.168.0.9';
  var portNum = 32768;
  var path = "/cgi-bin/ptz.cgi?"
  var userPassword = "microservice:abc123";

  path += "action=" + req.query.action;
  path += "&channel=0";
  path += "&code=" + req.query.direction;
  path += "&arg1=0";
  path += "&arg2=" + req.query.velocity;
  path += "&arg3=0&arg4=0";

  var options = {
    hostname: hostName,
    port: portNum,
    method: "GET",
    auth: userPassword,
    path: path
  };

  var req = http.request(options, function(res) {
    debug("  Sending command to camera...");
    // Apparently need to handle this event otherwise the 'end' event
    // doesn't get emitted.
    res.on('data', (chunk) => {
      debug("  Response data from camera '" + new String(chunk) + "'");
    });
    res.on('end', () => {
      debug("  Code from camera: " + res.statusCode);
      if (res.statusCode == 200) {
	next();
      }
      else {
	next("Error from camera. Code: " + res.statusCode);
      }
    });
  });

  req.on('error', function(e) {
    next("Error sending command to camera: " + e.message);
  });

  req.end();
};

/**
 * Sends success response back to the user.
 */
var finishRequest = function(req, res) {
  debug("  Finishing request");
  res.status(200).send("PTZ Done!");
};

app.post("/ptz/v1",
	 actionChecker,
	 directionChecker,
	 velocityChecker,
	 cameraIdChecker,
	 sendCommand,
	 finishRequest);

/**
 * Error handler routine (has to come last).
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
