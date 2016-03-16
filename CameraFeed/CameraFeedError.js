/*
 * @author Michael Albers
 * For CSCI 5799
 */

"use strict";

/*
 * Error codes
 */
exports.InvalidPath = 100;
exports.InvalidResource = 101;
exports.InvalidAPI = 102;
exports.InvalidAction = 103;
exports.InvalidCameraId = 104;

/*
 * Builds the JSON object for CameraFeedError
 *
 * @param errorCode One of the error codes listed above.
 * @param message String describing the error
 * @return object for CameraFeedError
 */
function buildJSON(errorCode, message) {
  return {
    "code": errorCode,
    "message": message
  }
}
exports.buildJSON = buildJSON;

/*
 * General purpose error class for CameraFeed microservice.
 */
class CameraFeedError extends Error {
  constructor(name, httpCode, returnJSON) {
    super(name);
    this.httpCode = httpCode;
    this.returnJSON = returnJSON;
  }

  getHttpCode() {
    return this.httpCode;
  }

  getJSON() {
    return this.returnJSON;
  }
}

exports.CameraFeedError = CameraFeedError;
