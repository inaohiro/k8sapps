package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	_ "embed"
)

//go:embed config.json
var configjson string

var (
	auth_url *url.URL
	app_url  *url.URL
)

type Route struct {
	Pattern string `json:"pattern"`
	Auth    bool   `json:"auth"`
}

func main() {
	var routes []Route
	if err := json.Unmarshal([]byte(configjson), &routes); err != nil {
		log.Fatal(err)
	}

	var err error
	_auth_url := os.Getenv("AUTH_URL")
	if _auth_url == "" {
		log.Fatal("AUTH_URL required")
	}
	auth_url, err = url.Parse(_auth_url)
	if err != nil {
		log.Fatal(err)
	}
	_app_url := os.Getenv("APP_URL")
	if _app_url == "" {
		log.Fatal("APP_URL required")
	}
	app_url, err = url.Parse(_app_url)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to convert HTTP port number from HTTP_PORT environment variable into int: %v", err))
	}

	mux := http.NewServeMux()
	for _, route := range routes {
		if route.Auth {
			mux.HandleFunc(route.Pattern, authHandler)
		} else {
			mux.HandleFunc(route.Pattern, handler)
		}
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

type VerifyTokenResponse struct {
	Namespace string `json:"namespace"`
	Error     string `json:"error"`
}

// token の検証, webapp へリクエスト送信
func authHandler(w http.ResponseWriter, r *http.Request) {

	// トークン検証
	auth_url.Path = path.Join(auth_url.Path, "tokens")
	defer func() {
		auth_url.Path = ""
	}()
	req, err := http.NewRequest(http.MethodGet, auth_url.String(), nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	// Authorization ヘッダが付いているはず
	// トークン検証 API につけて送信
	req.Header = r.Header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	defer resp.Body.Close()
	_body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	var verifyTokenResponse VerifyTokenResponse
	if err := json.Unmarshal(_body, &verifyTokenResponse); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	if resp.StatusCode >= 400 {
		writeJSON(w, resp.StatusCode, map[string]string{"message": verifyTokenResponse.Error})
		return
	}

	// トークン検証が成功したらアプリケーションにリクエスト送信
	app_url.Path = path.Join(app_url.Path, "api", verifyTokenResponse.Namespace, r.URL.Path)
	defer func() {
		app_url.Path = ""
	}()
	req, err = http.NewRequest(r.Method, app_url.String(), r.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	slog.Info("")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// token 発行
func handler(w http.ResponseWriter, r *http.Request) {

	// トークン発行
	auth_url.Path = path.Join(auth_url.Path, "tokens")
	defer func() {
		auth_url.Path = ""
	}()
	req, err := http.NewRequest(http.MethodPost, auth_url.String(), r.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	req.Header = r.Header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	// そのままレスポンスを返す
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
	}
}
