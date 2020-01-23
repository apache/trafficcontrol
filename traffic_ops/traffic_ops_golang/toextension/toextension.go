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
