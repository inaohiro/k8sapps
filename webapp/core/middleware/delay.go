package middleware

import (
	"math/rand/v2"
	"net/http"
	"time"
)

// この middleware では以下の場合で遅延をはさみます
// - 毎 5 分ごと, 1 分間
// - 50%, 10% の確率
func Delay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, minute, _ := time.Now().Clock()

		// POST メソッドの場合、常に 300ms の遅延が発生する
		if r.Method == http.MethodPost {
			time.Sleep(time.Duration(300) * time.Millisecond)
		}

		// 5 分ごとに, 1 分間, 500ms ~ 1000ms の遅延が発生する
		if minute%5 == 0 {
			time.Sleep(time.Duration(rand.IntN(500)+500) * time.Millisecond)
		}

		// 50% の確率で 300ms ~ 500ms の遅延
		if rand.IntN(100) > (100 - 50) {
			time.Sleep(time.Duration(rand.IntN(200)+300) * time.Millisecond)

			// 10% の確立で 500ms ~ 1000ms の遅延
		} else if rand.IntN(100) > (100 - 10) {
			time.Sleep(time.Duration(rand.IntN(500)+500) * time.Millisecond)
		}

		next.ServeHTTP(w, r)
	})
}
