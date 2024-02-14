package handler

import (
	"net/http"

	"github.com/evenboee/httputil/bind"
)

type ErrorFunc func(http.ResponseWriter, *http.Request, error)
type Func[T any] func(w http.ResponseWriter, r *http.Request, params T)

func NewFunc[T any](fn Func[T], opts ...funcConfigOpt) http.HandlerFunc {
	config := NewFuncConfig()
	for _, opt := range opts {
		opt(config)
	}

	return NewFuncWith(fn, config)
}

func NewFuncWith[T any](fn Func[T], config *FuncConfig) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		var req T
		err := bind.All(&req, r)
		if err != nil {
			config.BindError(w, r, err)
			return
		}

		fn(w, r, req)
	}

	if config.Wrapper != nil {
		f = config.Wrapper(f)
	}

	return f
}

type FuncConfig struct {
	BindError ErrorFunc
	Wrapper   WrapperFunc
}

type funcConfigOpt func(*FuncConfig)

func NewFuncConfig() *FuncConfig {
	return &FuncConfig{
		BindError: DefaultBindErrorFunc,
	}
}

func WithConfig(config *FuncConfig) funcConfigOpt {
	return func(c *FuncConfig) {
		c.BindError = config.BindError
	}
}

func (config *FuncConfig) WithWrapper(w WrapperFunc) *FuncConfig {
	config.Wrapper = w
	return config
}

func WithWrapper(w WrapperFunc) funcConfigOpt {
	return func(c *FuncConfig) {
		c.Wrapper = w
	}
}

func (config *FuncConfig) WithBindErrorFunc(fn ErrorFunc) *FuncConfig {
	config.BindError = fn
	return config
}

func WithBindErrorFunc(fn ErrorFunc) funcConfigOpt {
	return func(c *FuncConfig) {
		c.BindError = fn
	}
}

var DefaultBindErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, http.StatusBadRequest, err)
}
