package middleware

import (
	"mime"
	"net/http"
	"strings"
	"time"

	http2 "github.com/compuzest/zlifecycle-event-service/internal/web/http"
	"github.com/compuzest/zlifecycle-event-service/internal/zlog"
)

func EnforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http2.ErrorResponse(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http2.ErrorResponse(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		} else {
			http2.ErrorResponse(w, "Content-Type header must be present", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TimeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 60*time.Second, "Request timed out")
}

func LoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zlog.NewCtxEntry(r.Context())
		if strings.HasPrefix(r.URL.Path, "/health") {
			h.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		l.Debugf("REQUEST START %s %s", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
		l.Debugf("REQUEST END %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l := zlog.NewCtxEntry(r.Context())
				l.Errorf("panic: %+v", err)
				http2.ErrorResponse(w, "An unknown error occurred.", 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
