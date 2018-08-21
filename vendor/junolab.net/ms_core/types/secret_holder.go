package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))
var typeOfJson = reflect.TypeOf(json.RawMessage{})
var typeOfJsonPtr = reflect.TypeOf((*json.RawMessage)(nil))
var typeOfStringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

//ToSecuredString prepares struct for logging without password
func ToSecuredString(val interface{}) string {
	value := reflect.ValueOf(val)
	// return invalid values without analysis
	if !value.IsValid() {
		return fmt.Sprintf("%#v", val)
	}

	// get value from pointer
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = reflect.Indirect(value)
	}

	// check value type
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		if value.Len() > 0 {
			arrayVal := value.Index(0)
			if isSimpleKind(arrayVal.Kind()) {
				return fmt.Sprintf("%#v", value)
			}

			buf := []string{}
			for i := 0; i < value.Len(); i++ {
				buf = append(buf, ToSecuredString(value.Index(i).Interface()))
			}
			return fmt.Sprintf("%s", buf)
		} else {
			return value.Type().Name() + "{}"
		}
	case reflect.Struct:
	default:
		return fmt.Sprintf("%#v", value)
	}

	// analyze struct
	var result []string

	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		if fieldType.PkgPath != "" && !fieldType.Anonymous { // unexported
			continue // //https://tip.golang.org/doc/go1.6#reflect
		}

		fieldValue := value.Field(i)
		// return invalid values without analysis
		if !fieldValue.IsValid() {
			return fmt.Sprintf("%#v", val)
		}

		// get value from pointer
		if fieldValue.IsValid() && fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = reflect.Indirect(fieldValue)
		}

		sc := fieldType.Tag.Get("secured")
		name := fieldType.Name

		// check field type and secured tag and go deeper
		if sc != "true" {
			switch fieldValue.Kind() {
			case reflect.Struct:
				if fieldType.Type.Implements(typeOfStringer) {
					result = append(result, fmt.Sprintf("%s:%s", name, fieldValue.Interface()))
				} else {
					result = append(result, fmt.Sprintf("%s:%s", name, ToSecuredString(fieldValue.Interface())))
				}
				continue
			case reflect.Array, reflect.Slice:
				if fieldType.Type == typeOfBytes || fieldType.Type == typeOfJson || fieldType.Type == typeOfJsonPtr {
					result = append(result, fmt.Sprintf("%s:%s", name, fieldValue))
					continue
				}
				buf := []string{}
				for i := 0; i < fieldValue.Len(); i++ {
					buf = append(buf, ToSecuredString(fieldValue.Index(i).Interface()))
				}
				result = append(result, fmt.Sprintf("%s:[%s]", name, strings.Join(buf, ", ")))
				continue
			}
		}

		// obfuscate value
		if sc == "true" {
			result = append(result, fmt.Sprintf("%s:<secured>", name))
		} else {
			result = append(result, fmt.Sprintf("%s:%#v", name, fieldValue))
		}
	}

	return fmt.Sprintf("%s{%s}", value.Type().String(), strings.Join(result, ", "))
}

func isSimpleKind(k reflect.Kind) bool {
	return k != reflect.Array &&
		k != reflect.Chan &&
		k != reflect.Func &&
		k != reflect.Interface &&
		k != reflect.Map &&
		k != reflect.Ptr &&
		k != reflect.Slice &&
		k != reflect.Struct &&
		k != reflect.UnsafePointer
}
