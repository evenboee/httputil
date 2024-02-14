package handler

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type ServeMux struct {
	*http.ServeMux
	// WrapperH is used to wrap handler on startup
	WrapperH WrapperHandler
	// WrapperF is used to wrap functions added with Handle and HandleWith
	WrapperF WrapperFunc
	// Logger is used to log server messages and for logging requests if LoggerFunc is not set
	Logger *slog.Logger
	// LogLevel is the level at which the logger logs at
	LogLevel slog.Level
	// LoggerFunc is used to log the request
	LoggerFunc LoggerFunc
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

func (mux *ServeMux) WithLogLevel(l slog.Level) *ServeMux {
	mux.LogLevel = l
	return mux
}

func (mux *ServeMux) WithLoggerFunc(l LoggerFunc) *ServeMux {
	mux.LoggerFunc = l
	return mux
}

func (mux *ServeMux) GetHandler() http.Handler {
	var loggerFunc WrapperFunc
	if loggerFunc != nil {
		loggerFunc = LogHandlerWith(mux.LoggerFunc)
	} else {
		loggerFunc = SlogLoggerHandler(mux.Logger, slog.LevelInfo)
	}

	w := WrapperH(loggerFunc)
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

func (mux *ServeMux) Register(pattern string, f http.HandlerFunc) {
	if pattern == "/" {
		pattern = "/{$}"
	} else if !strings.HasSuffix(pattern, "/") {
		pattern += "/{$}"
	}

	mux.HandleFunc(pattern, mux.WrapperF(f))
}

func Handle[T any](mux *ServeMux, pattern string, f Func[T], opts ...funcConfigOpt) {
	mux.Register(pattern, NewFunc(f, opts...))
}

func HandleWith[T any](mux *ServeMux, pattern string, f Func[T], config *FuncConfig) {
	mux.Register(pattern, NewFuncWith(f, config))
}
