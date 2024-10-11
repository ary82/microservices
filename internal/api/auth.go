package api

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserType int32 `json:"user_type"`
	jwt.RegisteredClaims
}

type JwtInfo struct {
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
		Email:    email,
	}
	return userInfo, nil
}
