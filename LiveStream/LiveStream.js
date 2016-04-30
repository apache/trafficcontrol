/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

var express = require('express');
var app = express();
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
  debug("---New LiveStream request---");
  next();
});

var returnURL = function(req, res, next) {
  debug("  Returning URL...");
  var user = req.params.user;
  var cameraId = req.params.cameraId;
  debug("  User: " + user + ", Camera: " + cameraId);

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
	// TODO: will need port
	var url = cameraData.CameraData[0].url +
	  "/cgi-bin/mjpg/video.cgi?channel=0&subtype=1";
	res.status(200).send(url);
      }
      else {
	res.status(400).send("No such camera.");
      }
    });
  }).on('error', (e) => {
    res.status(500).send("Error getting camera info: " + e);
  });
};

app.get("/livestream/:user/:cameraId",
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
