package middlewares

import "net/http"

// Middleware для CORS
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")

		// ОБЯЗАТЕЛЬНО добавить x-user-agent !!!
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Authorization, x-user-agent, X-Grpc-Web, grpc-timeout, ngrok-skip-browser-warning")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
