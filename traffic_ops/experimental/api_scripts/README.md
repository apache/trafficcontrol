# Overview

## How to run go scripts:
1. Create a folder called /Downloads/scripts
2. Store your CDN_API_Credentials.txt in this folder (user, pw on separate lines)
3. Create a folder called/Downloads/scripts/usage_reports to store usage_reporting_CustomerCopy.go outputs

## Usage Reporting:
This program is intended to get usage data (Sum in Gbs and 95th Percentile in MBs) from various services via this api: https://cdnportal.comcast.net/api/1.1/deliveryservices/999/server_types/edge/metric_types/kbps/start_date/1453223126/end_date/1453309526

The usage data will be printed on the screen as well as stored in a csv file

The program does the following things:
1. Parses various user inputs such as start date, end date, xmlidFilter, longDesc1Filter, sumByDay
2. Reads user credentials from a file
3. Makes an API request using #2 and gets a mojo key
4. The mojo key is used to authenticate all other API requests
5. Makes an API request to get all the assigned services using the mojo key
6. Makes an API request to get usage data for all the obtained services from #5
7. If the script run successfully, prints the usage data on the screen
8. Creates a csv file and store it under /Downloads/scripts/usage_reports

## To run usage_reporting_CustomerCopy.go
1. go run usage_reporting_CustomerCopy.go -startDate="02/01/2017" -endDate="03/01/2017" -xmlidFilter="" -longDesc1Filter="" -sumByDay=""
2. Required inputs: -startDate, -endDate will get the usage data for the specified dates
3. Optional Filters: xmlidFilter, longDesc1Filter, sumByDay
4. -sumByDay will get a daily breakdown of the usage data or get usage data for each day between the start and end date.
5. -xmlidFilter will return services that match the xmlid in your list of services. Example: -xmlidFilter="servicename-subname" will return services that contain that string provided in xmlid. If left blank, it will return all the services.
6. -longDesc1Filter will only return services where the customer field contains "cts; prod". If left blank, it will return all the services that are in production, trials or demo services.
7. To get all services assigned to your username, only specify the start and end dates
