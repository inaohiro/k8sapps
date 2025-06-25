package controller

import (
	"context"
	"net/http"

	"k8soperation/token/internal/service"
)

type ctxNamespaceKey struct{}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		tokenString := auth[len("Bearer "):]

		namespace, err := service.ParseNamespaceFromToken(tokenString)
		if err != nil {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), ctxNamespaceKey{}, namespace)),
		)
	})
}

func GetNamespace(ctx context.Context) string {
	return ctx.Value(ctxNamespaceKey{}).(string)
}
