package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET no está definido en el .env")
	}
	return []byte(secret), nil
}

func GenerateJWT(username string) (string, error) {
	secretKey, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["nombre_usuario"] = username
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error generando JWT: %w", err)
	}

	return tokenString, nil
}

func GetUsernameFromToken(tokenString string) (string, error) {
	secretKey, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token no válido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("no se pudieron leer los claims")
	}

	authorized, ok := claims["authorized"].(bool)
	if !ok || !authorized {
		return "", errors.New("usuario no autorizado")
	}

	username, ok := claims["nombre_usuario"].(string)
	if !ok {
		return "", errors.New("nombre_usuario inválido")
	}

	return username, nil
}
