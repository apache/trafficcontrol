package cdnfederation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOCDNFederation struct {
	ReqInfo *api.APIInfo `json:"-"`
	v13.CDNFederationNullable
}

//Used for all CRUD routes
func GetTypeSingleton() api.CRUDFactory {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOCDNFederation{reqInfo, v13.CDNFederationNullable{}}
		return &toReturn
	}
}

//Fufills `Identifier' interface
func (fed TOCDNFederation) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

//Fufills `Identifier' interface
func (fed TOCDNFederation) GetKeys() (map[string]interface{}, bool) {
	if fed.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *fed.ID}, true
}

//Fufills `Identifier' interface
func (fed TOCDNFederation) GetAuditName() string {
	if fed.CName != nil {
		return *fed.CName
	}
	if fed.ID != nil {
		return strconv.Itoa(*fed.ID)
	}
	return "unknown"
}

//Fufills `Identifier' interface
func (fed TOCDNFederation) GetType() string {
	return "cdnfederation"
}

//Fufills `Create' interface
func (fed *TOCDNFederation) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //non-panicking type assertion
	fed.ID = &i
}

//Fulfills `Validate' interface
func (fed *TOCDNFederation) Validate() error {

	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")
	endsWithDot := validation.NewStringRule(
		func(str string) bool {
			return strings.HasSuffix(str, ".")
		}, "must end with a period")

	//cname regex: (^\S*\.$), ttl regex: (^\d+$)
	validateErrs := validation.Errors{
		"cname": validation.Validate(fed.CName, validation.Required, endsWithDot, noSpaces),
		"ttl":   validation.Validate(fed.Ttl, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

//This separates out errors depending on whether or not some constraint prevented
//the operation from occuring.
func parseQueryError(parseErr error, method string) (error, tc.ApiErrorType) {
	if pqErr, ok := parseErr.(*pq.Error); ok {
		err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
		if eType == tc.DataConflictError {
			return errors.New("a federation with " + err.Error()), eType
		}
		return err, eType
	} else {
		log.Errorf("received error: %++v from %s execution", parseErr, method)
		return tc.DBError, tc.SystemError
	}
}

//fed.ReqInfo.Params["name"] is not used on creation, rather the cdn name
//is connected when the federations/:id/deliveryservice links a federation
//Note: cdns and deliveryservies have a 1-1 relationship
func (fed *TOCDNFederation) Create() (error, tc.ApiErrorType) {

	if fed.ReqInfo.User.PrivLevel < auth.PrivLevelAdmin {
		return errors.New("Forbidden"), tc.ForbiddenError
	}

	// Boilerplate code below
	resultRows, err := fed.ReqInfo.Tx.NamedQuery(insertQuery(), fed)
	if err != nil {
		return parseQueryError(err, "create")
	}
	defer resultRows.Close()

	var id int
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err = resultRows.Scan(&id); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no federation was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from fed insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	fed.SetKeys(map[string]interface{}{"id": id})

	return nil, tc.NoError
}

//Concerning efficiency, maybe it would be better to conditionally run different queries?
//This way the handler wouldn't be doing a full outer join for every query.
func (fed *TOCDNFederation) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {

	_, ok_id := parameters["id"]
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		//db tag                                         symbol from query
		"id":          dbhelpers.WhereColumnInfo{Column: "federation.id", Checker: api.IsInt},
		"cname":       dbhelpers.WhereColumnInfo{Column: "cname", Checker: nil},
		"ttl":         dbhelpers.WhereColumnInfo{Column: "ttl", Checker: api.IsInt},
		"description": dbhelpers.WhereColumnInfo{Column: "description", Checker: nil},
		"xmlId":       dbhelpers.WhereColumnInfo{Column: "xml_id", Checker: nil},
		"ds_id":       dbhelpers.WhereColumnInfo{Column: "deliveryservice.id", Checker: api.IsInt},

		//"name":     dbhelpers.WhereColumnInfo{Column: "cdn.name", Checker: nil}, //would narrow scope of search
		"cdn_name": dbhelpers.WhereColumnInfo{Column: "cdn.name", Checker: nil}, //used to verify cdn name exists
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := fed.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying federations: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	cdnFound := false
	federations := []interface{}{}
	for rows.Next() {
		var fed v13.CDNFederationNullable
		if err = rows.StructScan(&fed); err != nil {
			log.Errorf("error parsing federation rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}

		//if we have an id, the cdn name is fooable
		if ok_id || fed.Name != nil && *fed.Name == parameters["name"] {
			cdnFound = true

			//if we are getting by id, there may not be an attached deliveryservice
			//DeliveryServiceIDsNullable will not be nil itself, due to the struct scan
			if ok_id && fed.DeliveryServiceIDsNullable.ID == nil {
				fed.DeliveryServiceIDsNullable = nil
			}

			//Never return cnd and deliveryService only information. Just federation data.
			if fed.ID != nil {
				federations = append(federations, fed)
			}
		}
	}

	if !cdnFound {
		return nil, []error{errors.New("Resource not found.")}, tc.DataMissingError
	}

	return federations, []error{}, tc.NoError
}

func (fed *TOCDNFederation) Update() (error, tc.ApiErrorType) {

	if fed.ReqInfo.User.PrivLevel < auth.PrivLevelAdmin {
		return errors.New("Forbidden"), tc.ForbiddenError
	}

	result, err := fed.ReqInfo.Tx.NamedExec(updateQuery(), fed)
	if err != nil {
		return parseQueryError(err, "update")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}

	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no federation found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//In the perl version, :name is ignored. It is not even verified whether or not
//:name is a real cdn that exists. This mimicks the perl behavior.
func (fed *TOCDNFederation) Delete() (error, tc.ApiErrorType) {

	if fed.ReqInfo.User.PrivLevel < auth.PrivLevelAdmin {
		return errors.New("Forbidden"), tc.ForbiddenError
	}

	log.Debugf("about to run exec query: %s with federation: %++v", deleteQuery(), fed)
	result, err := fed.ReqInfo.Tx.NamedExec(deleteQuery(), fed)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no federation with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this delete affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//Full outer joins are used so that we can get a federation without an attached deliveryservice (fetching by id)
//and so that we can retrieve all cdn names to determine if the requested cdn exists.
func selectQuery() string {
	return `
  SELECT federation.id as id, cname, ttl, description, deliveryservice.id as ds_id, xml_id, cdn.name as cdn_name FROM federation
  FULL OUTER JOIN federation_deliveryservice as fd ON federation.id = fd.federation
  FULL OUTER JOIN deliveryservice ON deliveryservice.id = fd.deliveryservice
  FULL OUTER JOIN cdn ON cdn.id = cdn_id`
}

func updateQuery() string {
	return `
  UPDATE federation SET
  cname=:cname
  ttl=:ttl
  description=:description
  WHERE id=:id`
}

func insertQuery() string {
	return `
  INSERT INTO federation (
    cname,
    ttl,
    description
  ) VALUES (
    :cname,
    :ttl,
    :description
  ) RETURNING id`
}

func deleteQuery() string {
	return `DELETE FROM federation WHERE id = :id`
}
