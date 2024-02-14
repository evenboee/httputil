package bind

import (
	"reflect"
)

// inspired by: https://github.com/gin-gonic/gin/blob/bb3519d26f52835cf00e5e430b52651a9c378c97/binding/form_mapping.go

var emptyField = reflect.StructField{}
var DefaultTagKey = "bind"

type getterFunc = func(string) ([]string, bool)

func Bind(obj any, f getterFunc) error {
	value := reflect.ValueOf(obj)

	_, err := bind(value, emptyField, f, DefaultTagKey)
	return err
}

func BindWith(obj any, f getterFunc, tagKey string) error {
	value := reflect.ValueOf(obj)

	_, err := bind(value, emptyField, f, tagKey)
	return err
}

func bind(value reflect.Value, field reflect.StructField, f getterFunc, tagKey string) (bool, error) {
	tagValue := field.Tag.Get(tagKey)
	if tagValue == "-" {
		return false, nil
	}

	vKind := value.Kind()

	if vKind == reflect.Pointer {
		isNew := false
		vPtr := value
		if value.IsNil() {
			vPtr = reflect.New(value.Type().Elem())
			isNew = true
		}

		isSet, err := bind(vPtr.Elem(), emptyField, f, tagKey)
		if err != nil {
			return false, err
		}
		if isNew && isSet {
			value.Set(vPtr)
		}
		return isSet, nil
	}

	if vKind != reflect.Struct || !field.Anonymous {
		ok, err := setValue(value, field, f, tagKey)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if vKind == reflect.Struct {
		tValue := value.Type()
		isSet := false
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			sField := tValue.Field(i)
			if sField.PkgPath != "" && !sField.Anonymous {
				continue
			}

			ok, err := bind(value.Field(i), sField, f, tagKey)
			if err != nil {
				return false, err
			}
			isSet = isSet || ok
		}
		return isSet, nil
	}

	return false, nil
}
