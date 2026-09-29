package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super_secret_jwt_key_pt_anugrah_maha_tunggal_2026"
	}
	return []byte(secret)
}

func GenerateToken(userID uint, username, role, name string) (string, error) {
	expHours := 24
	if envExp := os.Getenv("JWT_EXPIRATION_HOURS"); envExp != "" {
		if val, err := strconv.Atoi(envExp); err == nil && val > 0 {
			expHours = val
		}
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "erp-pt-anugrah-maha-tunggal",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token otorisasi tidak valid")
}
