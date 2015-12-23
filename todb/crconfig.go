package todb

import (
	"fmt"
)

type ContentRouter struct {
	profile  string `db:"profile"`
	apiport  int64  `db:"apiport"`
	location string `db:"location"`
	ip       string `db:"ip"`
	status   string `db:"status"`
	ip6      string `db:"ip6"`
	port     int64  `db:"port"`
	fqdn     string `db:"fqdn"`
	cdnname  string `db:"cdnname"`
}

func CRConfig(cdnName string) (interface{}, error) {

	crQuery := "select * from content_routers where cdnname=\"" + cdnName + "\""
	fmt.Println(crQuery)
	crs := []ContentRouter{}
	err := globalDB.Select(&crs, crQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return crs, nil
}
