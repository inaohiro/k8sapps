package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// この middleware は、middeware が実行される時間の秒数が 03 のとき
// 503 エラーを返します
func IntentionalError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := fmt.Sprintf("%d", time.Now().Unix())

		if "03" == t[len(t)-2:] {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "service temporary unavailable. please wait a second}`))

			return
		}
		next.ServeHTTP(w, r)
	})
}
