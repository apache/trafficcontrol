package toextension

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

// CreateUpdateServercheck handles creating or updating an existing servercheck
func CreateTOExtension(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.User.UserName != "extension" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("invalid user for this API. Only the \"extension\" user can use this"), nil)
		return
	}

	toExt := tc.TOExtensionNullable{}

	// Validate request body
	if err := api.Parse(r.Body, inf.Tx.Tx, &toExt); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	// Get Type ID
	typeID, exists, err := dbhelpers.GetTypeIDByName(*toExt.Type, inf.Tx.Tx)
	if !exists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("expected type %v does not exist in type table", *toExt.Type))
		return
	} else if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	toExt.TypeID = &typeID

	id := 0
	successMsg := "Extension Loaded."
	errCode = http.StatusInternalServerError
	if strings.Contains(*toExt.Type, "CHECK_EXTENSION") {
		successMsg = "Check Extension Loaded."
		id, userErr, sysErr = createCheckExt(toExt, inf.Tx)
		if userErr != nil {
			errCode = http.StatusBadRequest
		}
	} else {
		id, sysErr = createNonCheckExt(toExt, inf.Tx)
	}
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	resp := tc.TOExtensionPostResponse{
		Response: tc.TOExtensionID{ID: id},
	}
	api.CreateChangeLogRawTx(api.ApiChange, successMsg, inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, successMsg, resp)
}

func createNonCheckExt(toExt tc.TOExtensionNullable, tx *sqlx.Tx) (int, error) {
	resultRows, err := tx.NamedQuery(insertQuery(), toExt)
	if err != nil {
		return 0, fmt.Errorf("inserting extension: %v", err)
	}
	defer resultRows.Close()

	id := 0
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id); err != nil {
			return 0, fmt.Errorf("scanning id from TO Extension insert: %v", err)
		}
	}
	if rowsAffected == 0 {
		return 0, errors.New("no rows affects from TO Extension insert")
	} else if rowsAffected > 1 {
		return 0, errors.New("too many ids returned from TO Extension insert")
	}

	return id, nil
}

func createCheckExt(toExt tc.TOExtensionNullable, tx *sqlx.Tx) (int, error, error) {
	id := 0
	dupErr, sysErr := checkDupTOCheckExtension("name", *toExt.Name, tx)
	if dupErr != nil || sysErr != nil {
		return 0, dupErr, sysErr
	}

	dupErr, sysErr = checkDupTOCheckExtension("servercheck_short_name", *toExt.ServercheckShortName, tx)
	if dupErr != nil || sysErr != nil {
		return 0, dupErr, sysErr
	}

	// Get open slot
	scc := ""
	if err := tx.Tx.QueryRow(`
	SELECT id, servercheck_column_name
	FROM to_extension 
	WHERE type in 
		(SELECT id FROM type WHERE name = 'CHECK_EXTENSION_OPEN_SLOT')
	ORDER BY servercheck_column_name
	LIMIT 1`).Scan(&id, &scc); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("No open slots left for checks, delete one first."), nil

		}
		return 0, nil, fmt.Errorf("querying open slot to_extension: %v", err)
	}
	toExt.ID = &id
	_, err := tx.NamedExec(updateQuery(), toExt)
	if err != nil {
		return 0, nil, fmt.Errorf("update open extension slot to new check extension: %v", err)
	}

	_, err = tx.Tx.Exec(fmt.Sprintf("UPDATE servercheck set %v = 0", scc))
	if err != nil {
		return 0, nil, fmt.Errorf("reset servercheck table for new check extension: %v", err)
	}
	return id, nil, nil

}

func updateQuery() string {
	return `
	UPDATE to_extension SET
	name=:name,
	version=:version,
	info_url=:info_url,
	script_file=:script_file,
	isactive=:isactive,
	additional_config_json=:additional_config_json,
	description=:description,
	servercheck_short_name=:servercheck_short_name,
	type=:type
	WHERE id=:id
	`
}

func insertQuery() string {
	return `
	INSERT INTO to_extension (
	  name,
	  version,
	  info_url,
	  script_file,
	  isactive,
	  additional_config_json,
	  description,
	  servercheck_short_name,
	  servercheck_column_name,
	  type
	)
	VALUES (
	  :name,
	  :version,
	  :info_url,
	  :script_file,
	  :isactive,
	  :additional_config_json,
	  :description,
	  :servercheck_short_name,
	  :servercheck_column_name,
	  :type
	)
	RETURNING id
	`
}

func checkDupTOCheckExtension(colName, value string, tx *sqlx.Tx) (error, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM to_extension WHERE %v =$1)", colName)
	exists := false
	err := tx.Tx.QueryRow(query, value).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("checking if to_extension %v already exists: %v", colName, err)
	}
	if exists {
		return fmt.Errorf("A Check extension is already loaded with %v %v", value, colName), nil
	}
	return nil, nil
}

// type TOExtension struct {
// 	api.APIInfoImpl `json:"-"`
// 	tc.TOExtensionNullable
// }

// func (e *TOExtension) SetLastUpdated(t tc.TimeNoMod) {
// 	e.LastUpdated = &t
// }

// func (e *TOExtension) NewReadObj() interface{} {
// 	return &tc.TOExtensionNullable{}
// }

// // func (v *TOExtension) InsertQuery() string {
// // 	return `
// // INSERT INTO server_capability (
// //   name
// // )
// // VALUES (
// //   :name
// // )
// // RETURNING last_updated
// // `
// // }

// // func (v *TOServerCapability) SelectQuery() string {
// // 	return `
// // SELECT
// //   name,
// //   last_updated
// // FROM
// //   server_capability sc
// // `
// // }

// // func (v *TOServerCapability) DeleteQuery() string {
// // 	return `
// // DELETE FROM server_capability WHERE name=:name
// // `
// // }

// func (e *TOExtension) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
// 	return map[string]dbhelpers.WhereColumnInfo{
// 		"name": {"sc.name", nil},
// 	}
// }

// func (e TOExtension) GetKeyFieldsInfo() []api.KeyFieldInfo {
// 	return []api.KeyFieldInfo{{"name", api.GetStringKey}}
// }

// // Implementation of the Identifier, Validator interface functions
// func (e TOExtension) GetKeys() (map[string]interface{}, bool) {
// 	return map[string]interface{}{"name": v.Name}, true
// }

// func (e *TOExtension) SetKeys(keys map[string]interface{}) {
// 	//v.Name, _ = keys["name"].(string)
// }

// func (e *TOExtension) GetAuditName() string {
// 	return *e.Name
// }

// func (e *TOExtension) GetType() string {
// 	return "to extension"
// }

// func (e *TOExtension) Validate() error {
// 	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
// 	errs := validation.Errors{
// 		"name": validation.Validate(e.Name, validation.Required, rule),
// 	}
// 	return util.JoinErrs(tovalidate.ToErrors(errs))
// }

// func (e *TOExtension) Read() ([]interface{}, error, error, int) {
// 	return api.GenericRead(e)
// }

// func (e *TOExtension) Create() (error, error, int) {
// 	if e.APIInfo().User.UserName != "extension" {
// 		return errors.New("invalid user for this API. Only the \"extension\" user can use this"), nil, http.StatusForbidden
// 	}
// 	return nil, nil, http.StatusOK
// }

// func (e *TOExtension) Delete() (error, error, int) {
// 	if e.APIInfo().User.UserName != "extension" {
// 		return errors.New("invalid user for this API. Only the \"extension\" user can use this"), nil, http.StatusForbidden
// 	}
// 	return nil, nil, http.StatusOK
// }

// // if inf.User.UserName != "extension" {
// // 	api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("invalid user for this API. Only the \"extension\" user can use this"), nil)
// // 	return
// // }

// // CreateUpdateServercheck handles creating or updating an existing servercheck
// func CreateUpdateServercheck(w http.ResponseWriter, r *http.Request) {
// 	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
// 	if userErr != nil || sysErr != nil {
// 		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
// 		return
// 	}
// 	defer inf.Close()

// 	serverCheckReq := tc.ServercheckRequestNullable{}

// 	if err := api.Parse(r.Body, inf.Tx.Tx, &serverCheckReq); err != nil {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	id, exists, err := getServerID(serverCheckReq.ID, serverCheckReq.HostName, inf.Tx.Tx)
// 	if err != nil {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server id: "+err.Error()))
// 		return
// 	}
// 	if !exists {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("Server not found"), nil)
// 		return
// 	}

// 	col, exists, err := getColName(serverCheckReq.Name, inf.Tx.Tx)
// 	if err != nil {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting servercheck column name: "+err.Error()))
// 		return
// 	}
// 	if !exists {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("Server Check Extension %v not found - Do you need to install it?", *serverCheckReq.Name), nil)
// 		return
// 	}

// 	err = createUpdateServerCheck(id, col, *serverCheckReq.Value, inf.Tx.Tx)
// 	if err != nil {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("updating servercheck: "+err.Error()))
// 		return
// 	}

// 	successMsg := "Server Check was successfully updated"
// 	api.CreateChangeLogRawTx(api.ApiChange, successMsg, inf.User, inf.Tx.Tx)
// 	api.WriteRespAlert(w, r, tc.SuccessLevel, successMsg)
// }
