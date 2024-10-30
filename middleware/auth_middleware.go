package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("98991a0ce985e1e48f796693b95a33e52511cbc993e24ccbd0fdbc91cec73e21b37cfe5f50b35c002346fc06406104745b3c5f1426d81c632311be1a72e00743d044cb37af3024af41fb746a4daf908a5ba12c3f34805f18a4e8229026e70278916542bfe0475e6c7fd765c928988b29ce1d27876add3b284f04ef330ddb8b266fffa34191790000d0bc19a4ac1c28276363173a046e60015d7777bd3f04b682995f0a4f7c49eead65e1819a4e548108ae32b280228c3ec95827b8f0d9f645717afbe2d85412ebf1aa9cba7920c7a39c78bfe5f512ab493f006eb7205ab315cd509a6849fed1addb7c76bff2b821b17d91c5db2b7185f9f5644dfd66a9a5e0aa")

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user ID to context for access in other handlers
		ctx := context.WithValue(r.Context(), "email", claims["email"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
