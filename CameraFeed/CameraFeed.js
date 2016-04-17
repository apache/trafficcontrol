/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

var express = require('express');
var app = express();
var https = require('https');
var fs = require('fs');
var spawn = require('child_process').spawn;

/** HTTP server listen port */
const PORT = 8080;

/** Map used to determine if a user's camera is recording or not. */
var recordingMap = new Object();

/** Debug flag for development*/
const debugFlag = process.env.DEBUG;

var debug = function(debugString) {
  if (debugFlag) {
    console.log(debugString);
  }
};

// Generic entry point for all requests.
app.use(function(req, res, next) {
  debug("---New CameraFeed request---");
  next();
});

var actionChecker = function(req, res, next) {
  debug("  Checking action...");
  var action = req.query.action;
  if (action) {
    if (action.toUpperCase() == "START" ||
	action.toUpperCase() == "STOP") {
      next();
    }
    else {
      next({code: 400, message: "Invalid 'action' parameter, " +
	    "must be either 'start' or 'stop'."});
    }
  }
  else {
    next({code: 400, message: "Missing required 'action' parameter."});
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

const user = "user"; // TODO: need real user name

/**
 * Starts recording the feed from the camera.
 */
var startRecord = function(req, res, next) {
  debug("  Starting recording executable...");
  if (! recordingMap[user]) {
    recordingMap[user] = new Object();
  }

  if (recordingMap[user][req.query.camera_id]) {
    // TODO: already recording should this be a success or failure?
    debug("    Already recording!");
  }
  else {
    // TODO: need actual username/password
    var args = ["--username=microservice",
		"--password=abc123"];
    // The camera recording exe spits out tons of debug
    if (debugFlag) {
      var debugArg = "--debug";
      if (/^\d+$/.test(debugFlag)) {
      	debugArg += "=" + debugFlag;
      }
      args.push(debugArg);
    }
    // TODO: need real IP/port
    args.push("192.168.0.9:32768");

    var options = {
    };

    debug("    " + args);

    var child = spawn("src/AmcrestIPM-721S_StreamReader", args, options);
    recordingMap[user][req.query.camera_id] = child;

    child.on('error', function(err) {
      console.log("Error with process: " + err);
    });

    child.on('exit', function(code, signal) {
      // In real application this would send an email to the user to
      // notify them something went wrong or restart the application.
      // Not here since this is just a toy.
      debug("   Process exited with code " + code + ", by signal " + signal);
      recordingMap[user][req.query.camera_id] = null;
    });

    child.stdout.on('data', function(data) {
      debug("" + data);
    });

    child.stderr.on('data', function(data) {
      debug("" + data);
    });
  }

  res.status(200).send("Recording started.");
};

/**
 * Stops recording the feed from the camera.
 */
var stopRecord = function(req, res, next) {
  debug("  Stoping recording executable...");
  if (! recordingMap[user]) {
    recordingMap[user] = new Object();
  }

  if (recordingMap[user][req.query.camera_id]) {
    recordingMap[user][req.query.camera_id].kill('SIGTERM');
    recordingMap[user][req.query.camera_id] = null;
  }
  else {
    // TODO: not recording: success or failure?
    debug("    Already not recording!");
  }
  res.status(200).send("Recording stopped.");
};

const v1 = "/v1";

// Start record
app.post(v1,
	 actionChecker,
	 cameraIdChecker,
	 startRecord);

// Stop record
app.delete(v1,
	   actionChecker,
	   cameraIdChecker,
	   stopRecord);

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
