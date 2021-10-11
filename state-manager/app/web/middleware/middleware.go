package middleware

import (
	"github.com/compuzest/zlifecycle-state-manager/app/web/controllers"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"mime"
	"net/http"
	"time"
)

func EnforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				controllers.ErrorResponse(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				controllers.ErrorResponse(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		} else {
			controllers.ErrorResponse(w, "Content-Type header must be present", http.StatusBadRequest)
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
		start := time.Now()
		zlog.Logger.Printf(">> %s %s", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
		zlog.Logger.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				zlog.Logger.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
