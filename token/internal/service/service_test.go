package service

import (
	"k8soperation/token/internal/models"
	"testing"
)

func TestIssueAndParseNamespaceFromToken(t *testing.T) {
	namespace := "test-namespace"
	m := models.IssueToken{Namespace: namespace}

	token, err := IssueToken(m)
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	parsedNamespace, err := ParseNamespaceFromToken(token)
	if err != nil {
		t.Fatalf("ParseNamespaceFromToken failed: %v", err)
	}

	if parsedNamespace != namespace {
		t.Errorf("expected namespace %q, got %q", namespace, parsedNamespace)
	}
}
