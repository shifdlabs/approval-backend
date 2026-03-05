package helper

import (
	"Microservice/config"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ValidateTokenFormat(token string, publicKey string) (*jwt.Token, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	extractedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	return extractedToken, nil
}

// Utilize puiblic key to verify JWT and Extract Corresponding Payload
func ExtractIdentifierFromToken(tokenString string, pubKey *rsa.PublicKey) (*string, error) {
	var err error

	// 1) If caller didn't supply a key, load it now
	if pubKey == nil {
		pubKey, err = getPublicKey()
		if err != nil {
			return nil, fmt.Errorf("could not load public key: %w", err)
		}
	}

	// 2) Parse & verify signature
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}
	if !parsedToken.Valid {
		return nil, errors.New("validate: invalid token")
	}

	// 3) Extract the MapClaims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("validate: cannot convert claims to MapClaims")
	}

	// 4) Drill into "data" → "id"
	dataRaw, ok := claims["data"]
	if !ok {
		return nil, errors.New(`validate: "data" claim missing`)
	}
	dataMap, ok := dataRaw.(map[string]interface{})
	if !ok {
		return nil, errors.New(`validate: "data" is not an object`)
	}
	idRaw, ok := dataMap["id"]
	if !ok {
		return nil, errors.New(`validate: "id" missing in data claim`)
	}
	id, ok := idRaw.(string)
	if !ok {
		return nil, errors.New(`validate: "id" is not a string`)
	}

	return &id, nil
}

func GetAccessTokenPublicKey() string {
	env, _ := config.LoadConfig(".")
	return env.AccessTokenPublicKey
}

func GetRefreshTokenPublicKey() string {
	env, _ := config.LoadConfig(".")
	return env.RefreshTokenPublicKey
}

func getPublicKey() (*rsa.PublicKey, error) {
	b64 := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	pemBytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("unable to base64-decode public key: %w", err)
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public key PEM: %w", err)
	}
	return pub, nil
}
