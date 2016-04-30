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

/** Debug flag for development*/
const debugFlag = process.env.DEBUG;

/** HTTP server listen port */
const PORT = debugFlag ? 8080 : 443;

var debug = function(debugString) {
  if (debugFlag) {
    console.log(debugString);
  }
};

// Read in /cameras service information
const cameras = JSON.parse(fs.readFileSync('cameras.json', 'utf8'));
debug("Cameras: " + JSON.stringify(cameras));

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
 * Retrieves camera URL/user/password information.
 */
var cameraData = function(req, res, next) {
  debug("  Getting Camera data...");
  var user = req.params.user;
  var cameraId = req.params.cameraId;

  var path = "https://" + cameras.host + ":" + cameras.port + "/cameras/" +
    user + "/" + cameraId;
  // http://stackoverflow.com/questions/10888610/ignore-invalid-self-signed-ssl-certificate-in-node-js-with-https-request
  // "Cheap and insecure"
  process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";

  https.get(path, (cameraRes) => {
    var body = "";

    cameraRes.on('data', function(chunk) {
      body += chunk;
    });

    cameraRes.on('end', function() {
      var cameraData = JSON.parse(body);
      debug(cameraData);
      if (cameraData.Status && cameraData.Status.toUpperCase() == "SUCCESS") {
	req.CameraURL = cameraData.CameraData[0].url;
	req.CameraUser = cameraData.CameraData[0].username;
 	req.CameraPassword = cameraData.CameraData[0].password;
	next();
      }
      else {
	next({code: 400, message: "No such camera."});
      }
    });
  }).on('error', (e) => {
    next({code: 500, message: "Error getting camera info: " + e});
  });
};

/**
 * Builds and sends final control command to the camers.
 */
var sendCommand = function(req, res, next) {
  debug("  Building command for camera...");

  var hostName = req.CameraURL;
  var portNum = 32768; // TODO: not on camera info
  var path = "/cgi-bin/ptz.cgi?"
  var userPassword = req.CameraUser + ":" + req.CameraPassword;

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

  debug(options);

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

app.post("/control/:user/:cameraId",
	 actionChecker,
	 directionChecker,
	 velocityChecker,
	 cameraData,
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
