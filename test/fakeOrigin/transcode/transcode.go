package transcode

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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/endpoint"
)

type TranscodeMeta struct {
	CrcHash         string   `json:"crc_hash"`
	Cmd             string   `json:"cmd"`
	Args            []string `json:"args"`
	LastTranscodeDT string   `json:"last_transcode_datetime"`
}

func Do(ep *endpoint.Endpoint, cmdStr string, args []string) error {
	fmt.Printf("Creating directory: %s\n", ep.OutputDirectory)
	err := os.MkdirAll(ep.OutputDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println("Directory creation complete")

	isMatch := false
	isMatch, err = checkMeta((*ep), cmdStr, args)
	if err != nil {
		return err
	}
	if !isMatch {
		fmt.Println("Transcode Meta Check Failed")
		fmt.Println("Beginning Transcode")
		err = RunSynchronousCmd(cmdStr, args)
		fmt.Println("Transcode complete")
		if err != nil {
			return err
		}
		err = generateMasterManifest(ep)
		if err != nil {
			return err
		}
		err = generateMeta((*ep), cmdStr, args)
	} else {
		fmt.Println("Transcode Meta Check Matched")
	}

	return err
}

func RunSynchronousCmd(baseCmd string, args []string) error {
	thisDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return errors.New("getting absolute filepath '" + os.Args[0] + "': " + err.Error())
	}
	cmd := exec.Command(baseCmd, args...)
	cmd.Dir = thisDir
	fmt.Printf("Transcode Cmd: %#v\n", baseCmd)
	fmt.Printf("Transcode Args: %#v\n", args)
	fmt.Println("Working Directory:", cmd.Dir)
	writer := io.MultiWriter(os.Stdout)
	cmd.Stderr = writer
	cmd.Stdout = writer
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("running baseCmd '%+v' args '%+v': %+v", baseCmd, args, err)
	}
	return nil
}

func hashFileCrc32(filePath string, polynomial uint32) (string, error) {
	out := ""

	file, err := os.Open(filePath)
	if err != nil {
		return out, err
	}
	defer file.Close()

	tp := crc32.MakeTable(polynomial)
	hash := crc32.New(tp)

	if _, err := io.Copy(hash, file); err != nil {
		return out, err
	}

	hashInBytes := hash.Sum(nil)[:]
	out = hex.EncodeToString(hashInBytes)

	return out, nil
}

func generateMeta(ep endpoint.Endpoint, cmdStr string, args []string) error {
	// CRC-32 reversed polynomial https://en.wikipedia.org/wiki/Cyclic_redundancy_check
	crc, err := hashFileCrc32(ep.Source, 0xEDB88320)
	if err != nil {
		return err
	}

	metainfo := TranscodeMeta{
		Cmd:             cmdStr,
		Args:            args,
		CrcHash:         crc,
		LastTranscodeDT: time.Now().Format(time.RFC1123),
	}

	bytes, err := json.MarshalIndent(metainfo, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(ep.OutputDirectory+"/"+ep.DiskID+".meta.json", bytes, 0755)

	return err
}

// GetMeta retrieves the disk metadata for the given endpoint
func GetMeta(ep endpoint.Endpoint) (TranscodeMeta, error) {
	raw, err := ioutil.ReadFile(ep.OutputDirectory + "/" + ep.DiskID + ".meta.json")
	if err != nil {
		raw = []byte("{}")
	}
	var sourcemetainfo TranscodeMeta
	if err = json.Unmarshal(raw, &sourcemetainfo); err != nil {
		return TranscodeMeta{}, err
	}
	return sourcemetainfo, nil
}

func checkMeta(ep endpoint.Endpoint, cmdStr string, args []string) (bool, error) {
	// CRC-32 reversed polynomial https://en.wikipedia.org/wiki/Cyclic_redundancy_check
	crc, err := hashFileCrc32(ep.Source, 0xEDB88320)
	if err != nil {
		return false, err
	}

	destmetainfo := TranscodeMeta{
		Cmd:     cmdStr,
		Args:    args,
		CrcHash: crc,
	}

	sourcemetainfo, err := GetMeta(ep)
	if err != nil {
		return false, err
	}

	if sourcemetainfo.CrcHash == destmetainfo.CrcHash && sourcemetainfo.Cmd == destmetainfo.Cmd && reflect.DeepEqual(sourcemetainfo.Args, destmetainfo.Args) {
		return true, nil
	}

	return false, nil
}

func generateMasterManifest(ep *endpoint.Endpoint) error {
	files, err := ioutil.ReadDir(ep.OutputDirectory)
	if err != nil {
		return err
	}
	type manifest struct {
		resolution string
		bandwidth  string
		name       string
	}
	var manifests []manifest
	for _, file := range files {
		var r *regexp.Regexp
		r, err = regexp.Compile(ep.DiskID + `.*?(?:(?P<res>\d+x\d+)-(?P<bw>\d+))?\.m3u8`)
		if err != nil {
			return err
		}
		if r.MatchString(file.Name()) {
			match := r.FindStringSubmatch(file.Name())
			ep.ABRManifests = append(ep.ABRManifests, file.Name())
			manifests = append(manifests, manifest{
				resolution: match[1],
				bandwidth:  match[2],
				name:       file.Name(),
			})
		}
	}
	if len(manifests) == 0 {
		return errors.New("Master manifest detection failed")
	} else if len(manifests) == 1 {
		return nil
	}

	out := "#EXTM3U\n#EXT-X-VERSION:3\n"
	for _, layer := range manifests {
		fmt.Printf("DEBUG: %+v\n", layer)
		out = out + "#EXT-X-STREAM-INF:BANDWIDTH=" + layer.bandwidth + ",RESOLUTION=" + layer.resolution + "\n" + layer.name + "\n"
	}
	err = ioutil.WriteFile(ep.OutputDirectory+"/"+ep.DiskID+".m3u8", []byte(out), 0755)
	return err
}
