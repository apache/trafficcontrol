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
	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

// we need a type alias to define functions on
type TOCDNFederation struct {
	ReqInfo *api.APIInfo `json:"-"`
	v13.CDNFederation
}

// Used for all CRUD routes
func GetTypeSingleton() api.CRUDFactory {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOCDNFederation{reqInfo, v13.CDNFederation{}}
		return &toReturn
	}
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetKeys() (map[string]interface{}, bool) {
	if fed.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *fed.ID}, true
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetAuditName() string {
	if fed.CName != nil {
		return *fed.CName
	}
	if fed.ID != nil {
		return strconv.Itoa(*fed.ID)
	}
	return "unknown"
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetType() string {
	return "cdnfederation"
}

// Fufills `Create' interface
func (fed *TOCDNFederation) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) // non-panicking type assertion
	fed.ID = &i
}

// Fulfills `Validate' interface
func (fed *TOCDNFederation) Validate() error {

	isDNSName := validation.NewStringRule(govalidator.IsDNSName, "must be a valid hostname")
	endsWithDot := validation.NewStringRule(
		func(str string) bool {
			return strings.HasSuffix(str, ".")
		}, "must end with a period")

	// cname regex: (^\S*\.$), ttl regex: (^\d+$)
	validateErrs := validation.Errors{
		"cname": validation.Validate(fed.CName, validation.Required, endsWithDot, isDNSName),
		"ttl":   validation.Validate(fed.TTL, validation.Required, validation.Min(0)),
	}
	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

// This separates out errors depending on whether or not some constraint prevented
// the operation from occuring.
func parseQueryError(parseErr error, method string) (error, tc.ApiErrorType) {
	if pqErr, ok := parseErr.(*pq.Error); ok {
		err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
		if eType == tc.DataConflictError {
			log.Errorf("data conflict error: %v", err)
			return errors.New("a federation with " + err.Error()), eType
		}
		return err, eType
	} else {
		log.Errorf("received error: %++v from %s execution", parseErr, method)
		return tc.DBError, tc.SystemError
	}
}

// fed.ReqInfo.Params["name"] is not used on creation, rather the cdn name
// is connected when the federations/:id/deliveryservice links a federation
// Note: cdns and deliveryservies have a 1-1 relationship
func (fed *TOCDNFederation) Create() (error, tc.ApiErrorType) {

	// Deliveryservice IDs should not be included on create.
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
	var lastUpdated tc.TimeNoMod

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err = resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id and last_updated from insert: %s\n", err)
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
	fed.LastUpdated = &lastUpdated

	return nil, tc.NoError
}

// returning true indicates the data related to the given tenantID should be visible
// `tenantIDs` is presumed to be unsorted, and a nil tenantID is viewable by everyone
func checkTenancy(tenantID *int, tenantIDs []int) bool {
	if tenantID == nil {
		return true
	}
	for _, id := range tenantIDs {
		if id == *tenantID {
			return true
		}
	}
	return false
}

func (fed *TOCDNFederation) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {

	// Cannot perform query on tenantID while "rows" aren't closed (limitation of
	// psql), so we need to get the valid tenentIDs ahead of time.
	tenantIDs, err := tenant.GetUserTenantIDListTx(fed.ReqInfo.Tx.Tx, fed.ReqInfo.User.TenantID)
	if err != nil {
		log.Errorf("getting tenant list for user: %v\n", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}

	var query string
	_, id := parameters["id"]
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"id": dbhelpers.WhereColumnInfo{Column: "federation.id", Checker: api.IsInt},
	}
	if !id {
		queryParamsToQueryCols["name"] = dbhelpers.WhereColumnInfo{Column: "cdn.name", Checker: nil}
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)

	if id {
		// Can't use `AddTenancyCheck` for id because then we won't know what caused
		// an empty response, so the tenancy check will be performed below.
		query = selectByID()
	} else { // searching by name
		query = selectByCDNName()
		where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "(tenant_id", tenantIDs)
		where += " OR tenant_id IS NULL)"
	}

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

		var tenantID *int
		fed.DeliveryServiceIDs = &v13.DeliveryServiceIDs{}
		if err = rows.Scan(&tenantID, &fed.ID, &fed.CName, &fed.TTL, &fed.Description, &fed.LastUpdated, &fed.DsId, &fed.XmlId); err != nil {
			log.Errorf("parsing federation rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		if fed.ID == nil {
			log.Errorf("unexpected nil id")
			return nil, []error{tc.DBError}, tc.SystemError
		}

		// if we are getting by id, there may not be an attached deliveryservice
		if id && fed.DsId == nil {
			fed.DeliveryServiceIDs = nil
		}

		// append if by cdn or if tenancy check for id passes
		if !id || checkTenancy(tenantID, tenantIDs) {
			federations = append(federations, *fed)
		} else { // id is true and the tenancy check failed
			return nil, []error{errors.New("user not authorized for requested federation")}, tc.ForbiddenError
		}
	}

	// if federations yields "response": []
	if len(federations) == 0 {

		if id {
			return nil, []error{errors.New("resource not found")}, tc.DataMissingError
		}

		if yes, err := dbhelpers.CDNExists(parameters["name"], fed.ReqInfo.Tx); yes {
			return federations, []error{}, tc.NoError
		} else if err != nil { // internal server error
			log.Errorf("verifying cdn exists: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		// the query ran as expected and the cdn does not exist
		return nil, []error{errors.New("resource not found")}, tc.DataMissingError
	}

	return federations, []error{}, tc.NoError
}

func (fed *TOCDNFederation) Update() (error, tc.ApiErrorType) {

	if ok, err := fed.isTenantAuthorized(); err != nil {
		log.Errorf("checking tenancy: %v", err)
		return tc.DBError, tc.SystemError
	} else if !ok {
		return tc.TenantUserNotAuthError, tc.ForbiddenError
	}

	// Deliveryservice IDs should not be included on update.
	if fed.DeliveryServiceIDs != nil {
		fed.DsId = nil
		fed.XmlId = nil
		fed.DeliveryServiceIDs = nil
	}

	resultRows, err := fed.ReqInfo.Tx.NamedQuery(updateQuery(), fed)
	defer resultRows.Close()
	if err != nil {
		return parseQueryError(err, "update")
	}

	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	fed.LastUpdated = &lastUpdated

	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no federation found with this id"), tc.DataMissingError
		}
		return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
	}

	return nil, tc.NoError
}

// In the perl version, :name is ignored. It is not even verified whether or not
// :name is a real cdn that exists. This mimicks the perl behavior.
func (fed *TOCDNFederation) Delete() (error, tc.ApiErrorType) {

	if ok, err := fed.isTenantAuthorized(); err != nil {
		log.Errorf("checking tenancy: %v", err)
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
		}
		return fmt.Errorf("this delete affected too many rows: %d", rowsAffected), tc.SystemError
	}
	return nil, tc.NoError
}

// Function not exported because although DELETE and UPDATE have normal tenancy check,
// CREATE does not. No ds is associated on create. This isn't used for READ because
// psql doesn't like nested queries within the same transaction.
func (fed TOCDNFederation) isTenantAuthorized() (bool, error) {

	tenantID, err := getTenantIDFromFedID(*fed.ID, fed.ReqInfo.Tx.Tx)
	if err != nil {
		// If nobody has claimed a tenant, that federation is publicly visible.
		// This logically follows /federations/:id/deliveryservices
		if err == sql.ErrNoRows {
			return true, nil
		}
		log.Errorf("getting tenant id from federation: %v", err)
		return false, err
	}

	// TODO: After IsResourceAuthorizedToUserTx is updated to no longer have `use_tenancy`,
	// that will probably be better to use. For now, use the list. Issue #2602
	list, err := tenant.GetUserTenantIDListTx(fed.ReqInfo.Tx.Tx, fed.ReqInfo.User.TenantID)
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
	SELECT ds.tenant_id FROM federation AS f
	JOIN federation_deliveryservice AS fd ON f.id = fd.federation
	JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice
	WHERE f.id = $1`
	err := tx.QueryRow(query, id).Scan(&tenantID)
	return tenantID, err
}

func selectByID() string {
	return `
	SELECT
	ds.tenant_id,
	federation.id AS id,
	federation.cname,
	federation.ttl,
	federation.description,
	federation.last_updated,
	ds.id AS ds_id,
	ds.xml_id
	FROM federation
	LEFT JOIN federation_deliveryservice AS fd ON federation.id = fd.federation
	LEFT JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice`
	// WHERE federation.id = :id (determined by dbhelper)
}

func selectByCDNName() string {
	return `
	SELECT
	ds.tenant_id,
	federation.id AS id,
	federation.cname,
	federation.ttl,
	federation.description,
	federation.last_updated,
	ds.id AS ds_id,
	ds.xml_id
	FROM federation
	JOIN federation_deliveryservice AS fd ON federation.id = fd.federation
	JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice
	JOIN cdn ON cdn.id = cdn_id`
	// WHERE cdn.name = :cdn_name (determined by dbhelper)
}

func updateQuery() string {
	return `
	UPDATE federation SET
	cname=:cname,
	ttl=:ttl,
	description=:description
	WHERE id=:id
	RETURNING last_updated`
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
	) RETURNING id, last_updated`
}

func deleteQuery() string {
	return `DELETE FROM federation WHERE id = :id`
}
