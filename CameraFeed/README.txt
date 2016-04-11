####################################
 README For CameraFeed Microservice
####################################

Description
===========
This microservice handles reading the live video stream from the camera.

API
===
/<version>?[args...]

API Versions
------------
+V1
 /v1?action=[start|stop]&camera_id=[id]

 action - whether to start or stop recording.
 camera_id - ID string of the camera from which to start or stop recording

 Example:
 /v1/action=start&camera_id=12345

