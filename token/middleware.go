package token

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ctxNamespaceKey struct{}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		tokenString := auth[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte("hmacSampleSecret"), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "token invalid", http.StatusUnauthorized)
			return
		}

		namespace, ok := claims["namespace"]
		if !ok {
			http.Error(w, "token invalid", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), ctxNamespaceKey{}, namespace)),
		)
	})
}

func getNamespace(ctx context.Context) string {
	return ctx.Value(ctxNamespaceKey{}).(string)
}
