######################################
 README For CameraControl Microservice
######################################

Description
===========
This microservice sends PTZ (pan/tilt/zoom) commands to the camera.

Set the 'DEBUG' environment variable to any non-empty value for debugging
information.

API
===
/ptz/[args...]

API Versions
------------
+V1
 /ptz/v1?action=[start,stop]&direction=[Up,Down,Left,Right]&velocity[1-8]&camera_i=[camera id]

 action     starts or stops the camera moving
 direction  direction in which the camera is to move
 velocity   how quickly the camera should move
 camera_id  which camera to move

 On stop commands, direction and velocity are required but ignored.

 Example:
 /v1?action=start&direction=Up&velocity=5&camera_id=12345
 /v1?action=stop&direction=Up&velocity=5&camera_id=12345
