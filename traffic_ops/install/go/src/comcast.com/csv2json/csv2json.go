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

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Profile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Parameter struct {
	Name        string `json:"name"`
	Config_File string `json:"config_file"`
	Value       string `json:"value"`
}

type ProfileParameter struct {
	Profile    string `json:"profile"`
	Parameter  string `json:"parameter"`
	ConfigFile string `json:"config_file"`
	Value      string `json:"value"`
}

func main() {
	createProfileData()
	createParameterData()
	createProfileParameterData()
}

func writeToFile(message []byte, fileName string) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := string(message) + "\n"
	n, err := io.WriteString(f, s)
	if err != nil {
		fmt.Println(n, err)
		return
	}
	f.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//TODO:  break out reading of csv file into it's own func

func createProfileData() {
	fileName := "/opt/traffic_ops/install/data/csv/profile.csv"
	fmt.Println("converting profile data json csv")
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	lineCount := 0
	for {
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("name", record[0], "description", record[1])
		//convert to json
		p := Profile{record[0], record[1]}
		b, err := json.Marshal(p)
		//write json file
		writeToFile(b, "/opt/traffic_ops/install/data/json/profile.json")
		lineCount += 1
	}
}

func createParameterData() {
	fileName := "/opt/traffic_ops/install/data/csv/parameter.csv"
	fmt.Println("converting parameter data json csv")
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	lineCount := 0
	for {
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("name=", record[0], "config_file=", record[1], "value=", record[2])
		//convert to json
		p := Parameter{record[0], record[1], record[2]}
		b, err := json.Marshal(p)
		//write json file
		writeToFile(b, "/opt/traffic_ops/install/data/json/parameter.json")
		lineCount += 1
	}
}

func createProfileParameterData() {
	fileName := "/opt/traffic_ops/install/data/csv/profile_parameter.csv"
	fmt.Println("converting profile_parameter data json csv")
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	lineCount := 0
	for {
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("profile=", record[0], "parameter=", record[1], "config_file=", record[2], "value", record[3])
		//convert to json
		p := ProfileParameter{record[0], record[1], record[2], record[3]}
		b, err := json.Marshal(p)
		//write json file
		writeToFile(b, "/opt/traffic_ops/install/data/json/profile_parameter.json")
		lineCount += 1
	}
}
