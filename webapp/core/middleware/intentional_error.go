package middleware

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

// この middleware は、middeware が実行される時間の秒数が 0 で終わるとき
// 503 エラーを返します
// また、50% の確率で 503 エラーとなります
func IntentionalError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xskiperr := r.Header.Get("X-Skip-Error")
		if xskiperr != "" {
			next.ServeHTTP(w, r)
			return
		}

		t := fmt.Sprintf("%d", time.Now().Unix())

		if "0" == t[len(t)-1:] {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "seconds end with 0. please wait a second}`))

			return
		}

		if rand.IntN(100) > (100 - 50) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "50/50 error. please try again}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
