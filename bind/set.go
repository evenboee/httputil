package bind

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// inspired by: https://github.com/gin-gonic/gin/blob/bb3519d26f52835cf00e5e430b52651a9c378c97/binding/form_mapping.go

type StringUnmarshaler interface {
	UnmarshalString(string) error
}

type TextUnmarshaler interface {
	encoding.TextUnmarshaler
}

type BindError struct {
	Field string
	Type  string
	Err   error
}

func (e *BindError) Error() string {
	return fmt.Sprintf("field %q (%s): %s", e.Field, e.Type, e.Err)
}

func (e *BindError) Unwrap() error {
	return e.Err
}

var (
	ErrUnknownType = errors.New("unknown type")

	// ErrConvertMapStringSlice can not convert to map[string][]string
	ErrConvertMapStringSlice = errors.New("can not convert to map slices of strings")

	// ErrConvertToMapString can not convert to map[string]string
	ErrConvertToMapString = errors.New("can not convert to map of strings")
)

var DefaultValueSeparator = " "

func setValue(value reflect.Value, field reflect.StructField, f getterFunc, tagKey string) (bool, error) {
	tagValue := field.Tag.Get(tagKey)
	if tagValue == "" {
		tagValue = field.Name
	}

	if tagValue == "" {
		return false, nil
	}

	tg := parseTag(tagValue)

	vs, ok := f(tg.Name)

	if !ok && tg.Default == "" {
		if tg.Required {
			return false, fmt.Errorf("missing required field: %q", tg.Name)
		}
		return false, nil
	}

	if tg.CommaSep && len(vs) != 0 {
		vs = strings.Split(vs[0], ",")
	}

	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			vs = strings.Split(tg.Default, DefaultValueSeparator)
		}
		return true, setSlice(vs, value, field)
	case reflect.Array:
		if !ok {
			vs = strings.Split(tg.Default, DefaultValueSeparator)
		}
		if len(vs) != value.Len() {
			return false, fmt.Errorf("%q is not a valid length for %s", vs, value.Type().String())
		}
		return true, setArray(vs, value, field)
	default:
		var val string
		if !ok {
			val = tg.Default
		}

		if len(vs) > 0 {
			val = vs[0]
		}

		err := setWithProperType(val, value, field)
		if err != nil {
			return true, &BindError{
				Field: tg.Name,
				Type:  value.Type().String(),
				Err:   err,
			}
		}

		return true, nil
	}
}

func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		case multipart.FileHeader:
			return nil
		}

		switch t := value.Addr().Interface().(type) {
		case StringUnmarshaler:
			return t.UnmarshalString(val)
		case TextUnmarshaler:
			return t.UnmarshalText([]byte(val))
		}

		return json.Unmarshal([]byte(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	case reflect.Pointer:
		if !value.Elem().IsValid() {
			value.Set(reflect.New(value.Type().Elem()))
		}

		return setWithProperType(val, value.Elem(), field)
	default:
		return ErrUnknownType
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func setFormMap(ptr any, form map[string][]string) error {
	el := reflect.TypeOf(ptr).Elem()

	if el.Kind() == reflect.Slice {
		ptrMap, ok := ptr.(map[string][]string)
		if !ok {
			return ErrConvertMapStringSlice
		}
		for k, v := range form {
			ptrMap[k] = v
		}

		return nil
	}

	ptrMap, ok := ptr.(map[string]string)
	if !ok {
		return ErrConvertToMapString
	}
	for k, v := range form {
		ptrMap[k] = v[len(v)-1] // pick last
	}

	return nil
}
