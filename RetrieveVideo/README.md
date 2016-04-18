# RetrieveVideo Microservice

## Description

This microservice pulls recorded video and serves it to the user as a single
MP4 file.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information. (Example: export DEBUG=1)

## Connection
HTTPS on port 8080.

## API
    /video/v1?[args...]

## API Versions
### Version 1
    GET /video/v1?start=[start time]&stop=[stop time]&camera_id=[camera id]

#### Parameters
* **start** Video start time. Format is one accepted by Javascript [Date.parse][parse].
* **stop**  Video stop time. Format is one accepted by Javascript [Date.parse][parse].
* **camera_id**  Id of the camera from which to retrieve recorded data.

#### Response
On success, HTTP response with MP4 embedded as file attachment. On error, HTTP response with error message in the body.

#### Example

    GET /video/v1/start=2016-04-10T00:00:00&stop=2016-04-11T00:00:00.0&camera_id=Camera123

[parse]:https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/parse
