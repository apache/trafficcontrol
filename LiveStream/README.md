# LiveStream Microservice

## Description

This microservice returns a URL to retrieve the live MJPG stream from the camera.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information. (Example: export DEBUG=1)

## Connection
HTTPS on port 443, 8080 when in debug mode (gets around having to run it as root to open 443).

## API
    /livestream

## API Versions
### Version 1
    GET /livestream/{user}/{camera}

#### Arguments
* **user** name of the user making the request
* **camera** identifer of the user's camera

#### Response
On success, HTTP response with the live stream URL as the message body. On error, HTTP response with error message in the body.

#### Example

    GET /livestream/tony/OfficeCamera
