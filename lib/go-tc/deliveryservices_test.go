package tc

import (
	"encoding/json"
	"testing"
)

func ya (t *testing.T) {
	// make sure we get the right default when marshalled within a new delivery service
	var ds DeliveryService
	_, err := json.MarshalIndent(ds, ``, `  `)
	if err != nil {
		t.Errorf(err.Error())
	}
	c := ds.DeepCachingType.String()
	if c != "NEVER" {
		t.Errorf(`Default "%s" expected to be "NEVER"`, c)
	}

	// make sure null values are handled properly
	byt := []byte(`{"deepCachingType": null}`)
	err = json.Unmarshal(byt, &ds)
	if err != nil {
		t.Errorf(`Expected to be able to unmarshall a null deepCachingType as "NEVER", instead got error "%v"`, err.Error())
	}
	c = ds.DeepCachingType.String()
	if c != "NEVER" {
		t.Errorf(`null deepCachingType expected to be "NEVER", got "%s"`, c)
	}
}
