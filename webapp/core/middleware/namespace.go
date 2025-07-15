package middleware

import (
	"k8soperation/core"
	"log/slog"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// path から namespace を取得
		if !strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/images") ||
			strings.HasPrefix(r.URL.Path, "/api/flavors") ||
			strings.HasPrefix(r.URL.Path, "/api/namespace") {
			next.ServeHTTP(w, r)
			return
		}
		namespace := strings.Split(r.URL.Path, "/")[2]

		// client-go 取得
		clientset, err := core.GetKubeClient()
		if err != nil {
			core.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		// namespace 作成
		nsName := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err = clientset.CoreV1().Namespaces().Create(r.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			slog.Error("failed to create namespace", slog.String("err", err.Error()))
		}

		next.ServeHTTP(w, r)
	})
}
