package todb

import (
	"fmt"
)

func getAsn() ([]Asn, error) {
	ret := []Asn{}
	queryStr := "select * from asn"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getCachegroup() ([]Cachegroup, error) {
	ret := []Cachegroup{}
	queryStr := "select * from cachegroup"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getCachegroupParameter() ([]CachegroupParameter, error) {
	ret := []CachegroupParameter{}
	queryStr := "select * from cachegroup_parameter"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getCdn() ([]Cdn, error) {
	ret := []Cdn{}
	queryStr := "select * from cdn"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getDeliveryservice() ([]Deliveryservice, error) {
	ret := []Deliveryservice{}
	queryStr := "select * from deliveryservice"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getDeliveryserviceRegex() ([]DeliveryserviceRegex, error) {
	ret := []DeliveryserviceRegex{}
	queryStr := "select * from deliveryservice_regex"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getDeliveryserviceServer() ([]DeliveryserviceServer, error) {
	ret := []DeliveryserviceServer{}
	queryStr := "select * from deliveryservice_server"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getDeliveryserviceTmuser() ([]DeliveryserviceTmuser, error) {
	ret := []DeliveryserviceTmuser{}
	queryStr := "select * from deliveryservice_tmuser"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getDivision() ([]Division, error) {
	ret := []Division{}
	queryStr := "select * from division"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getFederation() ([]Federation, error) {
	ret := []Federation{}
	queryStr := "select * from federation"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getFederationDeliveryservice() ([]FederationDeliveryservice, error) {
	ret := []FederationDeliveryservice{}
	queryStr := "select * from federation_deliveryservice"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getFederationFederationResolver() ([]FederationFederationResolver, error) {
	ret := []FederationFederationResolver{}
	queryStr := "select * from federation_federation_resolver"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getFederationResolver() ([]FederationResolver, error) {
	ret := []FederationResolver{}
	queryStr := "select * from federation_resolver"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getFederationTmuser() ([]FederationTmuser, error) {
	ret := []FederationTmuser{}
	queryStr := "select * from federation_tmuser"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getGooseDbVersion() ([]GooseDbVersion, error) {
	ret := []GooseDbVersion{}
	queryStr := "select * from goose_db_version"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getHwinfo() ([]Hwinfo, error) {
	ret := []Hwinfo{}
	queryStr := "select * from hwinfo"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getJob() ([]Job, error) {
	ret := []Job{}
	queryStr := "select * from job"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getJobAgent() ([]JobAgent, error) {
	ret := []JobAgent{}
	queryStr := "select * from job_agent"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getJobResult() ([]JobResult, error) {
	ret := []JobResult{}
	queryStr := "select * from job_result"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getJobStatus() ([]JobStatus, error) {
	ret := []JobStatus{}
	queryStr := "select * from job_status"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getLog() ([]Log, error) {
	ret := []Log{}
	queryStr := "select * from log"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getParameter() ([]Parameter, error) {
	ret := []Parameter{}
	queryStr := "select * from parameter"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getPhysLocation() ([]PhysLocation, error) {
	ret := []PhysLocation{}
	queryStr := "select * from phys_location"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getProfile() ([]Profile, error) {
	ret := []Profile{}
	queryStr := "select * from profile"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getProfileParameter() ([]ProfileParameter, error) {
	ret := []ProfileParameter{}
	queryStr := "select * from profile_parameter"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getRegex() ([]Regex, error) {
	ret := []Regex{}
	queryStr := "select * from regex"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getRegion() ([]Region, error) {
	ret := []Region{}
	queryStr := "select * from region"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getRole() ([]Role, error) {
	ret := []Role{}
	queryStr := "select * from role"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getServer() ([]Server, error) {
	ret := []Server{}
	queryStr := "select * from server"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getServercheck() ([]Servercheck, error) {
	ret := []Servercheck{}
	queryStr := "select * from servercheck"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getStaticdnsentry() ([]Staticdnsentry, error) {
	ret := []Staticdnsentry{}
	queryStr := "select * from staticdnsentry"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getStatsSummary() ([]StatsSummary, error) {
	ret := []StatsSummary{}
	queryStr := "select * from stats_summary"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getStatus() ([]Status, error) {
	ret := []Status{}
	queryStr := "select * from status"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getTmUser() ([]TmUser, error) {
	ret := []TmUser{}
	queryStr := "select * from tm_user"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getToExtension() ([]ToExtension, error) {
	ret := []ToExtension{}
	queryStr := "select * from to_extension"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func getType() ([]Type, error) {
	ret := []Type{}
	queryStr := "select * from type"
	err := globalDB.Select(&ret, queryStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func GetTable(tableName string) (interface{}, error) {
	if tableName == "asn" {
		return getAsn()
	}
	if tableName == "cachegroup" {
		return getCachegroup()
	}
	if tableName == "cachegroup_parameter" {
		return getCachegroupParameter()
	}
	if tableName == "cdn" {
		return getCdn()
	}
	if tableName == "deliveryservice" {
		return getDeliveryservice()
	}
	if tableName == "deliveryservice_regex" {
		return getDeliveryserviceRegex()
	}
	if tableName == "deliveryservice_server" {
		return getDeliveryserviceServer()
	}
	if tableName == "deliveryservice_tmuser" {
		return getDeliveryserviceTmuser()
	}
	if tableName == "division" {
		return getDivision()
	}
	if tableName == "federation" {
		return getFederation()
	}
	if tableName == "federation_deliveryservice" {
		return getFederationDeliveryservice()
	}
	if tableName == "federation_federation_resolver" {
		return getFederationFederationResolver()
	}
	if tableName == "federation_resolver" {
		return getFederationResolver()
	}
	if tableName == "federation_tmuser" {
		return getFederationTmuser()
	}
	if tableName == "goose_db_version" {
		return getGooseDbVersion()
	}
	if tableName == "hwinfo" {
		return getHwinfo()
	}
	if tableName == "job" {
		return getJob()
	}
	if tableName == "job_agent" {
		return getJobAgent()
	}
	if tableName == "job_result" {
		return getJobResult()
	}
	if tableName == "job_status" {
		return getJobStatus()
	}
	if tableName == "log" {
		return getLog()
	}
	if tableName == "parameter" {
		return getParameter()
	}
	if tableName == "phys_location" {
		return getPhysLocation()
	}
	if tableName == "profile" {
		return getProfile()
	}
	if tableName == "profile_parameter" {
		return getProfileParameter()
	}
	if tableName == "regex" {
		return getRegex()
	}
	if tableName == "region" {
		return getRegion()
	}
	if tableName == "role" {
		return getRole()
	}
	if tableName == "server" {
		return getServer()
	}
	if tableName == "servercheck" {
		return getServercheck()
	}
	if tableName == "staticdnsentry" {
		return getStaticdnsentry()
	}
	if tableName == "stats_summary" {
		return getStatsSummary()
	}
	if tableName == "status" {
		return getStatus()
	}
	if tableName == "tm_user" {
		return getTmUser()
	}
	if tableName == "to_extension" {
		return getToExtension()
	}
	if tableName == "type" {
		return getType()
	}
	return nil, nil
}
