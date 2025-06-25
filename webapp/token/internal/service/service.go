package service

import (
	"fmt"
	"k8soperation/token/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const hmacSampleSecret = "hmacSampleSecret"

func IssueToken(m models.IssueToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"namespace": m.Namespace,
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", fmt.Errorf("signedString :%w", err)
	}

	return tokenString, nil
}

func ParseNamespaceFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(hmacSampleSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("token claims invalid")
	}

	namespace, ok := claims["namespace"].(string)
	if !ok {
		return "", fmt.Errorf("namespace claim missing or invalid")
	}

	return namespace, nil
}
