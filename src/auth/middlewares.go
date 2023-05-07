package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
)

var JWT_SECRET = []byte("Avada Kedavra")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT token from authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return secret key for verification
			return JWT_SECRET, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Verify token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add claims to request context
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
	})
}

func CurrentUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Missing JWT claims", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok || username != chi.URLParam(r, "username") {
			http.Error(w, "unauthorized access", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Missing JWT claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		fmt.Println("Role: ", claims)
		if !ok || role != "admin" {
			http.Error(w, "Admins only", http.StatusForbidden)
			return
		}

		// Pass the request to the next middleware/handler in chain
		next.ServeHTTP(w, r)
	})
}
