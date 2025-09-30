package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func ZeroLogLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// method := r.Method
		// requestUri := r.RequestURI
		// remoteAddr := r.RemoteAddr
		// status := ww.Status()
		// bytesWritten := ww.BytesWritten()
		// duration := time.Since(start)

		defer func() {
			log.Info().
				Str("method", r.Method).
				Str("path", r.RequestURI).
				Str("remote_addr", r.RemoteAddr).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", time.Since(start)).
				Msg("Request")
		}()

		next.ServeHTTP(ww, r)
	})
}
