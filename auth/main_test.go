package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestIssueToken_Success(t *testing.T) {
	body := `{"namespace":"test-ns"}`
	req := httptest.NewRequest(http.MethodPost, "/tokens", strings.NewReader(body))
	w := httptest.NewRecorder()

	issueToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tok Token
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if tok.Token == "" {
		t.Error("expected token in response")
	}
}

func TestIssueToken_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tokens", strings.NewReader("{invalid json"))
	w := httptest.NewRecorder()

	issueToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestVerifyToken_Success(t *testing.T) {
	// First, issue a token
	ns := "test-ns"
	token := getTokenForNamespace(t, ns)

	req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	verifyToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res["namespace"] != ns {
		t.Errorf("expected namespace %q, got %q", ns, res["namespace"])
	}
}

func TestVerifyToken_MissingAuthHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
	w := httptest.NewRecorder()

	verifyToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")
	w := httptest.NewRecorder()

	verifyToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestVerifyToken_MissingNamespaceClaim(t *testing.T) {
	// Create a token without "namespace" claim
	token := getTokenForNamespace(t, "")

	req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	verifyToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// Helper to issue a token for a namespace using the same signing method as main.go
func getTokenForNamespace(t *testing.T, ns string) string {
	claims := map[string]any{}
	if ns != "" {
		claims["namespace"] = ns
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return tokenString
}
