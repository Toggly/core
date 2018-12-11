package rest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

// CtxValue type
type CtxValue int

// CtxValue enum
const (
	CtxAPIVersion CtxValue = iota
	CtxValueOwner
	CtxValueRequestID
	CtxValueAuth
)

// Headers
const (
	XTogglyRequestID string = "X-Toggly-Request-Id"
	XTogglyOwnerID   string = "X-Toggly-Owner-Id"
	XServiceName     string = "X-Service-Name"
	XServiceVersion  string = "X-Service-Version"
)

// OwnerFromContext returns context value for project owner
func OwnerFromContext(r *http.Request) string {
	owner := r.Context().Value(CtxValueOwner)
	return owner.(string)
}

// VersionCtx adds api version to context
func VersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), CtxAPIVersion, version))
			next.ServeHTTP(w, r)
		})
	}
}

type wrappedWriter struct {
	Code         int
	Bytes        int
	Response     []byte
	statusSetted bool
	writer       http.ResponseWriter
	tee          io.Writer
	level        zerolog.Level
}

func (w *wrappedWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *wrappedWriter) Write(buf []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	if w.level == zerolog.DebugLevel {
		w.Response = append(w.Response, buf...)
	}
	n, err := w.writer.Write(buf)
	if w.tee != nil {
		_, err2 := w.tee.Write(buf[:n])
		if err == nil {
			err = err2
		}
	}
	w.Bytes += n
	return n, err
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	if !w.statusSetted {
		w.statusSetted = true
		w.Code = statusCode
		w.writer.WriteHeader(statusCode)
	}
}

// Logger middleware
func Logger(logger zerolog.Logger, level zerolog.Level) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := &wrappedWriter{writer: w, level: level}
			log := WithRequest(logger, r)
			defer func() {
				elapsed := time.Since(t1)
				var e *zerolog.Event
				switch ww.Code / 100 {
				case 5:
					e = log.Error()
				default:
					e = log.Info()
				}
				status := fmt.Sprintf("%d %s", ww.Code, http.StatusText(ww.Code))
				duration := fmt.Sprintf("%v", elapsed)
				e.Str("status", status).Int("bytes", ww.Bytes).Str("duration", duration)
				if level == zerolog.DebugLevel {
					log.Debug().Bytes("body", ww.Response).Msg("HTTP Response")
				}
				e.Msgf("<- %s %s", r.Method, r.RequestURI)
			}()
			log.Info().Str("proto", r.Proto).Str("remote", r.RemoteAddr).Msgf("-> %s %s", r.Method, r.RequestURI)
			if level == zerolog.DebugLevel {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Error().Err(err).Msg("Can't parse request body")
				}
				log.Debug().Bytes("body", body).Msg("HTTP Request")
			}
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// WithRequest looger with request context
func WithRequest(l zerolog.Logger, r *http.Request) zerolog.Logger {
	reqID, ok := r.Context().Value(CtxValueRequestID).(string)
	if !ok || reqID == "" {
		return l
	}
	return l.With().Str("req", reqID).Logger()
}

// RequestIDCtx adds request id to context
func RequestIDCtx(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(http.CanonicalHeaderKey(XTogglyRequestID))
			warn := false
			if rid == "" {
				rid = fmt.Sprintf("%d", middleware.NextRequestID())
				warn = true
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxValueRequestID, rid)
			w.Header().Set(http.CanonicalHeaderKey(XTogglyRequestID), rid)
			req := r.WithContext(ctx)
			if warn {
				log := WithRequest(logger, req)
				log.Warn().Msg("Header Toggly-Request-Id missed. Autogenerated.")
			}
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(fn)
	}
}

// OwnerCtx adds auth data to context
func OwnerCtx(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			owner := r.Header.Get(http.CanonicalHeaderKey(XTogglyOwnerID))
			if owner == "" {
				logger.Warn().Msg("Header X-Toggly-Owner-Id missed")
				NotFoundResponse(w, r, "Owner not found")
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxValueOwner, owner)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// ServiceInfo adds service information to the response header
func ServiceInfo(name string, version string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(http.CanonicalHeaderKey(XServiceName), name)
			w.Header().Set(http.CanonicalHeaderKey(XServiceVersion), version)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}