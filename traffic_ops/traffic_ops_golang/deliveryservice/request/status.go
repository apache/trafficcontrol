
func GetStatusSingleton() api.Updater {
	return &deliveryServiceRequestStatus{}
}

// deliveryServiceRequestStatus implements interfaces needed to update the request status only
type deliveryServiceRequestStatus struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceRequestV15
}

func (req *deliveryServiceRequestStatus) GetAuditName() string {
	if req != nil && req.ID != nil {
		return strconv.Itoa(*req.ID)
	}
	return "UNKNOWN"
}

func (req *deliveryServiceRequestStatus) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

func (req *deliveryServiceRequestStatus) GetKeys() (map[string]interface{}, bool) {
	keys := map[string]interface{}{"id": 0}
	success := false
	if req.ID != nil {
		keys["id"] = *req.ID
		success = true
	}
	return keys, success
}

func (req *deliveryServiceRequestStatus) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int)
	req.ID = &i
}

func (*deliveryServiceRequestStatus) GetType() string {
	return "deliveryservice_request"
}

func (req *deliveryServiceRequestStatus) Update() (error, error, int) {
	// req represents the state the deliveryservice_request is to transition to
	// we want to limit what changes here -- only status can change,  and only according to the established rules
	// for status transition
	if req.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	var current tc.DeliveryServiceRequestV30
	err := req.APIInfo().Tx.QueryRowx(selectQuery+` WHERE r.id = $1`, *req.ID).StructScan(&current)
	if err != nil {
		return nil, errors.New("dsr status querying existing: " + err.Error()), http.StatusInternalServerError
	}

	if err = current.Status.ValidTransition(*req.Status); err != nil {
		return err, nil, http.StatusBadRequest // TODO verify err is secure to send to user
	}

	// keep everything else the same -- only update status
	st := req.Status
	req.DeliveryServiceRequestV15 = current.Downgrade()
	req.Status = st

	// LastEditedBy field should not change with status update

	if _, err = req.APIInfo().Tx.Tx.Exec(`UPDATE deliveryservice_request SET status = $1 WHERE id = $2`, *req.Status, *req.ID); err != nil {
		return api.ParseDBError(err)
	}

	if err = req.APIInfo().Tx.QueryRowx(selectQuery+` WHERE r.id = $1`, *req.ID).StructScan(req); err != nil {
		return nil, errors.New("dsr status update querying: " + err.Error()), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

// Validate is not needed when only Status is updated
func (req deliveryServiceRequestStatus) Validate() error {
	return nil
}

// ChangeLogMessage implements the api.ChangeLogger interface for a custom log message
func (req deliveryServiceRequestStatus) ChangeLogMessage(action string) (string, error) {
	XMLID := "UNKNOWN"
	if req.XMLID != nil {
		XMLID = *req.XMLID
	} else if req.DeliveryService != nil && req.DeliveryService.XMLID != nil {
		XMLID = *req.DeliveryService.XMLID
	}
	status := "UNKNOWN"
	if req.Status != nil {
		status = req.Status.String()
	}
	return fmt.Sprintf("Changed status of '%s' Delivery Service Request to '%s'", XMLID, status), nil
}
