######################################
 README For RetrieveVideo Microservice
######################################

Description
===========
This microservice pulls recorded video and serves it to the user as a single
MP4 file.

API
===
/video?[args...]

API Versions
------------
+V1
 /video?start=[start time]&stop=[stop time]&camera_id=[camera id]

 start Video start time. Format is one accepted by Javascript Date.parse.
 stop  Video stop time. Format is one accepted by Javascript Date.parse.

 Example:
 /start=2016-04-10T00:00:00&stop=2016-04-11T00:00:00.0&camera_id=Camera123
