package bind

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type BodyUnmarshaler interface {
	UnmarshalBody([]byte) error
}

type FieldError struct {
	Field string
	Err   error
}

func (f *FieldError) Error() string {
	return fmt.Sprintf("FieldError: %s: %s", f.Field, f.Err)
}

func (f *FieldError) Unwrap() error {
	return f.Err
}

func All(obj any, r *http.Request) error {
	v := reflect.ValueOf(obj).Elem()

	if f := v.FieldByName("Path"); f.IsValid() {
		if err := Path(f.Addr().Interface(), r); err != nil {
			return &FieldError{"uri", err}
		}
	}

	if f := v.FieldByName("Query"); f.IsValid() {
		if err := Query(f.Addr().Interface(), r); err != nil {
			return &FieldError{"query", err}
		}
	}

	if f := v.FieldByName("PostForm"); f.IsValid() {
		if err := PostForm(f.Addr().Interface(), r); err != nil {
			return &FieldError{"postform", err}
		}
	}

	if f := v.FieldByName("Header"); f.IsValid() {
		if err := Header(f.Addr().Interface(), r); err != nil {
			return &FieldError{"header", err}
		}
	}

	if f := v.FieldByName("Body"); f.IsValid() {
		switch i := f.Addr().Interface().(type) {
		case BodyUnmarshaler:
			c, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return i.UnmarshalBody(c)
		}

		if err := Body(f.Addr().Interface(), r); err != nil {
			return &FieldError{"body", err}
		}
	}

	return nil
}
