#!/bin/bash
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

set -e

# brew install ffmpeg --with-rtmp-dump
# brew install mp4box

# Sample Invokation:
#   ./shard.sh ./example/video/kelloggs.mp4 ./out/shard.hls.vod HLS_VOD
#   ./shard.sh ./example/video/kelloggs.mp4 ./out/shard.hls.live HLS_LIVE
#   ./shard.sh ./example/video/kelloggs.mp4 ./out/shard.dash.vod DASH_VOD
#   ./shard.sh ./example/video/kelloggs.mp4 ./out/shard.dash.live DASH_LIVE

# This script is to serve as an example of external transcoder for VOD assets that may be plumbed inside fakeOrigin as a VOD endpoint or static Directory of files


src=$1
destination=$2
format=$3
filename=$(basename "${destination}")

if [ $# -ne 3 ]; then
  echo 1>&2 "Usage: $0 <Sourcefile> <Destdirectory> <Type (HLS_VOD|HLS_LIVE|DASH_VOD|DASH_LIVE)>"
  exit 3
fi

mkdir -p $destination

set -x
# HLS VOD
if [[ $format == "HLS_VOD" ]]; then
ffmpeg -y -i "${src}" \
  -vf scale=w=640:h=360:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod  -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k -hls_segment_filename "${destination}/${filename}_seq_%03d_640x360-10000.ts" "${destination}/${filename}_640x360-10000.m3u8" \
  -vf scale=w=842:h=480:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 1400k -maxrate 1498k -bufsize 2100k -b:a 128k -hls_segment_filename "${destination}/${filename}_seq_%03d_842x480-20000.ts" "${destination}/${filename}_842x480-20000.m3u8" \
  -vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 2800k -maxrate 2996k -bufsize 4200k -b:a 128k -hls_segment_filename "${destination}/${filename}_seq_%03d_1280x720-40000.ts" "${destination}/${filename}_1280x720-40000.m3u8" \
  -vf scale=w=1920:h=1080:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 5000k -maxrate 5350k -bufsize 7500k -b:a 192k -hls_segment_filename "${destination}/${filename}_seq_%03d_1920x1080-80000.ts" "${destination}/${filename}_1920x1080-80000.m3u8"
fi

# HLS LIVE
if [[ $format == "HLS_LIVE" ]]; then
ffmpeg -y -i "${src}" \
  -vf scale=w=640:h=360:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k -hls_segment_filename "${destination}/${filename}_seq_%03d_640x360-10000.ts" "${destination}/${filename}_640x360-10000.m3u8" \
  -vf scale=w=842:h=480:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -b:v 1400k -maxrate 1498k -bufsize 2100k -b:a 128k -hls_segment_filename "${destination}/${filename}_seq_%03d_842x480-20000.ts" "${destination}/${filename}_842x480-20000.m3u8" \
  -vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -b:v 2800k -maxrate 2996k -bufsize 4200k -b:a 128k -hls_segment_filename "${destination}/${filename}_seq_%03d_1280x720-40000.ts" "${destination}/${filename}_1280x720-40000.m3u8" \
  -vf scale=w=1920:h=1080:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -b:v 5000k -maxrate 5350k -bufsize 7500k -b:a 192k -hls_segment_filename "${destination}/${filename}_seq_%03d_1920x1080-80000.ts" "${destination}/${filename}_1920x1080-80000.m3u8"
fi

# DASH Live
if [[ $format == "DASH_LIVE" ]]; then
ffmpeg -i "${src}" -c:a copy -vn "${destination}/${filename}-audio.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 5300k -maxrate 5300k -bufsize 2650k -vf 'scale=-1:1080' "${destination}/${filename}-1080.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 2400k -maxrate 2400k -bufsize 1200k -vf 'scale=-1:720' "${destination}/${filename}-720.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 1060k -maxrate 1060k -bufsize 530k -vf 'scale=-1:480' "${destination}/${filename}-480.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 600k -maxrate 600k -bufsize 300k -vf 'scale=-1:360' "${destination}/${filename}-360.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 260k -maxrate 260k -bufsize 130k -vf 'scale=-1:240' "${destination}/${filename}-240.mp4"

MP4Box -dash 1000 -rap -frag-rap -profile live -out "${destination}/$filename.mpd" "${destination}/${filename}-1080.mp4" "${destination}/${filename}-720.mp4" "${destination}/${filename}-480.mp4" "${destination}/${filename}-360.mp4" "${destination}/${filename}-240.mp4" "${destination}/${filename}-audio.mp4"
fi

# DASH VOD
if [[ $format == "DASH_VOD" ]]; then
ffmpeg -i "${src}" -c:a copy -vn "${destination}/${filename}-audio.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 5300k -maxrate 5300k -bufsize 2650k -vf 'scale=-1:1080' "${destination}/${filename}-1080.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 2400k -maxrate 2400k -bufsize 1200k -vf 'scale=-1:720' "${destination}/${filename}-720.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 1060k -maxrate 1060k -bufsize 530k -vf 'scale=-1:480' "${destination}/${filename}-480.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 600k -maxrate 600k -bufsize 300k -vf 'scale=-1:360' "${destination}/${filename}-360.mp4"
ffmpeg -i "${src}" -an -c:v libx264 -x264opts 'keyint=24:min-keyint=24:no-scenecut' -b:v 260k -maxrate 260k -bufsize 130k -vf 'scale=-1:240' "${destination}/${filename}-240.mp4"

MP4Box -dash 1000 -rap -frag-rap -profile onDemand -out "${destination}/$filename.mpd" "${destination}/${filename}-1080.mp4" "${destination}/${filename}-720.mp4" "${destination}/${filename}-480.mp4" "${destination}/${filename}-360.mp4" "${destination}/${filename}-240.mp4" "${destination}/${filename}-audio.mp4"
fi
