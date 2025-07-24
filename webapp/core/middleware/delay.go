package middleware

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

// この middleware は、middeware が実行される時間の秒数が 5 で終わるとき
// 100ms ~ 999ms の遅延はさみます
// また、50% の確率で 100ms = 999ms の遅延をはさみます
func Delay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := fmt.Sprintf("%d", time.Now().Unix())

		if "5" == t[len(t)-1:] {
			time.Sleep(time.Duration(rand.IntN(999-100)+100) * time.Millisecond)
		}

		if rand.IntN(100) > (100 - 50) {
			time.Sleep(time.Duration(rand.IntN(999-100)+100) * time.Millisecond)
		}

		next.ServeHTTP(w, r)
	})
}
