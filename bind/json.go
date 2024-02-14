package bind

import (
	"encoding/json"
	"net/http"
)

type ContentType string

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
)

func JSON(obj any, r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func Body(obj any, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")

	switch ContentType(contentType) {
	case ContentTypeJSON:
		return JSON(obj, r)
	case ContentTypeForm:
		return PostForm(obj, r)
	// TODO: add more content types
	default:
		return JSON(obj, r)
	}
}
