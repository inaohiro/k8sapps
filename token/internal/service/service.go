package service

import (
	"fmt"
	"k8soperation/token/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func IssueToken(m models.IssueToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"namespace": m.Namespace,
	})

	tokenString, err := token.SignedString("hmacSampleSecret")
	if err != nil {
		return "", fmt.Errorf("signedString :%w", err)
	}

	return tokenString, nil
}
