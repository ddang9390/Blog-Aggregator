package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	jwt.RegisteredClaims
}

func jwtCreation(user User, secret string) (string, error) {
	claims := customClaims{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-aggregator",
			Subject:   fmt.Sprint(user.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func getToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		fmt.Println("authorization header is required")
		return "", fmt.Errorf("authorization header is required")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenString, nil
}

func jwtValidate(r *http.Request, secret string) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("authorization header is required")
		return "", fmt.Errorf("authorization header is required")
	}

	tokenString, err := getToken(r.Header)
	if err != nil {
		return "", fmt.Errorf("issue getting token")
	}
	claims := &customClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userId, nil
}
