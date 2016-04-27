/*
 * @author Michael Albers
 * For CSCI 5799
 */

'use strict';

var express = require('express');
var app = express();
var https = require('https');
var mongoClient = require('mongodb').MongoClient;
var sprintf = require('sprintf-js');
var avconv = require('avconv');
const fs = require('fs');

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

/**
 * Entry point for all requests
 */
app.use(function(req, res, next) {
  debug("--New Video Retrieval--");
  next();
});

/**
 * Verifies the existence of the start & stop parameters then converts them
 * to the format needed for the MongoDB query.
 */
var dateConvert = function(req, res, next) {
  var start = req.query.start;
  var stop = req.query.stop;
  debug("  Converting dates...");
  if (start && stop) {
    req.query.start = Date.parse(start);
    req.query.stop = Date.parse(stop);
    if (isNaN(req.query.start) || isNaN(req.query.stop))
    {
      next({code: 400, message: "Invalid start or stop date format."});
    }
    else {
      debug("   start: " + req.query.start + ", stop: " + req.query.stop);
      next();
    }
  }
  else if (! start) {
    next({code: 400, message: "Missing required 'start' parameter."});
  }
  else if (! stop) {
    next({code: 400, message: "Missing required 'stop' parameter."});
  }
  else {
    next({code: 400,
	  message: "Missing required 'start' and 'stop' parameters."});
  }
};

/**
 * Pulls the JPEG data from the MongoDB and writes them to a temporary
 * directory.
 */
var jpegRetriever = function(req, res, next) {
  debug("  Retrieving video...");

  var user = req.params.user;
  var cameraId = req.params.cameraId;
  debug("    User: " + user + ", Camera: " + cameraId);

  var template = "/tmp/RetrieveVideo_" + cameraId + '_' + user;
  req.tmpDir = fs.mkdtempSync(template);
  req.fileList = new Array();

  var url = "mongodb://" + mongo.host + ":" + mongo.port + "/CSCI5799";

  mongoClient.connect(url, function(err, db) {
    if (err) {
      next("Error connecting to MongoDB " + url + ". Error: " + err);
    }
    else {
      debug("  Connected to MongoDB " + url + ".");
      var collection = db.collection('CameraFeed');
      // https://docs.mongodb.org/getting-started/node/query/
      // https://docs.mongodb.org/manual/reference/operator/query
      var query = {
	"user": user,
	"camera_id": cameraId,
	$and : [
	  { $and: [ {"msSinceEpoch": {$gte: req.query.start}},
		    {"msSinceEpoch": {$lte: req.query.stop}} ] } ]
      };

      var cursor = collection.find(query).sort({"msSinceEpoch": 1});
      var width = 0;
      cursor.count(function(err, count) {
	var numberImages = 0;
	if (! err) {
	  numberImages = count;
	  width = Math.ceil(Math.log(numberImages) / Math.log(10));
	  req.fileFormat = "input%0" + width + "d.jpg";
	  debug("  Found " + numberImages + " images...");
	}

	// This has to be in the count callback otherwise numberImages
	// will likely be zero due to asynchronous operations.
	if (numberImages > 0) {
	  var imageNumber = 0;
	  cursor.each(function(err, doc) {
	    if (err) {
	      next("Error retrieving video from DB: " + err);
	    }
	    else if (doc) {
	      imageNumber++;
	      var fileName = sprintf.sprintf(req.fileFormat, imageNumber);
	      fileName = req.tmpDir + "/" + fileName;
	      var jpegFile = fs.createWriteStream(fileName);
	      jpegFile.write(doc.jpeg.buffer);
	      jpegFile.end();
	      req.fileList.push(fileName);
	    }
	    else {
	      debug("  Done writing " + numberImages + " image(s) to disk...");
	      db.close();
	      next();
	    }
	  });
	}
	else {
      	  next({code: 400, message: "No video found for given time range."});
	}
      }); // End of cursor.count
    }
  });
};

/**
 * Assembles the JPEGs into an MP4.
 */
var mpeg4Builder = function(req, res, next) {
  debug("  Building MP4 in " + req.tmpDir + "...");
  req.tmpFile = req.tmpDir + "/" + "Video.mp4";
  // avconv -r 25 -i input%03d.jpg test.mp4
  var avconvParams = [
    "-r", 25,
    "-i", req.tmpDir + "/" + req.fileFormat,
    req.tmpFile
    ];

  var avconvStream = avconv(avconvParams);
  avconvStream.on('exit', function(exitCode, signal, metadata) {
    if (exitCode == 0) {
      debug("  Finished encoding movie...");
      next();
    }
    else {
      next("Error with avconv: " + exitCode + ", signal: " + signal +
	   ", meta: " + metadata);
    }
  });
};

/**
 * Cleans up the temporary files/directory used to create the video.
 */
var cleanup = function(req) {
  if (req.fileList) {
    for (var ii = 0; ii < req.fileList.length; ++ii) {
      fs.unlinkSync(req.fileList[ii]);
    }
  }
  if (req.tmpFile) {
    fs.unlinkSync(req.tmpFile);
  }
  if (req.tmpDir) {
    fs.rmdirSync(req.tmpDir);
  }
};

/**
 * Ships the MP4 to the caller and on completion removes temporary files.
 */
var downloadVideo = function (req, res) {
  debug("  Downloading video...");
  res.download(req.tmpFile, "Video.mp4", function(err){
    if (err) {
      console.log("Error: " + err);
    }
    else {
      debug("  Video downloaded...");
    }
    cleanup(req);
  });
};

app.get('/video/:user/:cameraId',
	dateConvert,
	jpegRetriever,
	mpeg4Builder,
	downloadVideo);

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
  cleanup(req);
};

app.use(errorHandler);

const options = {
  key: fs.readFileSync('certs/key.pem'),
  cert: fs.readFileSync('certs/cert.pem')
};

https.createServer(options, app).listen(PORT);
