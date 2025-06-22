package rest

import (
	"net/http"
	"time"

	"github.com/abcfe/abcfe-node/common/logger"
)

// LoggingMiddleware HTTP 요청 로깅 미들웨어
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 다음 핸들러 호출
		next.ServeHTTP(w, r)

		// 요청 로깅
		duration := time.Since(start)
		logger.Info("Request:", r.Method, r.URL.Path, "Duration:", duration)
	})
}

// RecoveryMiddleware 패닉 복구 미들웨어
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("API Panic recovered:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
