package clistruct

import (
	"reflect"
	"strings"
)

func getStructField(v interface{}, fieldName string) (reflect.Value, error) {
	field := indirectValue(reflect.ValueOf(v)).
		FieldByName(fieldName)

	if !field.IsValid() || !field.CanSet() {
		return reflect.Value{}, NewErrInvalid(v)
	}

	return field, nil
}

func setStructField(v interface{}, fieldName string, value interface{}) error {
	field, err := getStructField(v, fieldName)
	if err != nil {
		return err
	}

	reflectValue := reflect.ValueOf(value)

	if field.Type() != reflectValue.Type() {
		return NewErrTypeMistmatch(
			field.Type().String(),
			reflectValue.Type().String(),
		)
	}

	field.Set(reflectValue)

	return nil
}

func checkValue(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return NewErrPtrRequired(v)
	}

	return nil
}

func shouldBeStruct(reflectValue reflect.Value) error {
	if reflectValue.Kind() != reflect.Struct {
		return NewErrInvalidKind(
			reflect.Struct,
			reflectValue.Kind(),
		)
	}

	return nil
}

func isStructFieldExported(field reflect.StructField) bool {
	// From reflect docs:
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	return field.PkgPath == ""
}

func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func typeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func getStructFieldTag(field reflect.StructField, name string) string {
	return strings.TrimSpace(field.Tag.Get(name))
}

func getStructFieldTagSlice(field reflect.StructField, name string) []string {
	tagValue := strings.TrimSpace(field.Tag.Get(name))
	if tagValue == "" {
		return nil
	}

	if strings.HasPrefix(tagValue, "[") && strings.HasSuffix(tagValue, "]"){
		tagValue = tagValue[1:len(tagValue)-1]
		return cleanTagValues(strings.Split(tagValue, ","))
	} else {
		return []string{tagValue}
	}
}

func cleanTagValues(values []string) []string {
	res := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.Trim(v, "'")
		v = strings.TrimLeft(v, "'")
		res = append(res, v)
	}
	return res
}
