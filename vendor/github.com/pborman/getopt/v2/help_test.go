package getopt

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type flagValue bool

func (f *flagValue) Set(value string, opt Option) error {
	switch strings.ToLower(value) {
	case "true", "t", "on", "1":
		*f = true
	case "false", "f", "off", "0":
		*f = false
	default:
		return fmt.Errorf("invalid flagValue %q", value)
	}
	return nil
}
func (f *flagValue) String() string {
	return fmt.Sprint(bool(*f))
}

func TestHelpDefaults(t *testing.T) {
	HelpColumn = 40
	set := New()
	bf := false
	bt := true
	set.FlagLong(&bf, "bool_false", 'f', "false bool value")
	set.FlagLong(&bt, "bool_true", 't', "true bool value")
	i := int(0)
	i8 := int8(0)
	i16 := int16(0)
	i32 := int32(0)
	i64 := int64(0)
	si := int(1)
	si8 := int8(8)
	si16 := int16(16)
	si32 := int32(32)
	si64 := int64(64)
	ui := uint(0)
	ui8 := uint8(0)
	ui16 := uint16(0)
	ui32 := uint32(0)
	ui64 := uint64(0)
	sui := uint(1)
	sui8 := uint8(8)
	sui16 := uint16(16)
	sui32 := uint32(32)
	sui64 := uint64(64)

	set.FlagLong(&i, "int", 0, "int value")
	set.FlagLong(&si, "int_set", 0, "set int value")
	set.FlagLong(&i8, "int8", 0, "int8 value")
	set.FlagLong(&si8, "int8_set", 0, "set int8 value")
	set.FlagLong(&i16, "int16", 0, "int16 value")
	set.FlagLong(&si16, "int16_set", 0, "set int16 value")
	set.FlagLong(&i32, "int32", 0, "int32 value")
	set.FlagLong(&si32, "int32_set", 0, "set int32 value")
	set.FlagLong(&i64, "int64", 0, "int64 value")
	set.FlagLong(&si64, "int64_set", 0, "set int64 value")

	set.FlagLong(&ui, "uint", 0, "uint value")
	set.FlagLong(&sui, "uint_set", 0, "set uint value")
	set.FlagLong(&ui8, "uint8", 0, "uint8 value")
	set.FlagLong(&sui8, "uint8_set", 0, "set uint8 value")
	set.FlagLong(&ui16, "uint16", 0, "uint16 value")
	set.FlagLong(&sui16, "uint16_set", 0, "set uint16 value")
	set.FlagLong(&ui32, "uint32", 0, "uint32 value")
	set.FlagLong(&sui32, "uint32_set", 0, "set uint32 value")
	set.FlagLong(&ui64, "uint64", 0, "uint64 value")
	set.FlagLong(&sui64, "uint64_set", 0, "set uint64 value")

	f32 := float32(0)
	f64 := float64(0)
	sf32 := float32(3.2)
	sf64 := float64(6.4)

	set.FlagLong(&f32, "float32", 0, "float32 value")
	set.FlagLong(&sf32, "float32_set", 0, "set float32 value")
	set.FlagLong(&f64, "float64", 0, "float64 value")
	set.FlagLong(&sf64, "float64_set", 0, "set float64 value")

	d := time.Duration(0)
	sd := time.Duration(time.Second)

	set.FlagLong(&d, "duration", 0, "duration value")
	set.FlagLong(&sd, "duration_set", 0, "set duration value")

	str := ""
	sstr := "string"
	set.FlagLong(&str, "string", 0, "string value")
	set.FlagLong(&sstr, "string_set", 0, "set string value")

	var fv flagValue
	set.FlagLong(&fv, "vbool", 0, "value bool").SetFlag()

	var fvo flagValue = true
	set.FlagLong(&fvo, "vbool_on", 0, "value bool").SetFlag()

	want := `
     --duration=value      duration value
     --duration_set=value  set duration value [1s]
 -f, --bool_false          false bool value
     --float32=value       float32 value
     --float32_set=value   set float32 value [3.2]
     --float64=value       float64 value
     --float64_set=value   set float64 value [6.4]
     --int=value           int value
     --int16=value         int16 value
     --int16_set=value     set int16 value [16]
     --int32=value         int32 value
     --int32_set=value     set int32 value [32]
     --int64=value         int64 value
     --int64_set=value     set int64 value [64]
     --int8=value          int8 value
     --int8_set=value      set int8 value [8]
     --int_set=value       set int value [1]
     --string=value        string value
     --string_set=value    set string value [string]
 -t, --bool_true           true bool value [true]
     --uint=value          uint value
     --uint16=value        uint16 value
     --uint16_set=value    set uint16 value [16]
     --uint32=value        uint32 value
     --uint32_set=value    set uint32 value [32]
     --uint64=value        uint64 value
     --uint64_set=value    set uint64 value [64]
     --uint8=value         uint8 value
     --uint8_set=value     set uint8 value [8]
     --uint_set=value      set uint value [1]
     --vbool               value bool
     --vbool_on            value bool [true]
`[1:]

	var buf bytes.Buffer
	set.PrintOptions(&buf)
	if got := buf.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
