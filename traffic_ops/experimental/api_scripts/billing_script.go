package main

// Summary:
//to run -> go run billing_script.go -startDate="04/01/2017" -endDate="05/01/2017" -xmlidFilter="" -longDesc1Filter="" -sumByDay=""

// This program is intended to get usage data (Sum Gbs and 95th Percentile MBs) from various services
// via this api - https://cdnportal.comcast.net/api/1.1/deliveryservices/999/server_types/edge/metric_types/kbps/start_date/1453223126/end_date/1453309526
// The usage date will be printed on the screen and will also be stored in a csv file
//
// The program does the following things:
// 1. Parses various user inputs such as start date, end date, xmlidFilter, longDesc1Filter, sumByDay
// 2. Reads user credentials from a file
// 3. Makes an API request using #2 and gets a mojo key
// 4. The mojo key is used to authenticate all other API requests
// 5. Makes an API request to get all the assigned services using the mojo key
// 6. Makes an API request to get usage data for all the obtained services from #5
// 7. Prints the usage data on the screen
// 8. Stores the usage data in a csv file

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Define command line variables
var startDate string
var endDate string
var xmlidFilter string     //service name filter
var longDesc1Filter string //longdesc1=cts; prod filter to get services in prod
var sumByDay string        //filter to get daily breakdown of the usage data

//------------------------------------------------------------------------------
//Initialize command line variables
func init() {
	flag.StringVar(&startDate, "startDate", "01/01/2017", "Start Date")
	flag.StringVar(&endDate, "endDate", "02/01/2017", "End Date")
	flag.StringVar(&xmlidFilter, "xmlidFilter", "", "service name filter")
	flag.StringVar(&longDesc1Filter, "longDesc1Filter", "", "longdesc1 ='cts; prod' aka services in production filter")
	flag.StringVar(&sumByDay, "sumByDay", "", "sumByDay filter")

}

//------------------------------------------------------------------------------

//main function:
// 1. parse the command line inputs
// 2. get the user creds
// 3. get the mojo key
// 4. call the first API to get assigned services based on those creds
// 5. call the second API to get usage data for those services and send it to a channel
// 6. send the channel response to a receiver function that parses the channel response and returns a slice of slice string
// 7. get the return value from receiver function and write the output to screen and csv file

func main() {

	//initialize wait group
	wg := &sync.WaitGroup{}

	//create an empty channel of HTTPResponse stuct
	httpCh := make(chan HTTPResponse)

	//Parse cmd flags
	flag.Parse()

	//GET CREDENTIALS FROM FILE IF NOT ENTERED ON COMMAND LINE
	fmt.Println("Getting user credentials from file...")
	user, pw := read_creds_file()

	//GET MOJOLICIOUS KEY FROM API
	fmt.Println("Sending API call to get mojolicious key...")
	mojo_key := api_request(user, pw)

	//GET Assigned Services for the user
	fmt.Println("Getting assigned services from API...")
	api_req_output := assigned_services(mojo_key)

	// range over the assigned services, get usage data concurrently and store it in a channel
	for xmlID, dsID := range api_req_output {
		wg.Add(1) //increment the waitgroup counter
		go services_api_req(mojo_key, xmlID, dsID, httpCh, wg)
	}

	//closing channel
	go func() {
		wg.Wait()     //wait for all goroutine to finish
		close(httpCh) //close the channel
	}()

	fmt.Println("Getting usage data for each assigned service ...")
	output := receiver(httpCh)

	// WRITE OUTPUT TO SCREEN
	fmt.Println("Writing output to screen...")
	write_screen(output)
	//

	//WRITE OUTPUT TO FILE
	fmt.Println("Writing output to file...")
	write_file(output)
}

//------------------------------------------------------------------------------

//receive channel response and return a slice of string slice
func receiver(httpCh chan HTTPResponse) (data [][]string) {

	//Define Output csv and add header row
	data = append(data, []string{
		"Delivery Service",
		"dsid",
		"start date",
		"end date",
		"sum GBs",
		"95th Percentile MBs",
	})

	for result := range httpCh { //loop through the channel response
		// fmt.Println("received:", result.xmlId, result.dsId, result.startdate, result.enddate, result.sum, result.Nine5th)

		//add data to the output csv slice
		data = append(data, []string{
			result.xmlId,
			result.dsId,
			result.startdate,
			result.enddate,
			result.sum,
			result.nine5th,
		})

	}

	return

}

//------------------------------------------------------------------------------

//ERROR HANDLER FOR ALL FUNCTIONS -------------------------------------------------------------
func checkError(message string, err error) { //pass in an err to check and a string to print in case it fails
	if err != nil {
		log.Fatal(message, err)
	}
}

//------------------------------------------------------------------------------

//takes a string date(01/01/2017) and converts it to unix timestamp(1483228800)
func convertToUnixTimeStamps(userDate string) (unixDate string) {
	layout := "01/02/2006"
	t, err := time.Parse(layout, userDate)
	checkError("Failed time conversion: ", err)
	unixDate = strconv.FormatInt(t.Unix(), 10)
	return
}

//------------------------------------------------------------------------------

//READ IN CREDENTIALS FROM FILE--------------------------------------------------------
func read_creds_file() (user, pw string) {
	path := os.Getenv("HOME") + "/Downloads/scripts/"
	file, err := os.Open(path + "CDN_API_Credentials.txt")
	checkError("Failed to open creds file", err)
	contents, err := ioutil.ReadAll(file) //read file to bytes
	checkError("Failed to read creds file", err)
	creds := strings.Split(string(contents), "\n") //create a slice of the contents, must convert from bytes to string
	file.Close()                                   //closes file
	user = creds[0]                                // extract user from file strings
	pw = creds[1]

	return
}

//------------------------------------------------------------------------------

//API REQUEST FOR MOJOLICIOUS KEY REQUIRED FOR ALL API DATA REQUESTS
func api_request(user, pw string) (mojo_key string) {
	user_pw := map[string]string{"p": pw, "u": user} // create a map of user and pw for inclusion in request
	user_pw_json, _ := json.Marshal(user_pw)         // convert user/password map to json structure for inclusion in request
	key_url := "https://cdnportal.comcast.net/api/1.2/user/login"
	// add key to url
	key_req, err := http.NewRequest("POST", key_url, bytes.NewBuffer(user_pw_json)) // post request to get mojolicious cookie
	key_req.Header.Add("Accept", "application/json")                                //adding header to request
	key_client := &http.Client{}
	key_resp, err := key_client.Do(key_req) //make request and get response
	checkError("Failed key http request: ", err)
	key_resp.Body.Close() //close response body

	//Extract Mojolicious cookie from response headers map
	pattern := regexp.MustCompile(`mojolicious=([A-Za-z0-9\-\_]+);`) //compile regex to extract cookie from response
	mojo_cookie := key_resp.Header["Set-Cookie"]                     // get cookie header from http response
	mojo_key = pattern.FindStringSubmatch(mojo_cookie[0])[1]         //extract cookie from header text
	return
}

//------------------------------------------------------------------------------

//API Request to get assigned services
func assigned_services(mojo_key string) (api_output map[string]string) {

	// API REQUEST FOR ASSIGNED SERVERS
	api_url := ("https://cdnportal.comcast.net/api/1.1/deliveryservices/" + ".json")
	api_req, err := http.NewRequest("GET", api_url, nil)  // post request to get mojolicious cookie
	api_req.Header.Add("Accept", "application/json")      // add header to request
	api_req.Header.Add("Cookie", "mojolicious="+mojo_key) // add header to request with mojo key
	api_client := &http.Client{}
	api_resp, err := api_client.Do(api_req) //make http request to api
	checkError("Failed api client request: ", err)
	// fmt.Printf("Server data API request status for DSID=%s is: %s \n", dsid, api_resp.Status) //print response status to screen
	api_body, _ := ioutil.ReadAll(api_resp.Body)
	api_resp.Body.Close()

	src := string(api_body)       //store api body to a string variable
	var structObj ServiceIdStruct //create new struct object

	//unmarshal api body of type []byte
	json.Unmarshal([]byte(src), &structObj)

	//empty map to store key=xmlid and value=dsid
	api_output = make(map[string]string)

	//range over the api response, get the xmlid and dsId, store them as key,value in api_output
	for _, v := range structObj.Response {

		//if else statements to check xmlidFilter and longDesc1Filter
		if xmlidFilter == "" && longDesc1Filter == "" {
			fmt.Println("Assigned Service:", v.XMLID)
			api_output[v.XMLID] = strconv.Itoa(v.ID)
		} else if xmlidFilter != "" && longDesc1Filter == "" { //match on xmlID
			if strings.Contains(v.XMLID, xmlidFilter) {
				fmt.Println("Assigned Service:", v.XMLID)

				api_output[v.XMLID] = strconv.Itoa(v.ID)

			}
		} else if xmlidFilter == "" && longDesc1Filter != "" { //match on LongDesc1
			if strings.Contains(v.LongDesc1, longDesc1Filter) {
				fmt.Println("Assigned Service:", v.XMLID)

				api_output[v.XMLID] = strconv.Itoa(v.ID)

			}
		} else if xmlidFilter != "" && longDesc1Filter != "" { //match on both xmlId and LongDesc1
			if strings.Contains(v.LongDesc1, longDesc1Filter) {
				if strings.Contains(v.XMLID, xmlidFilter) {
					fmt.Println("Assigned Service:", v.XMLID)

					api_output[v.XMLID] = strconv.Itoa(v.ID)

				}

			}

		}

	}

	return

}

//------------------------------------------------------------------------------

//struct to parse assigned services API response
type ServiceIdStruct struct {
	Response []struct {
		ID        int    `json:"id"`
		LongDesc1 string `json:"longDesc1"`
		XMLID     string `json:"xmlId"`
	} `json:"response"`
}

//struct that parses the usage data API response
type UsageDataStruct struct {
	Response []struct {
		Stats struct {
			Nine5ThPercentile interface{} `json:"95thPercentile"`
			Sum               interface{} `json:"sum"`
		} `json:"stats"`
	} `json:"response"`
}

//struct to hold channel response
type HTTPResponse struct {
	xmlId     string
	dsId      string
	startdate string
	enddate   string
	sum       string
	nine5th   string
}

//------------------------------------------------------------------------------
//get usage data and send response to channel
func services_api_req(mojo_key string, xmlID string, dsID string, httpCh chan<- HTTPResponse, wg *sync.WaitGroup) {

	defer wg.Done() //decrement the waitgroup counter

	//Initialize various variables
	apiurlValues := []string{}
	var finalStartDate string //date to be apppended to the api url
	var finalEndDate string   //date to be appended to the api url
	var sumStart time.Time    //date to be used if the sumByDay filter is on
	var sumFinish time.Time   //date to be used if the sumByDay filter is on
	layout := "01/02/2006"

	//parse user input dates in this format: 01/02/2006
	cmdStartDate, err := time.Parse(layout, startDate)
	cmdFinishDate, err := time.Parse(layout, endDate)

	checkError("Failed to parse input dates: ", err)

	//if else statement to check if the subByDay filter is specified
	if sumByDay == "" {
		finalStartDate = cmdStartDate.Format(layout)
		finalEndDate = cmdFinishDate.Format(layout)
		api_url_values := ("" + finalStartDate + " " + finalEndDate)

		apiurlValues = append(apiurlValues, api_url_values) //append start and end date

	} else {
		if sumByDay != "" {
			//get each day between start and end date and append it to the api url
			for i := cmdStartDate; i.Month() < cmdFinishDate.Month() || i.Day() < cmdFinishDate.Day(); i = i.AddDate(0, 0, 1) {

				sumStart = i
				sumFinish = i.AddDate(0, 0, 1)

				finalStartDate = sumStart.Format(layout)
				finalEndDate = sumFinish.Format(layout)

				api_url_values := ("" + finalStartDate + " " + finalEndDate)
				apiurlValues = append(apiurlValues, api_url_values) //append start and end date
			}
		}
	}

	//range over the dates and make api request
	for _, url := range apiurlValues {

		//parse start and end date
		s := strings.Split(url, " ")
		date1 := s[0]
		date2 := s[1]

		// convert dates to unix time stamp
		convertdate1 := convertToUnixTimeStamps(date1)
		convertdate2 := convertToUnixTimeStamps(date2)

		finalStartDate = convertdate1
		finalEndDate = convertdate2

		api_url := ("https://cdnportal.comcast.net/api/1.1/deliveryservices/" + dsID + "/server_types/edge/metric_types/kbps/start_date/" + finalStartDate + "/end_date/" + finalEndDate + ".json")
		api_req, err := http.NewRequest("GET", api_url, nil)  // post request to get mojolicious cookie
		api_req.Header.Add("Accept", "application/json")      // add header to request
		api_req.Header.Add("Cookie", "mojolicious="+mojo_key) // add header to request with mojo key
		api_client := &http.Client{}
		api_resp, err := api_client.Do(api_req) //make http request to api
		checkError("Failed api client request: ", err)
		fmt.Printf("Server data API request status for XMLID=%s is: %s startDate:%s endDate: %s \n", xmlID, api_resp.Status, date1, date2) //print response status to screen
		api_body, _ := ioutil.ReadAll(api_resp.Body)
		api_resp.Body.Close()

		//initialize various variables
		var structObj UsageDataStruct //create new struct object
		var Value1, Value2 interface{}
		var sumGBs, Nine5ThPercentileMBs string
		var sumFloat, Nine5ThPercentileFloat, sumGBsFloat, Nine5ThPercentileMBsFloat float64
		var ok bool

		//unmarshal api body of type []byte
		json.Unmarshal(api_body, &structObj)

		//parse api response using the UsageDataStruct
		if len(structObj.Response) == 0 {
			sumGBs = "0"
			Nine5ThPercentileMBs = "0"

		} else {

			Value1 = (structObj.Response[0].Stats.Sum)               //get Sum value from the UsageDataStruct
			Value2 = (structObj.Response[0].Stats.Nine5ThPercentile) //get 95th Percentile from the UsageDataStruct
			if Value1 != nil {
				if sumFloat, ok = Value1.(float64); !ok {
					log.Fatalf("Sum value has unexpected type %T", Value1)
				}
			}

			if Value2 != nil {
				if Nine5ThPercentileFloat, ok = Value2.(float64); !ok {
					log.Fatalf("95th Percentile value has unexpected type %T", Value2)
				}
			}

			sumGBsFloat = sumFloat / 1000000                          //convert kb to gb
			Nine5ThPercentileMBsFloat = Nine5ThPercentileFloat / 1000 //convert kb to mb

		}

		sumGBs = strconv.FormatFloat(sumGBsFloat, 'f', 1, 64) //convert float value to string
		Nine5ThPercentileMBs = strconv.FormatFloat(Nine5ThPercentileMBsFloat, 'f', 1, 64)

		//send values to the channel
		httpCh <- HTTPResponse{xmlID, dsID, date1, date2, sumGBs, Nine5ThPercentileMBs}
	}

}

//------------------------------------------------------------------------------

// WRITE OUTPUT TO FILE
func write_file(output [][]string) {
	t := time.Now()
	t_now := (t.Format("D-2006.01.02_T-15-04-05"))
	path := os.Getenv("HOME") + "/Downloads/scripts/usage_reports/"
	file, err := os.Create(path + "UsageReport_" + t_now + ".csv")
	defer file.Close()
	checkError("Failed to create file", err)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	fmt.Println("Output file location is = ", path+"UsageReport_"+t_now+".csv")
	for _, value := range output {
		err := writer.Write(value)
		checkError("Failed write to file", err)
	}
}

//------------------------------------------------------------------------------

//WRITE OUTPUT TO SCREEN
func write_screen(output [][]string) {
	for _, value := range output {
		fmt.Println(value)
	}
}
