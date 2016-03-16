/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

// Model URI:
// https://localhost:8080/CameraFeed/v1?action=[start|stop]&camera_id=[id]

var path = require('path');
var CameraFeedError = require('./CameraFeedError.js');

var moduleName = path.basename(module.id, ".js")

var debugFlag = false;

/* action values*/
const START = "start";
const STOP = "stop";

/**
 * Sets the debug flag for this module.
 * @param theDebugFlag debug value (true/false)
 */
function setDebug(theDebugFlag) {
  debugFlag = theDebugFlag;
}
exports.setDebug = setDebug;

/**
 * Pulls the action parameter from the query, verifying it is correct.
 *
 * @param query HTTP query object
 * @return action string
 * @throws CameraFeedError on any problem with the action parameter
 */
function getAction(query) {
  const actionKey = "action";
  if (query[actionKey])
  {
    var action = query[actionKey];
    if (START.localeCompare(action) != 0 &&
	STOP.localeCompare(action) != 0)
    {
      var returnJSON = CameraFeedError.buildJSON(
	CameraFeedError.InvalidAction,
	'Invalid "action" parameter provided. Must be either "' + START +
	  '" or "' + STOP + '".');
      throw new CameraFeedError.CameraFeedError(
	"Invalid Action Provided", 400, returnJSON);
    }
  }
  else
  {
    var returnJSON = CameraFeedError.buildJSON(
      CameraFeedError.InvalidAction,
      'No "action" parameter provided.');
    throw new CameraFeedError.CameraFeedError(
      "No Action Provided", 400, returnJSON);
  }

  return query[actionKey];
}

/**
 * Pulls the camera_id parameter from the query, verifying it is correct.
 *
 * @param query HTTP query object
 * @return camera id
 * @throws CameraFeedError on any problem with the camera_id parameter
 */
function getCameraId(query) {
  const cameraIdKey = "camera_id";
  if (query[cameraIdKey])
  {
    var cameraId = query[cameraIdKey];
    // TODO: need to determine the id format and verify it here
    // {
    //   var returnJSON = CameraFeedError.buildJSON(
    // 	CameraFeedError.InvalidCameraId,
    // 	'Invalid "cameraId" parameter provided. Must be either "' + START + '" or "' +
    // 	  STOP + '".');
    //   throw new CameraFeedError.CameraFeedError(
    // 	"Invalid CameraId Provided", 400, returnJSON);
    // }

    // TODO: then query DB to see if this is OK
  }
  else
  {
    var returnJSON = CameraFeedError.buildJSON(
      CameraFeedError.InvalidCameraId,
      'No "camera_id" parameter provided.');
    throw new CameraFeedError.CameraFeedError(
      "No Camera_id Provided", 400, returnJSON);
  }

  return query[cameraIdKey];
}

function executeStart(cameraId) {
  // TODO:
  // - check if camera already streaming
  //   if so, success
  //   if not, spawn child process to retrieve feed (C++, libcurl)
  //           update DB to say feed is started
  //    (need to errors from child process)
}

function executeStop(cameraId) {
}

/**
 * Processes the user's REST query for starting/stopping a camera feed.
 * @param query HTTP query object (see node.js URL.parse)
 * @param response HTTP response object
 * @throws CameraFeedError on any problem
 */
function processRequest(query, response) {
  if (debugFlag) {
    console.log("--" + moduleName +
		" Processing query:\n  " + JSON.stringify(query));
  }

  var action = getAction(query);
  var cameraId = getCameraId(query);

  if (debugFlag) {
    console.log("--Processed parameters:");
    console.log("  Action: " + action);
    console.log("  cameraId: " + cameraId);
  }

  if (START.localeCompare(action) == 0) {
    executeStart(cameraId);
  }
  else {
    executeStop(cameraId);
  }

  response.writeHead(200, {
    'Server': moduleName
  });
}
exports.processRequest = processRequest;
