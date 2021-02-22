<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->
# Endpoint Examples
Some of these examples make use of [ffmpeg](https://www.ffmpeg.org/), however any script or application that can satisfy the required tokens is sufficient if transcoding is needed.  The specific commands used in these examples may vary based on your individual requirements.  Refer to the documentation of your transcode command for more information.

## HLS VOD with Adaptive Bit Rate (ABR)
```json
{
  "id": "HLS_ABR_VOD",
  "source": "./example/video/kelloggs.mp4",
  "outputdir": "./out/HLS_ABR_VOD",
  "type": "vod",
  "manual_command": [
    "ffmpeg",
    "-y",
    "-i", "%SOURCE%",
    "-vf", "scale=w=1280:h=720:force_original_aspect_ratio=decrease",
    "-c:a", "aac",
    "-ar", "48000",
    "-c:v", "h264",
    "-profile:v", "main",
    "-crf", "20",
    "-sc_threshold", "0",
    "-g", "48",
    "-keyint_min", "48",
    "-hls_time", "4",
    "-hls_playlist_type", "vod",
    "-b:v", "2800k",
    "-b:a", "128k",
    "-hls_list_size", "0",
    "-hls_segment_filename", "%OUTPUTDIRECTORY%/%DISKID%_seq_%03d_1280x720.ts",
    "%OUTPUTDIRECTORY%/%DISKID%_1280x720-20000.m3u8",
    "-vf", "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
    "-c:a", "aac",
    "-ar", "48000",
    "-c:v", "h264",
    "-profile:v", "main",
    "-crf", "20",
    "-sc_threshold", "0",
    "-g", "48",
    "-keyint_min", "48",
    "-hls_time", "4",
    "-hls_playlist_type", "vod",
    "-b:v", "5000k",
    "-b:a", "192k",
    "-hls_list_size", "0",
    "-hls_segment_filename", "%OUTPUTDIRECTORY%/%DISKID%_seq_%03d_1920x1080.ts",
    "%OUTPUTDIRECTORY%/%DISKID%_1920x1080-40000.m3u8"
  ]
}

```
## Duplicate Endpoint
```json
{
  "id": "DUPLICATE_HLS_ABR_VOD",
  "override_disk_id": "HLS_ABR_VOD",
  "source": "./example/video/kelloggs.mp4",
  "outputdir": "./out/HLS_ABR_VOD",
  "type": "vod"
}
```
The important bits to notice here is that the `override_disk_id`, `outputdir`, & `type` match our first example.  Only the `id` and `default_headers` are different.  This will cause the transcoder metadata to match and therefore inform fakeOrigin that it shouldn't bother trying to transcode it again.

## HLS with ABR live manifests
```json
{
  "id": "HLS_ABR_LIVE",
  "source": "./example/video/kelloggs.mp4",
  "outputdir": "./out/HLS_ABR_LIVE",
  "type": "live",
  "manual_command": [
    "ffmpeg",
    "-y",
    "-i", "%SOURCE%",
    "-vf", "scale=w=1280:h=720:force_original_aspect_ratio=decrease",
    "-c:a", "aac",
    "-ar", "48000",
    "-c:v", "h264",
    "-profile:v", "main",
    "-crf", "20",
    "-sc_threshold", "0",
    "-g", "48",
    "-keyint_min", "48",
    "-hls_time", "4",
    "-hls_playlist_type", "vod",
    "-b:v", "2800k",
    "-b:a", "128k",
    "-hls_list_size", "0",
    "-hls_segment_filename", "%OUTPUTDIRECTORY%/%DISKID%_seq_%03d_1280x720.ts",
    "%OUTPUTDIRECTORY%/%DISKID%_1280x720-20000.m3u8",
    "-vf", "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
    "-c:a", "aac",
    "-ar", "48000",
    "-c:v", "h264",
    "-profile:v", "main",
    "-crf", "20",
    "-sc_threshold", "0",
    "-g", "48",
    "-keyint_min", "48",
    "-hls_time", "4",
    "-hls_playlist_type", "vod",
    "-b:v", "5000k",
    "-b:a", "192k",
    "-hls_list_size", "0",
    "-hls_segment_filename", "%OUTPUTDIRECTORY%/%DISKID%_seq_%03d_1920x1080.ts",
    "%OUTPUTDIRECTORY%/%DISKID%_1920x1080-40000.m3u8"
  ]
}
```
On live endpoint types, we intercept the m3u8 requests and rewrite them so that they loop.
## HLS without ABR VOD
```json
{
  "id": "HLS_VOD",
  "source": "./example/video/kelloggs.mp4",
  "outputdir": "./out/HLS_VOD",
  "type": "vod",
  "manual_command": [
    "ffmpeg",
    "-y",
    "-i", "%SOURCE%",
    "-vf", "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
    "-c:a", "aac",
    "-ar", "48000",
    "-c:v", "h264",
    "-profile:v", "main",
    "-crf", "20",
    "-sc_threshold", "0",
    "-g", "48",
    "-keyint_min", "48",
    "-hls_time", "4",
    "-hls_playlist_type", "vod",
    "-b:v", "5000k",
    "-b:a", "192k",
    "-hls_list_size", "0",
    "-hls_segment_filename", "%OUTPUTDIRECTORY%/%DISKID%_seq_%03d.ts",
    "%OUTPUTDIRECTORY%/%DISKID%.m3u8"
  ]
}
```

## HLS Event
```json
{
  "id": "HLS_EVENT",
  "source": "./example/video/kelloggs.mp4",
  "outputdir": "./out/HLS_EVENT",
  "type": "event",
  "manual_command": [
    "ffmpeg",
    "-y",
    "-i", "%SOURCE%",
    "-vf", "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
    "-c:a", "aac",
    "-ar", "48000",
    "-c:v", "h264",
    "-profile:v", "main",
    "-crf", "20",
    "-sc_threshold", "0",
    "-g", "48",
    "-keyint_min", "48",
    "-hls_time", "4",
    "-hls_playlist_type", "event",
    "-b:v", "5000k",
    "-b:a", "192k",
    "-hls_list_size", "0",
    "-hls_segment_filename", "%OUTPUTDIRECTORY%/%DISKID%_seq_%03d.ts",
    "%OUTPUTDIRECTORY%/%DISKID%.m3u8"
  ]
}
```
## Static file
```json
{
  "id": "SampleVideo",
  "source": "./example/video/kelloggs.mp4",
  "outputdir": "./out",
  "type": "static"
}
```
This shows how to serve single files with fakeOrigin.  Also, `id` must still be unique even if they are different source files.

## Directory
```json
{
  "id": "SampleDir",
  "source": "./example/video",
  "outputdir": "./out",
  "type": "dir"
}
```
This shows how to serve all files in a given directory recursively.  Also, `id` must still be unique even if they are different source files.  This type serves each file it finds at startup as a static file.

## Player Troubleshooting
If you're running into issues with a Javascript-based test player, there is a good chance you may need to get add some default CORS headers to your endpoint config.
```json
"default_headers": {
  "Access-Control-Allow-Headers": [
    "*"
  ],
  "Access-Control-Allow-Origin": [
    "*"
  ]
}
```
This is a wide open example and should be tailored to the domains of your test player appropriately.
