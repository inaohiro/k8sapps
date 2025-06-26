package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_TokenIssue(t *testing.T) {
	// Mock auth server for token issuance
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/tokens") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(201)
		w.Write([]byte(`{"token":"abc123"}`))
	}))
	defer authSrv.Close()

	auth_url, _ = url.Parse(authSrv.URL)
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("abc123")) {
		t.Errorf("expected token in response, got %s", string(body))
	}
}

func TestAuthHandler_Success(t *testing.T) {
	// Mock auth server for token verification
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/tokens") {
			w.WriteHeader(200)
			w.Write([]byte(`{"namespace":"test-ns"}`))
		}
	}))
	defer authSrv.Close()

	// Mock app server
	appSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/test-ns/") {
			w.WriteHeader(200)
			w.Write([]byte(`{"ok":true}`))
		}
	}))
	defer appSrv.Close()

	auth_url, _ = url.Parse(authSrv.URL)
	app_url, _ = url.Parse(appSrv.URL)

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	req.Header.Set("Authorization", "Bearer test")
	w := httptest.NewRecorder()
	authHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("ok")) {
		t.Errorf("expected ok in response, got %s", string(body))
	}
}

func TestAuthHandler_TokenVerifyFail(t *testing.T) {
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"invalid token"}`))
	}))
	defer authSrv.Close()

	auth_url, _ = url.Parse(authSrv.URL)
	app_url, _ = url.Parse("http://localhost") // not used

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	w := httptest.NewRecorder()
	authHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("invalid token")) {
		t.Errorf("expected error in response, got %s", string(body))
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, 418, map[string]string{"hello": "world"})
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != 418 {
		t.Errorf("expected 418, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("expected application/json, got %s", ct)
	}
	var m map[string]string
	json.NewDecoder(resp.Body).Decode(&m)
	if m["hello"] != "world" {
		t.Errorf("expected hello=world, got %v", m)
	}
}
