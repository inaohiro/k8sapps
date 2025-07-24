package middleware

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

// この middleware は、middeware が実行される時間の秒数が 00,10,20,30,40,50 のとき
// 503 エラーを返します
func IntentionalError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := fmt.Sprintf("%d", time.Now().Unix())

		if "0" == t[len(t)-1:] {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "seconds end with 0. please wait a second}`))

			return
		}

		n := rand.IntN(100)
		if n > 50 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "50/50 error. please try again}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
