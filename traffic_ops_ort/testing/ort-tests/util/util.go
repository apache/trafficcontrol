package util

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"errors"
	"os"
	"strings"
)

func readFile(fileName string) ([]byte, error) {
	if fileName == "" {
		return nil, errors.New("filename is empty.")
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	size := fileInfo.Size()

	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fileData := make([]byte, size)
	cnt, err := fd.Read(fileData)
	if err != nil || int64(cnt) != size {
		return nil, errors.New("unable to completely read from '" + fileName + "' : " + err.Error())
	}

	return fileData, nil
}

func filterOutComments(body []string) []string {

	var newlines []string

	newlines = make([]string, 0)

	for ii := range body {
		line := body[ii]
		if strings.HasPrefix(line, "#") {
			continue
		}
		newlines = append(newlines, line)
	}

	return newlines
}

func DiffFiles(firstFile string, secondFile string) (bool, error) {
	d1, err := readFile(firstFile)
	if err != nil {
		return false, err
	}
	d2, err := readFile(secondFile)
	if err != nil {
		return false, err
	}

	d1Data := strings.Split(string(d1), "\n")
	str1 := strings.Join(filterOutComments(d1Data), "\n")

	d2Data := strings.Split(string(d2), "\n")
	str2 := strings.Join(filterOutComments(d2Data), "\n")

	if str1 == str2 {
		return true, nil
	} else {
		return false, nil
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
