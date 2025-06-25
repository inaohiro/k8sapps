package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"k8soperation/token/internal/service"
)

type ctxNamespaceKey struct{}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, fmt.Sprintf("不正な値がAuthorization headerにセットされていますｓ。expected: Bearer <token>, got: %s", auth), http.StatusBadRequest)
			return
		}
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
