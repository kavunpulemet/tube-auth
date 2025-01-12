package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"auth/internal/utils"
)

func RecoveryMiddleware(ctx utils.MyContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v\n%s", err, debug.Stack())
				utils.NewErrorResponse(ctx, w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
