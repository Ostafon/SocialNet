package middlewares

import "net/http"

// Middleware для CORS
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешить все источники (для разработки)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		// Заголовки (с пробелами)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
		// Креды (на будущее)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Preflight (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
