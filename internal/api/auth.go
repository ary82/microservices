package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserType int32  `json:"user_type"`
	Id       string `json:"id"`
	jwt.RegisteredClaims
}

type JwtInfo struct {
	Id       string
	Email    string
	UserType int32
}

func ParseJWT(tokenStr string) (*JwtInfo, error) {
	// Parse claims to get the jwt token
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Check validity of claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims corrupted")
	}

	// Get email from claims
	email, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	// User struct for use in context
	userInfo := &JwtInfo{
		UserType: claims.UserType,
		Id:       claims.Id,
		Email:    email,
	}
	return userInfo, nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		key := TokenCtxKey("token")
		ctx := context.WithValue(r.Context(), key, header)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type TokenCtxKey string
