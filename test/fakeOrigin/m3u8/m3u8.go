package m3u8

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type M3U8 struct {
	Version           string
	HasVersion        bool
	TargetDuration    time.Duration
	HasTargetDuration bool
	MediaSequence     int64
	HasMediaSequence  bool
	TSes              []M3U8TS
	VARs              []M3U8VAR
}

type M3U8TS struct {
	Duration      time.Duration
	Discontinuity bool
	Filepath      string
}

type M3U8VAR struct {
	Bandwidth  int
	Resolution string
	Filepath   string
}

func Parse(fileBytes []byte) (M3U8, error) {
	lines := bytes.Split(fileBytes, []byte("\n"))

	// if bytes.TrimSpace(lines[0]) != []byte("#EXTM3U") {
	// 	fmt.Println("Error parsing m3u8: file must begin with '#EXTM3U'")
	// 	return
	// }

	m3u8 := M3U8{}

	inTarget := false
	ts := M3U8TS{}
	variantStream := M3U8VAR{}
linesLoop:
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if !inTarget {
			switch tagType(line) {
			case TagTypeHeader:
				continue
			case TagTypeVersion:
				m3u8.Version = string(bytes.TrimSpace(line[len(VersionPrefix):]))
				m3u8.HasVersion = true
				continue
			case TagTypeTargetDuration:
				targetDurationStr := string(bytes.TrimSpace(line[len(TargetDurationPrefix):]))
				targetDurationSeconds, err := strconv.ParseFloat(targetDurationStr, 64)
				if err != nil {
					return M3U8{}, errors.New("parsing m3u8: " + TargetDurationPrefix + " not a number")
				}
				m3u8.TargetDuration = time.Duration(targetDurationSeconds * float64(time.Second))
				m3u8.HasTargetDuration = true
				continue
			case TagTypeMediaSequence:
				mediaSequenceStr := string(bytes.TrimSpace(line[len(MediaSequencePrefix):]))
				mediaSequence, err := strconv.ParseInt(mediaSequenceStr, 10, 64)
				if err != nil {
					return M3U8{}, errors.New("parsing m3u8: " + MediaSequencePrefix + " not a number")
				}
				m3u8.MediaSequence = mediaSequence
				m3u8.HasMediaSequence = true
				continue
			case TagTypeTSInfo:
				fields := bytes.Split(line[len(TSInfoPrefix):], []byte(","))
				durationStr := bytes.TrimSpace(fields[0])
				if len(durationStr) == 0 {
					return M3U8{}, errors.New("parsing m3u8: malformed line '" + string(line))
				}
				durationSeconds, err := strconv.ParseFloat(string(durationStr), 64)
				if err != nil {
					return M3U8{}, errors.New("parsing m3u8: line '" + string(line) + "' duration '" + string(durationStr) + "' not a number")
				}
				ts.Duration = time.Duration(durationSeconds * float64(time.Second))
				inTarget = true
				continue
			case TagTypeVariantStream:
				fields := bytes.Split(line[len(VariantPrefix):], []byte(","))
				variantArgs := make(map[string]string)
				for _, kvp := range fields {
					temp := bytes.Split(kvp, []byte("="))
					k := string(temp[0])
					v := string(temp[1])
					k = strings.ToLower(k)
					v = strings.ToLower(v)
					if string(temp[0]) != "" {
						variantArgs[strings.ToLower(string(temp[0]))] = strings.ToLower(string(temp[1]))
					}
				}
				if _, ok := variantArgs["bandwidth"]; ok {
					var strconverror error
					variantStream.Bandwidth, strconverror = strconv.Atoi(variantArgs["bandwidth"])
					if strconverror != nil {
						return M3U8{}, errors.New("parsing m3u8: variant stream bandwidth malformed: " + variantArgs["bandwidth"] + " error: " + strconverror.Error())
					}
				} else {
					return M3U8{}, errors.New("parsing m3u8: variant stream bandwidth missing '" + string(line) + "'")
				}
				variantStream.Resolution = variantArgs["resolution"]
				if variantStream.Resolution == "" {
					return M3U8{}, errors.New("parsing m3u8: variant stream resolution missing '" + string(line) + "'")
				}

				inTarget = true
				continue
			case TagTypeEndList:
				break linesLoop
			// TODO handle #EXT-X-DISCONTINUITY
			default:
				if line[0] == '#' {
					fmt.Println("Warning: parsing m3u8: line '" + string(line) + "' unknown directive")
					continue
				}
				return M3U8{}, errors.New("parsing m3u8: unknown line '" + string(line))
			}
		} else {
			switch targetType(line) {
			case TargetTypeTS:
				ts.Filepath = string(line)
				m3u8.TSes = append(m3u8.TSes, ts)
				inTarget = false
				ts = M3U8TS{}
				variantStream = M3U8VAR{}
			case TargetTypeManifest:
				variantStream.Filepath = string(line)
				m3u8.VARs = append(m3u8.VARs, variantStream)
				inTarget = false
				ts = M3U8TS{}
				variantStream = M3U8VAR{}
			default:
				return M3U8{}, errors.New("parsing m3u8: unknown target line '" + string(line))
			}
		}
	}
	return m3u8, nil
}

type TagType int

const (
	TagTypeInvalid TagType = iota
	TagTypeHeader
	TagTypeVersion
	TagTypeTargetDuration
	TagTypeMediaSequence
	TagTypeTSInfo
	TagTypeEndList
	TagTypeVariantStream
)

type TargetType int

const (
	TargetTypeInvalid TargetType = iota
	TargetTypeManifest
	TargetTypeTS
)

const HeaderPrefix = "#EXTM3U"
const VersionPrefix = "#EXT-X-VERSION:"
const TargetDurationPrefix = "#EXT-X-TARGETDURATION:"
const MediaSequencePrefix = "#EXT-X-MEDIA-SEQUENCE:"
const TSInfoPrefix = "#EXTINF:"
const EndListPrefix = "#EXT-X-ENDLIST"
const VariantPrefix = "#EXT-X-STREAM-INF:"
const TSSuffix = ".ts"
const ManifestSuffix = ".m3u8"

func tagType(line []byte) TagType {
	switch {
	case bytes.HasPrefix(line, []byte(HeaderPrefix)):
		return TagTypeHeader
	case bytes.HasPrefix(line, []byte(VersionPrefix)):
		return TagTypeVersion
	case bytes.HasPrefix(line, []byte(TargetDurationPrefix)):
		return TagTypeTargetDuration
	case bytes.HasPrefix(line, []byte(MediaSequencePrefix)):
		return TagTypeMediaSequence
	case bytes.HasPrefix(line, []byte(TSInfoPrefix)):
		return TagTypeTSInfo
	case bytes.HasPrefix(line, []byte(EndListPrefix)):
		return TagTypeEndList
	case bytes.HasPrefix(line, []byte(VariantPrefix)):
		return TagTypeVariantStream
	}
	return TagTypeInvalid
}

func targetType(line []byte) TargetType {
	switch {
	case bytes.HasSuffix(line, []byte(TSSuffix)):
		return TargetTypeTS
	case bytes.HasSuffix(line, []byte(ManifestSuffix)):
		return TargetTypeManifest
	}
	return TargetTypeInvalid
}

// TransformVodToLive creates a live m3u8 from a given VOD m3u8.
// The created m3u8 will have TS files, of the lesser of MinLiveFiles and MinLiveDuration
func TransformVodToLive(vod M3U8, offset time.Duration, minFiles int64, minDuration time.Duration) (M3U8, error) {
	if !vod.HasVersion {
		return M3U8{}, errors.New("vod must have version") // TODO default version
	}
	if !vod.HasTargetDuration {
		return M3U8{}, errors.New("vod must have taget duration") // TODO compute target duration
	}
	if vod.HasMediaSequence && vod.MediaSequence != 0 {
		return M3U8{}, errors.New("vod must not have media sequence")
	}
	if len(vod.TSes) < 1 {
		return M3U8{}, errors.New("vod must have at least 1 TS")
	}

	live := M3U8{}
	live.Version = vod.Version
	live.TargetDuration = vod.TargetDuration
	live.MediaSequence = GetLiveMediaSequence(vod, offset)

	totalDuration := time.Duration(0)
	totalFiles := int64(0)
	for i := live.MediaSequence; totalFiles < minFiles && totalDuration < minDuration; i++ {
		vodIdx := i % int64(len(vod.TSes))
		ts := vod.TSes[vodIdx]
		ts.Discontinuity = vodIdx == 0
		live.TSes = append(live.TSes, ts)
		totalDuration += ts.Duration
		totalFiles++
	}
	return live, nil
}

// GetLiveMediaSequence gets the live media sequence, from the vod m3u8.
// That is, if vod has 3 TSes each 2 seconds long, and offset is 3s, this returns 2, because the 3rd second occurs within the 2nd TS (which starts at second 2 and is 2 seconds long, i.e. contains absolute seconds 2-4 of the video).
// The returned MediaSequence should be modded on the length of the TSes, to get the proper TSes to serve.
// Each time the mod is 0, the 0 TS should also insert a #EXT-X-DISCONTINUITY into the manifest, via M3U8TS.Discontinuity, in order to play correctly, assuming the TS files have correct continuity counters.
func GetLiveMediaSequence(vod M3U8, offset time.Duration) int64 {
	if len(vod.TSes) == 0 {
		return 0
	}
	i := int64(0)
	for sum := vod.TSes[0].Duration; sum <= offset; i++ {
		nextTS := vod.TSes[(i+1)%int64(len(vod.TSes))]
		sum += nextTS.Duration
	}
	return i
}

// SerializeLiveM3U8 serialized the given M3U8 object as a live m3u8 manifest file.
// It assumes the given m3u8 is valid, and does not check validity or existence of fields. Callers should check validity before calling (a M3U8 returned by TransformVodToLive without error is guaranteed valid).
func SerializeLive(m3u8 M3U8) []byte {
	b := []byte(`#EXTM3U
#EXT-X-VERSION:` + m3u8.Version + `
#EXT-X-TARGETDURATION:` + strconv.FormatFloat(float64(m3u8.TargetDuration)/float64(time.Second), 'f', -1, 64) + `
#EXT-X-MEDIA-SEQUENCE:` + strconv.FormatInt(m3u8.MediaSequence, 10) + `
`)
	for _, ts := range m3u8.TSes {
		if ts.Discontinuity {
			b = append(b, []byte(`#EXT-X-DISCONTINUITY
`)...)
		}
		b = append(b, []byte(`#EXTINF:`+strconv.FormatFloat(float64(ts.Duration)/float64(time.Second), 'f', 6, 64)+`,
`)...)
		b = append(b, []byte(ts.Filepath+`
`)...)
	}
	return b
}

// GetTotalTime returns the sum of the durations of all TSes in the given vod M3U8.
func GetTotalTime(vod M3U8) time.Duration {
	sum := time.Duration(0)
	for _, ts := range vod.TSes {
		sum += ts.Duration
	}
	return sum
}

func LoadTSes(m3u8 M3U8, m3u8Path string) (map[string][]byte, error) {
	tses := map[string][]byte{}
	dir := filepath.Dir(m3u8Path)
	for _, ts := range m3u8.TSes {
		path := filepath.Join(dir, ts.Filepath)
		bts, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.New("loading '" + path + "':" + err.Error())
		}
		tses[ts.Filepath] = bts
	}
	return tses, nil
}
