# CameraFeed Microservice

## Description

This microservice handles recording the live video stream from the camera.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information. (Example: export DEBUG=1)

## Connection
HTTPS on port 8080.

## API
	/CameraFeed/v1?[args...]

## API Versions
### Version 1
    POST /CameraFeed/v1?camera_id=[id]
    DELETE /CameraFeed/v1?camera_id=[id]

POST starts recording. DELETE stops recording. If already recording and a second POST is sent, the command is ignored and success returned. The same is true for sending a DELETE after recording has stopped.

#### Parameters
* **camera_id** ID of the camera from which to start or stop recording

#### Response
On success, HTTP response with or without any information in the body. On error, HTTP response with error message in the body.

#### Example
	POST /CameraFeed/v1?camera_id=12345
	DELETE /CameraFeed/v1?camera_id=12345
