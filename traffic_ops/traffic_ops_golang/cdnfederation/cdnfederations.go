package cdnfederation

import (
	"database/sql"
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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOCDNFederation struct {
	ReqInfo *api.APIInfo `json:"-"`
	v13.CDNFederation
}

//Used for all CRUD routes
func GetTypeSingleton() api.CRUDFactory {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOCDNFederation{reqInfo, v13.CDNFederation{}}
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

	//Deliveryservice IDs should not be included on create.
	if fed.DeliveryServiceIDs != nil {
		fed.DsId = nil
		fed.XmlId = nil
		fed.DeliveryServiceIDs = nil
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

func (fed *TOCDNFederation) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {

	//Cannot perform query on tenantID while "rows" aren't closed (limitation of
	//psql), so we need to get the valid tenentIDs ahead of time.
	tenantIDs, err := tenant.GetUserTenantIDListTx(fed.ReqInfo.Tx.Tx, fed.ReqInfo.User.TenantID)
	if err != nil {
		log.Errorf("getting tenant list for user: %v\n", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}

	var query string
	_, id := parameters["id"]
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		//db tag                                         symbol from query
		"id":          dbhelpers.WhereColumnInfo{Column: "federation.id", Checker: api.IsInt},
		"cname":       dbhelpers.WhereColumnInfo{Column: "cname", Checker: nil},
		"ttl":         dbhelpers.WhereColumnInfo{Column: "ttl", Checker: api.IsInt},
		"description": dbhelpers.WhereColumnInfo{Column: "description", Checker: nil},
		"xmlId":       dbhelpers.WhereColumnInfo{Column: "xml_id", Checker: nil},
		"ds_id":       dbhelpers.WhereColumnInfo{Column: "deliveryservice.id", Checker: api.IsInt},
	}
	if id {
		query = selectByID()
	} else { //searching by name
		queryParamsToQueryCols["name"] = dbhelpers.WhereColumnInfo{Column: "cdn.name", Checker: nil}
		query = selectByCDNName()
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "(tenant_id", tenantIDs)
	where += " OR tenant_id IS NULL)"
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query += where + orderBy
	log.Debugln("Query is ", query)

	rows, err := fed.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("querying federations: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	federations := []interface{}{}
	for rows.Next() {

		if err = rows.StructScan(&fed.CDNFederation); err != nil {
			log.Errorf("parsing federation rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		if fed.ID == nil {
			log.Errorf("unexpected nil id")
			return nil, []error{tc.DBError}, tc.SystemError
		}

		//if we are getting by id, there may not be an attached deliveryservice
		//DeliveryServiceIDs will not be nil itself, due to the struct scan
		if id && fed.DsId == nil {
			fed.DeliveryServiceIDs = nil
		}

		federations = append(federations, fed)
	}

	//if federations yields "response": []
	if len(federations) == 0 {

		if id {
			return nil, []error{tc.TenantDSUserNotAuthError}, tc.ForbiddenError
		}

		if yes, err := dbhelpers.CDNExists(parameters["name"], fed.ReqInfo.Tx); yes {
			return federations, []error{}, tc.NoError
		} else if err != nil { //internal server error
			log.Errorf("verifying cdn exists: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		} else { //the query ran as expected and the cdn does not exist
			return nil, []error{errors.New("Resource not found.")}, tc.DataMissingError
		}
	}

	return federations, []error{}, tc.NoError
}

func (fed *TOCDNFederation) Update() (error, tc.ApiErrorType) {

	if ok, err := fed.isTenantAuthorized(); err != nil {
		log.Errorf("checking tenacy: %v", err)
		return tc.DBError, tc.SystemError
	} else if !ok {
		return tc.TenantUserNotAuthError, tc.ForbiddenError
	}

	//Deliveryservice IDs should not be included on update.
	if fed.DeliveryServiceIDs != nil {
		fed.DsId = nil
		fed.XmlId = nil
		fed.DeliveryServiceIDs = nil
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
		}
		return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
	}

	return nil, tc.NoError
}

//In the perl version, :name is ignored. It is not even verified whether or not
//:name is a real cdn that exists. This mimicks the perl behavior.
func (fed *TOCDNFederation) Delete() (error, tc.ApiErrorType) {

	if ok, err := fed.isTenantAuthorized(); err != nil {
		log.Errorf("checking tenacy: %v", err)
		return tc.DBError, tc.SystemError
	} else if !ok {
		return tc.TenantUserNotAuthError, tc.ForbiddenError
	}

	log.Debugf("about to run exec query: %s with federation: %++v", deleteQuery(), fed)
	result, err := fed.ReqInfo.Tx.NamedExec(deleteQuery(), fed)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("getting rows affected: %v", err)
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

//Function not exported because although DELETE and UPDATE have normal tenacy check,
//CREATE does not. No ds is associated on create. This isn't used for READ because
//psql doesn't like nested queries.
func (fed TOCDNFederation) isTenantAuthorized() (bool, error) {
	if fed.ID == nil {
		log.Errorf("unexpected nil id\n")
	}
	//Note: the tenantID is not found via a recursive search. The recursive search
	//is done in IsResourceAuthorizedToUser below.
	tenantID, err := getTenantIDFromFedID(*fed.ID, fed.ReqInfo.Tx.Tx)
	if err != nil {
		//If nobody has claimed a tenant, that federation is publicly visible.
		//This logically follows /federations/:id/deliveryservices
		if err == sql.ErrNoRows {
			log.Errorf("no tenacy")
			return true, nil
		}
		log.Errorf("ran into error %v", err)
		return false, err
	}

	//After IsResourceAuthorizedToUserTx is updated to no longer have `use_tenancy`,
	//that will probably be better to use. For now, use the list.
	list, err := tenant.GetUserTenantIDListTx(fed.ReqInfo.Tx.Tx, fed.ReqInfo.User.ID)
	if err != nil {
		return false, err
	}
	for _, id := range list {
		if id == tenantID {
			return true, nil
		}
	}
	return false, nil
}

func getTenantIDFromFedID(id int, tx *sql.Tx) (int, error) {
	tenantID := 0
	query := `
	SELECT tenant_id from federation
	JOIN federation_deliveryservice as fd ON federation.id = fd.federation
	JOIN deliveryservice as ds ON ds.id = fd.deliveryservice
	WHERE federation.id = $1`
	err := tx.QueryRow(query, id).Scan(&tenantID)
	return tenantID, err
}

func selectByID() string {
	return `
	SELECT federation.id as id, cname, ttl, description, ds.id as ds_id, xml_id FROM federation
	LEFT JOIN federation_deliveryservice as fd ON federation.id = fd.federation
	LEFT JOIN deliveryservice as ds ON ds.id = fd.deliveryservice`
	//WHERE federation.id = :id (determined by dbhelper)
}

func selectByCDNName() string {
	return `
	SELECT federation.id as id, cname, ttl, description, ds.id as ds_id, xml_id FROM federation
	JOIN federation_deliveryservice as fd ON federation.id = fd.federation
	JOIN deliveryservice as ds ON ds.id = fd.deliveryservice
	JOIN cdn ON cdn.id = cdn_id`
	//WHERE cdn.name = :cdn_name (determined by dbhelper)
}

func updateQuery() string {
	return `
	UPDATE federation SET
	cname=:cname,
	ttl=:ttl,
	description=:description
	WHERE id=:id
	`
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
