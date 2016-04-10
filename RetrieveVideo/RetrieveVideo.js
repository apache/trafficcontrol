'use strict';

var express = require('express');
var app = express();
var mongoClient = require('mongodb').MongoClient;
var sprintf = require('sprintf-js');
var avconv = require('avconv');
const fs = require('fs');

const debugFlag = process.env.DEBUG;

var debug = function(debugString) {
  if (debugFlag) {
    console.log(debugString);
  }
};

var dateConvert = function(req, res, next) {
  var start = req.query.start;
  var stop = req.query.stop;
  debug("Converting... start: " + start + ", stop: " + stop);
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

app.use(dateConvert);

var cameraIdChecker = function(req, res, next) {
  debug("Checking camera_id parameter...");
  if (req.query.camera_id) {
    next();
  }
  else
  {
    next({code: 400, message: "Missing required 'camera_id' parameter."});
  }
};

app.use(cameraIdChecker);

var jpegRetriever = function(req, res, next) {
  debug("Retrieving video...");
  var user = "dummyUser"; // TODO: need user
  var template = "/tmp/RetrieveVideo_" + req.query.camera_id + '_' + user;
  req.tmpDir = fs.mkdtempSync(template);
  req.fileList = new Array();

  // TODO: need actual host/port
  var url = "mongodb://localhost:27017/CSCI5799";

  mongoClient.connect(url, function(err, db) {
    if (err) {
      next("Error connecting to MongoDB " + url + ". Error: " + err);
    }
    else {
      debug("Connected to MongoDB " + url + ".");
      var collection = db.collection('CameraFeed');
      // https://docs.mongodb.org/getting-started/node/query/
      // https://docs.mongodb.org/manual/reference/operator/query
      var query = {
	"user": user,
	"camera_id": req.query.camera_id,
	$and : [
	  { $and: [ {"msSinceEpoch": {$gte: req.query.start}},
		    {"msSinceEpoch": {$lte: req.query.stop}} ] } ]
      };
      debug("Query: " + JSON.stringify(query));
      var cursor = collection.find(query).sort({"msSinceEpoch": 1});
      var numberImages = 0;
      var width = 0;
      cursor.count(function(err, count) {
	if (! err) {
	  numberImages = count;
	  width = Math.ceil(Math.log(numberImages) / Math.log(10));
	  req.fileFormat = "input%0" + width + "d.jpg";
	  debug("Found " + numberImages + " images...");
	}
      });

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
	    debug("Done writing " + numberImages + " image(s) to disk...");
	    db.close();
	    next();
	  }
	});
      }
      else {
	next({code: 400, message: "No video found for given time range."});
      }
    }
  });
};

app.use(jpegRetriever);

var mpeg4Builder = function(req, res, next) {
  debug("Building MP4 in " + req.tmpDir + "...");
  req.tmpFile = req.tmpDir + "/" + "Video.mp4";
  // avconv -r 25 -i input%03d.jpg test.mp4
  var avconvParams = [
    "-r", 25,
    "-i", req.tmpDir + "/" + req.fileFormat,
    req.tmpFile
    ];
  debug(avconvParams);
  var avconvStream = avconv(avconvParams);
  avconvStream.on('exit', function(exitCode, signal, metadata) {
    if (exitCode == 0) {
      debug("Finished encoding movie...");
      next();
    }
    else {
      next("Error with avconv: " + exitCode + ", signal: " + signal +
	   ", meta: " + metadata);
    }
  });
};

app.use(mpeg4Builder);

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

// TODO: version
app.get('/video', function (req, res) {
  debug("Downloading video...");
  res.download(req.tmpFile, "Video.mp4", function(err){
    if (err) {
      console.log("Error: " + err);
    }
    else {
      debug("Video downloaded...");
    }
    for (var ii = 0; ii < req.fileList.length; ++ii) {
      fs.unlinkSync(req.fileList[ii]);
    }
    fs.unlinkSync(req.tmpFile);
    fs.rmdirSync(req.tmpDir);
  });
});

const PORT = 8080;

app.listen(PORT, function () {
  debug('RetrieveVideo listening on port ' + PORT + '!');
});
