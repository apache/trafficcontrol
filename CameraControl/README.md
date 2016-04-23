# CameraControl Microservice

## Description

This microservice sends PTZ (pan/tilt/zoom) commands to the camera. Camera movements are started or stopped. They are not performed in discrete motions. Sending a start command will cause the camera to move in that direction until a stop command is sent.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information.

## Connection
HTTPS on port 443, 8080 when in debug mode (gets around having to run it as root to open 443).

## API
    /control

## API Versions
### Version 1
    POST /control/{user}/{camera}?action=[start,stop]&direction=[Up,Down,Left,Right]&velocity[1-8]

#### Arguments
* **user** name of the user making the request
* **camera** identifer of the user's camera

#### Parameters
* **action**     starts or stops the camera moving
* **direction**  direction in which the camera is to move
* **velocity**   how quickly the camera should move

On stop commands, direction and velocity are required, and must be valid values, but are ultimately ignored. They do not have to correspond to the values from the start command.

#### Response
On success, HTTP response with or without any information in the body. On error, HTTP response with error message in the body.

#### Example

    POST /control/rhonda/LivingRoom?action=start&direction=Up&velocity=5
	POST /control/user123/1?action=stop&direction=Up&velocity=5
