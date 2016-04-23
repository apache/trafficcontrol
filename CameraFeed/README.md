# CameraFeed Microservice

## Description

This microservice handles recording the live video stream from the camera.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information. (Example: export DEBUG=1)

## Connection
HTTPS on port 443, 8080 when in debug mode (gets around having to run it as root to open 443).

## API
	/feed

## API Versions
### Version 1
    POST /feed/{user}/{camera}
    DELETE /feed/{user}/{camera}

POST starts recording. DELETE stops recording. If already recording and a second POST is sent, the command is ignored and success returned. The same is true for sending a DELETE after recording has stopped.

#### Arguments
* **user** name of the user making the request
* **camera** identifer of the user's camera

#### Parameters
None

#### Response
On success, HTTP response with or without any information in the body. On error, HTTP response with error message in the body.

#### Example
	POST /feed/jimmy_james/office_camera
	DELETE /feed/jimmy_james/office_camera
