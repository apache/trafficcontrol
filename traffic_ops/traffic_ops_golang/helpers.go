package main

import "reflect"

// Extract the tag annotations from a struct into a string array
func ColsFromStructByTag(tagName string, thing interface{}) []string {
	cols := []string{}
	t := reflect.TypeOf(thing)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Get the field tag value
		tag := field.Tag.Get(tagName)
		cols = append(cols, tag)
	}
	return cols
}
