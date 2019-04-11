package api

import (
	"bufio"
	"errors"
	"reflect"
)

// allExportedFields returns true if ALL fields are exported.
// The input is expected to be a struct pointer.
func allFieldsExported(obj interface{}) bool {
	for i := 0; i < reflect.ValueOf(obj).Elem().NumField(); i++ {
		if !reflect.ValueOf(obj).Elem().Field(i).CanInterface() {
			return false
		}
	}
	return true
}

// AllFields is a key value for the `StructFields` map returned.
const AllFields = ""

// StructFields is a map from an "info" struct tag to a list of pointers.
// It can be obtained by using the `GetFieldInfo` function.
type StructFields map[string][]interface{}

// GetFieldInfo returns a list of pointers to fields in the struct.
//
// An error is returned for structs that have non-exported fields or if
// a non-pointer is passed. It is not an error to pass a pointer to a
// non-struct, but the pointer value will just be returned back wrapped
// in the StructFields type, which isn't useful. EDIT: I think the last
// sentence is no longer correct. TODO: This comment should be updated
// after I verify the behavior while testing.
//
// Example:
//
//	type Example struct {
//		ID int `info:"key"`
//		F1 int
//		F2 int
//	}
//
//	obj := Example{}
//	fields := GetFieldInfo(&obj)
//	*fields["key"][0] = 42
//
//	// Note: The original object is changed
//	fmt.Printf("The key is %d\n", obj.ID)
//
// Within the GetFieldInfo, there are no extra semantics for certain label
// values. This means that any "info" struct tag can be specified multiple
// times.
//
// The map also contains an entry `AllFields` to get a list of all the
// struct field pointers in order of appearance. This may be useful for
// database functions such as QueryRow:
//
//	QueryRow(GetQuery(), fields[AllFields]...)
//
func GetFieldInfo(obj interface{}) (StructFields, error) {

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return nil, errors.New("unpacking struct: expected pointer to struct")
	}

	val := reflect.ValueOf(obj).Elem()
	if val.Kind() != reflect.Struct {
		return nil, errors.New("unpacking struct: expected pointer to struct")
	}

	fields := StructFields{}

	// infoTag is captured by the append function, and
	// can be updated in the main loop.
	var infoTag string = ""

	// AppendToResult appends the passed argument (presumably a pointer)
	// to the list of fields. Any info tag will create a new key, value
	// pair in the map if it doesn't exist already.
	AppendToResult := func(ptr ...interface{}) StructFields {
		if infoTag != "" {
			fields[infoTag] = append(fields[infoTag], ptr...)
		}
		fields[AllFields] = append(fields[AllFields], ptr...)
		return fields
	}

	// if obj implements scanner, read obj as atomic value
	if _, ok := obj.(bufio.Scanner); ok {
		return AppendToResult(obj), nil
	}

	// Can't reflect on private data
	// Beyond this point, CanAddr can be used on all fields
	if !allFieldsExported(obj) {
		return nil, errors.New("unpacking struct: can't iterate over non-exported fields")
	}

	type Helper func(f reflect.Value) error
	var TraverseStruct Helper
	var HandlePointer Helper
	var GetAddressOfField Helper

	TraverseStruct = func(ptr reflect.Value) error {

		fieldMap, err := GetFieldInfo(ptr.Interface())
		if err != nil {
			return err
		}

		for infoTag, ptrs := range fieldMap {
			fields[infoTag] = append(fields[infoTag], ptrs...)
		}

		return nil
	}

	GetAddressOfField = func(field reflect.Value) error {
		AppendToResult(field.Addr().Interface())
		return nil
	}

	HandlePointer = func(ptr reflect.Value) error {
		switch ptr.Elem().Kind() {
		case reflect.Struct:
			return TraverseStruct(ptr)
		case reflect.Ptr:
			return errors.New("unpacking struct: double pointer not supported")
		default:
			return GetAddressOfField(ptr)
		}
	}

	Struct := reflect.TypeOf(obj).Elem()

	// Traversing Struct Fields
	for i := 0; i < val.NumField(); i++ {

		var err error
		field := val.Field(i)

		infoTag = Struct.Field(i).Tag.Get("info")

		switch field.Kind() {
		case reflect.Struct:
			err = TraverseStruct(field.Addr())
		case reflect.Ptr:
			err = HandlePointer(field)
		default:
			err = GetAddressOfField(field)
		}

		if err != nil {
			return nil, err
		}
	}

	return fields, nil
}

// MustGetFieldInfo is the version of GetFieldInfo that will
// panic if an error occurs.
func MustGetFieldInfo(obj interface{}) StructFields {
	ptrs, err := GetFieldInfo(obj)
	if err != nil {
		panic(err)
	}
	return ptrs
}
