package middleware

import (
	"math/rand/v2"
	"net/http"
	"time"
)

// この middleware では以下の場合で遅延をはさみます
// - 毎 5 分ごと, 1 分間
// - 10% の確率
func Delay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, minute, _ := time.Now().Clock()

		// 毎 5 分ごとに, 1 分間遅延が発生する
		if minute%5 == 0 {
			time.Sleep(time.Duration(rand.IntN(999-300)+300) * time.Millisecond)
		}

		// 10% の確立でさらに遅延がはいる
		if rand.IntN(100) > (100 - 10) {
			time.Sleep(time.Duration(rand.IntN(999-500)+500) * time.Millisecond)
		}

		next.ServeHTTP(w, r)
	})
}
