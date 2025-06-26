package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to convert HTTP port number from HTTP_PORT environment variable into int: %v", err))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tokens", issueToken)
	mux.HandleFunc("GET /tokens", verifyToken)

	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

type IssueTokenRequest struct {
	Namespace string `json:"namespace"`
}

type Token struct {
	Token string `json:"token"`
}

func issueToken(w http.ResponseWriter, r *http.Request) {
	var body IssueTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("不正な値がリクエストボディにセットされています")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"namespace": body.Namespace,
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		slog.Error("JWT トークンの署名に失敗しました")
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, Token{Token: tokenString})
}

const hmacSampleSecret = "hmacSampleSecret"

func verifyToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getTokenFromAuthorizationHeader(r)
	if err != nil {
		slog.Error(err.Error())
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(hmacSampleSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("不正なトークンが渡されました")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("不正なトークンが渡されました。claims の取得に失敗しました")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid token"})
		return
	}

	namespace, ok := claims["namespace"].(string)
	if !ok {
		slog.Error("不正なトークンが渡されました。namespace の取得に失敗しました")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"namespace": namespace})
}

func getTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("不正な値がAuthorization headerにセットされています。expected: Bearer ${token}. got: %s", auth)
	}
	return auth[len("Bearer "):], nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
	}
}
