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

/** Map used to determine if a user's camera is recording or not. */
var recordingMap = new Object();

/** Debug flag for development*/
const debugFlag = process.env.DEBUG;

/** HTTP server listen port */
const PORT = debugFlag ? 8080 : 443;

var debug = function(debugString) {
  if (debugFlag) {
    console.log(debugString);
  }
};

// Read in Mongo DB information
const mongo = JSON.parse(fs.readFileSync('mongo.json', 'utf8'));
debug("Mongo: " + JSON.stringify(mongo));

// Read in /cameras service information
const cameras = JSON.parse(fs.readFileSync('cameras.json', 'utf8'));
debug("Cameras: " + JSON.stringify(cameras));

// Generic entry point for all requests.
app.use(function(req, res, next) {
  debug("---New CameraFeed request---");
  next();
});

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
 * Starts recording the feed from the camera.
 */
var startRecord = function(req, res, next) {
  debug("  Starting recording executable...");

  var user = req.params.user;
  var cameraId = req.params.cameraId;
  debug("    User: " + user + ", Camera: " + cameraId);

  if (! recordingMap[user]) {
    recordingMap[user] = new Object();
  }

  if (recordingMap[user][cameraId]) {
    debug("    Already recording!");
  }
  else {
    var mongoArg = "--mongo=" + mongo.host + ":" + mongo.port;

    var args = ["--username=" + req.CameraUser,
		"--password=" + req.CameraPassword,
		"--user=" + user,
		"--camera=" + cameraId,
		mongoArg];
    // The camera recording exe spits out tons of debug
    if (debugFlag) {
      var debugArg = "--debug";
      if (/^\d+$/.test(debugFlag)) {
      	debugArg += "=" + debugFlag;
      }
      args.push(debugArg);
    }

    // TODO: will neeed port number too
    args.push(req.CameraURL);

    var options = {
    };

    debug("    " + args);

    var child = spawn("src/AmcrestIPM-721S_StreamReader", args, options);
    recordingMap[user][cameraId] = child;

    child.on('error', function(err) {
      console.log("Error with process: " + err);
    });

    child.on('exit', function(code, signal) {
      // In real application this would send an email to the user to
      // notify them something went wrong or restart the application.
      // Not here since this is just a toy.
      debug("   Process exited with code " + code + ", by signal " + signal);
      recordingMap[user][cameraId] = null;
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

  var user = req.params.user;
  var cameraId = req.params.cameraId;
  debug("    User: " + user + ", Camera: " + cameraId);

  if (! recordingMap[user]) {
    recordingMap[user] = new Object();
  }

  if (recordingMap[user][cameraId]) {
    recordingMap[user][cameraId].kill('SIGTERM');
    recordingMap[user][cameraId] = null;
  }
  else {
    debug("    Already not recording!");
  }
  res.status(200).send("Recording stopped.");
};

const service = "/feed/:user/:cameraId";

// Start record
app.post(service,
	 cameraData,
	 startRecord);

// Stop record
app.delete(service,
	   cameraData,
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
