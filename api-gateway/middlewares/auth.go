package middlewares

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"net/http"
	"os"
	"strings"
)

var publicPaths = []string{
	"/api/v1/auth/login",
	"/api/v1/auth/register",
	"/api/v1/auth/refresh",
	"/api/v1/auth/update-password",
	"/api/v1/auth/forgot-password",
	"/api/v1/auth/reset-password",
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range publicPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userId, _ := claims["sub"].(string)
		username, _ := claims["name"].(string)

		// 🔹 добавляем user-id в metadata (чтобы gRPC видел)
		md := metadata.New(map[string]string{
			"user-id":  userId,
			"username": username,
		})

		// создаём context с metadata
		ctx := metadata.NewOutgoingContext(r.Context(), md)

		// передаём дальше по цепочке
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
