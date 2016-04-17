# LiveStream Microservice

## Description

This microservice returns a URL to retrieve the live MJPG stream from the camera.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information. (Example: export DEBUG=1)

## Connection
HTTPS on port 8080.

## API
    /LiveStream/v1?[args...]

### Version 1
    GET /LiveStream/v1?camera_id=[camera id]

#### Parameters
* **camera_id** ID of the camera from which to get live stream URL.

#### Response
On success, HTTP response with the live stream URL as the message body. On error, HTTP response with error message in the body.

#### Example

    GET /LiveStream/v1/camera_id=Camera123
