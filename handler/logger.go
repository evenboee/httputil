package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type LoggerFuncParams struct {
	Request *http.Request

	StatusCode int
	Method     string
	Path       string
	Query      string

	Timestamp time.Time
	Latency   time.Duration
}

type LoggerFunc func(LoggerFuncParams)

var DefaultLoggerFormatterTimeFormat = "2006-01-02 15:04:05"

var DefaultLoggerFormatter LoggerFunc = func(p LoggerFuncParams) {
	path := p.Path
	if p.Query != "" {
		path = path + "?" + p.Query
	}

	fmt.Printf("%s | %12s | %3d | %-7s | %s\n",
		p.Timestamp.Format(DefaultLoggerFormatterTimeFormat),
		p.Latency.String(), p.StatusCode, p.Method, path)
}

func SlogLogger(l *slog.Logger, level slog.Level) LoggerFunc {
	return func(p LoggerFuncParams) {
		path := p.Path
		if p.Query != "" {
			path = path + "?" + p.Query
		}

		l.Log(p.Request.Context(), level, "Request",
			"time", p.Timestamp.Format(DefaultLoggerFormatterTimeFormat),
			"latency", p.Latency.String(),
			"status", p.StatusCode,
			"method", p.Method,
			"path", path,
		)
	}
}

func SlogLoggerHandler(l *slog.Logger, level slog.Level) WrapperFunc {
	return LogHandlerWith(SlogLogger(l, level))
}

func LogHandlerWith(f LoggerFunc) WrapperFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startT := time.Now()
			lw := NewLogWriter(w)
			next(lw, r)
			elapsedT := time.Since(startT)

			status := lw.Status
			if status == 0 {
				status = http.StatusOK
			}

			f(LoggerFuncParams{
				Request:    r,
				StatusCode: lw.Status,
				Method:     r.Method,
				Path:       r.URL.Path,
				Query:      r.URL.RawQuery,
				Timestamp:  startT,
				Latency:    elapsedT,
			})
		}
	}
}

var LoggerHandler = LogHandlerWith(DefaultLoggerFormatter)

type LogWriter struct {
	Status int
	http.ResponseWriter
	Output io.Writer
}

func NewLogWriter(w http.ResponseWriter) *LogWriter {
	return &LogWriter{
		Status:         0,
		ResponseWriter: w,
		Output:         os.Stdout,
	}
}

func (w *LogWriter) WriteHeader(code int) {
	if w.Status != 0 {
		if w.Status != code {
			w.Output.Write(
				[]byte(fmt.Sprintf(
					"[ warning ] WriteHeader called with %d but status was %d\n",
					code, w.Status),
				),
			)
		}
		return
	}

	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = http.StatusOK
	}

	return w.ResponseWriter.Write(b)
}
