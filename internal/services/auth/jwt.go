package auth

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"

	"http-server/cmd/server/config"
)

type JWTClaim struct {
}

type JWT struct {
	Config *config.Config
}

type IJWT interface {
	GenerateToken(uuid string, exp int64) (string, error)
	RefreshToken(refreshToken string, exp int64) (string, error)
	ExtractIDFromToken(requestToken string) (string, error)
}

func NewJWT(config *config.Config) IJWT {
	return &JWT{
		Config: config,
	}
}

func (j *JWT) GenerateToken(uuid string, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid,
		"exp": exp,
	})

	tokenString, err := token.SignedString([]byte(j.Config.JWT_SECRET))

	return tokenString, err
}

func (j *JWT) RefreshToken(refreshToken string, exp int64) (string, error) {
	userId, err := j.ExtractIDFromToken(refreshToken)
	if err != nil {
		println("Cannot parse refresh token")
		return "", err
	}

	newToken, err := j.GenerateToken(userId, exp)
	if err != nil {
		println("Cannot gen refresh token")
		return "", err
	}

	return newToken, nil
}

func (j *JWT) ExtractIDFromToken(requestToken string) (string, error) {
	secretKey := j.Config.JWT_SECRET
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", err
	}
	return claims["sub"].(string), nil
}
