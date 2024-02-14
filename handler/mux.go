package handler

import (
	"log/slog"
	"net/http"
	"os"
)

type ServeMux struct {
	*http.ServeMux
	WrapperH WrapperHandler
	WrapperF WrapperFunc
	Logger   *slog.Logger
	LogLevel slog.Level
}

func New() *ServeMux {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return &ServeMux{
		ServeMux: http.NewServeMux(),
		WrapperH: WrapperH(DefaultRecoveryFunc),
		Logger:   logger,
		LogLevel: slog.LevelInfo,
		WrapperF: NopWrapper,
	}
}

func (mux *ServeMux) WithMux(m *http.ServeMux) *ServeMux {
	mux.ServeMux = m
	return mux
}

func (mux *ServeMux) WithWrapperH(w WrapperHandler) *ServeMux {
	mux.WrapperH = w
	return mux
}

func (mux *ServeMux) WithWrapperF(w WrapperFunc) *ServeMux {
	mux.WrapperF = w
	return mux
}

func (mux *ServeMux) WithLogger(l *slog.Logger) *ServeMux {
	mux.Logger = l
	return mux
}

func (mux *ServeMux) GetHandler() http.Handler {
	lh := SlogLoggerHandler(mux.Logger, slog.LevelInfo)
	w := WrapperH(lh)
	return w(mux.WrapperH(mux))
}

func (mux *ServeMux) Run(port string) error {
	mux.Logger.Info("Server started", "port", port)
	return http.ListenAndServe(":"+port, mux.GetHandler())
}

func (mux *ServeMux) MustRun(port string) {
	if err := mux.Run(port); err != nil {
		panic(err)
	}
}

func (mux *ServeMux) RunTLS(port string, certFile string, keyFile string) error {
	mux.Logger.Info("Server started", "port", port, "tls", "true")
	return http.ListenAndServeTLS(":"+port, certFile, keyFile, mux.GetHandler())
}

func (mux *ServeMux) MustRunTLS(port string, certFile string, keyFile string) {
	if err := mux.RunTLS(port, certFile, keyFile); err != nil {
		panic(err)
	}
}

func Handle[T any](mux *ServeMux, pattern string, f Func[T], opts ...funcConfigOpt) {
	mux.HandleFunc(pattern, mux.WrapperF(NewFunc(f, opts...)))
}

func HandleWith[T any](mux *ServeMux, pattern string, f Func[T], config *FuncConfig) {
	mux.HandleFunc(pattern, mux.WrapperF(NewFuncWith(f, config)))
}
