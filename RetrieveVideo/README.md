# RetrieveVideo Microservice

## Description

This microservice pulls recorded video and serves it to the user as a single
MP4 file.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information. (Example: export DEBUG=1)

## Connection
HTTPS on port 443, 8080 when in debug mode (gets around having to run it as root to open 443).

## API
    /video

## API Versions
### Version 1
    GET /video/{user}/{camera}?start=[start time]&stop=[stop time]

#### Arguments
* **user** name of the user making the request
* **camera** identifer of the user's camera

#### Parameters
* **start** Video start time.
* **stop**  Video stop time.

*Both times are in formats accepted by Javascript [Date.parse][parse]. These are treated as UTC time stamps. For user convenience it is best for the caller to appen the appropriate timezone correction to the timestamp. For example "-0600" for MDT when using the ISO 8601 format.*

#### Response
On success, HTTP response with MP4 embedded as file attachment. On error, HTTP response with error message in the body.

#### Example

    GET /video/user123/Camera2/start=2016-04-10T00:00:00&stop=2016-04-11T00:00:00.0

[parse]:https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/parse
