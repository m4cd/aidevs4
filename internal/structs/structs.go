package structs

import (
	"fmt"
	"reflect"
)

// func PrintStruct(s any) {
// 	v := reflect.ValueOf(s)
// 	t := reflect.TypeOf(s)

// 	// dereference pointer if needed
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 		t = t.Elem()
// 	}

// 	fmt.Printf("[+] %s\n", t.Name())
// 	for i := 0; i < v.NumField(); i++ {
// 		fmt.Printf("%s: %v\n", t.Field(i).Name, v.Field(i).Interface())
// 	}
// 	fmt.Println()
// }

func PrintStruct(s any) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// dereference pointer if needed
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	fmt.Printf("[+] %s\n", t.Name())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			fmt.Printf("%s:\n", fieldName)
			for j := 0; j < field.Len(); j++ {
				fmt.Printf("  [%d]: %v\n", j, field.Index(j).Interface())
			}
		} else {
			fmt.Printf("%s: %v\n", fieldName, field.Interface())
		}
	}
	fmt.Println()
}
