package middleware

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"slices"
	"time"
)

// この middleware は、middeware が実行される時間の秒数が 00 ~ 05 の範囲にあるとき、
// 100ms ~ 999ms の遅延はさみます
func Delay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := fmt.Sprintf("%d", time.Now().Unix())

		if slices.Contains([]string{
			"00",
			"01",
			"02",
			"03",
			"04",
			"05",
		}, t[len(t)-2:]) {
			time.Sleep(time.Duration(rand.IntN(999-100)+100) * time.Millisecond)
		}
		next.ServeHTTP(w, r)
	})
}
