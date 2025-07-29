package middleware

import (
	"math/rand/v2"
	"net/http"
)

func IntentionalError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-Error がついていればわざとエラーにする
		// 用意している frontend は X-Error をつけないため、エラーにならない
		// k6 からのリクエストのみを意図的なエラーの対象とする
		xerr := r.Header.Get("X-Error")
		if xerr == "" {
			next.ServeHTTP(w, r)
			return
		}

		// 15% の確率でエラーとする
		if rand.IntN(100) > (100 - 15) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "15% error. please try again}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
