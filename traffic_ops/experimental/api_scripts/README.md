# Overview


## How to run go scripts:
1. Create a folder called /Downloads/scripts
2. Store your CDN_API_Credentials.txt in this folder (user, pw on separate lines)

## To run cdn_api_mojokey.go:
1. go run cdn_api_mojokey.go
2. It will return a mojokey that you can use to authenticate all other API requests

## Billing Script:
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

## To run billing_script.go
1. go run billing_script.go -startDate="01/01/2017" -endDate="02/01/2017" -xmlidFilter="" -longDesc1Filter="" -sumByDay=""
2. Required inputs: -startDate, -endDate will get the usage data for the specified dates. To get usage data for the month of January 2017, startDate="01/01/2017" which will begin aggregation of data from Jan 1st, 2017 @ 12am UTC until the endDate="02/01/2017" or Feb 1st, 2017 12am UTC.
3. To get all services assigned to your username, only specify the start and end dates
4. Optional Filters: xmlidFilter, longDesc1Filter, sumByDay
5. -sumByDay will get a daily breakdown of the usage data in other words, get usage data for each day between the start and end date. To activate this filter, do: -sumByDay="y"
6. -xmlidFilter will return services that match the xmlid in your list of services. To activate this filter, do: -xmlidFilter="servicename-subname" which will return service(s) that contains this string: "servicename-subname". If left blank, it will return all the services.
7. -longDesc1Filter will only return the services that are in production and where the customer field contains "cts; prod". If left blank, it will return all the production, trials and demo services assigned to your username. To activate this filter, do: -longDesc1Filter="cts; prod"
