package helper

import (
	"Microservice/config"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	AccessTTL      = 1 * time.Minute
	RefreshTTL     = 365 * 24 * time.Hour
	RefreshRotateT = 30 * 24 * time.Hour
	Issuer         = "your.app.com"
)

// ─── Env‐based key loaders ───────────────────────────────────

func loadPEMEnv(base64Token string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64Token)
}

func getAccessPrivateKey() (*rsa.PrivateKey, error) {
	env, _ := config.LoadConfig(".")
	pem, err := loadPEMEnv(env.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(pem)
}

func getAccessPublicKey() (*rsa.PublicKey, error) {
	env, _ := config.LoadConfig(".")
	pem, err := loadPEMEnv(env.AccessTokenPublicKey)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(pem)
}

func getRefreshPrivateKey() (*rsa.PrivateKey, error) {
	env, _ := config.LoadConfig(".")
	pem, err := loadPEMEnv(env.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(pem)
}

func getRefreshPublicKey() (*rsa.PublicKey, error) {
	env, _ := config.LoadConfig(".")
	pem, err := loadPEMEnv(env.RefreshTokenPublicKey)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(pem)
}

// ─── Token generation ────────────────────────────────────────

func GenerateAccessToken(userID string) (string, error) {
	priv, err := getAccessPrivateKey()
	if err != nil {
		PrintValue(err, "Error Get Access Private Key")
		return "", err
	}
	claims := jwt.MapClaims{
		"iss":  Issuer,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(AccessTTL).Unix(),
		"data": map[string]string{"id": userID},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return tok.SignedString(priv)
}

func GenerateRefreshToken(userID string) (string, error) {
	priv, err := getRefreshPrivateKey()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"iss":  Issuer,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(RefreshTTL).Unix(),
		"data": map[string]string{"id": userID},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return tok.SignedString(priv)
}

// ─── Token validation & refresh ──────────────────────────────

func ValidateToken(accessToken string) *error {
	pubA, err := getAccessPublicKey()
	if err != nil {
		return &err
	}

	parsedA, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return pubA, nil
	})
	PrintObject(parsedA, "Parsed A")
	if err == nil && parsedA.Valid {
		// still valid
		return nil
	}

	return &err
}

func ExtractUserIDFromToken(accessToken string) (string, error) {
	// 1) Load your RSA public key
	pubA, err := getAccessPublicKey()
	if err != nil {
		return "", fmt.Errorf("could not load public key: %w", err)
	}

	// 2) Parse the token
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is RSA
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubA, nil
	})
	if err != nil {
		return "", fmt.Errorf("token parse error: %w", err)
	}

	// 3) Validate & extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// drill into "data"
		rawData, ok := claims["data"].(map[string]interface{})
		if !ok {
			return "", errors.New("missing or malformed 'data' claim")
		}
		// extract "id"
		rawID, ok := rawData["id"].(string)
		if !ok || rawID == "" {
			return "", errors.New("missing or invalid 'data.id' claim")
		}
		return rawID, nil
	}

	return "", errors.New("invalid token")
}

func ValidateOrRefreshAccess(accessToken, refreshToken string) (newAccess string, newRefresh string, err error) {
	// 1) Verify access
	pubA, err := getAccessPublicKey()
	if err != nil {
		return "", "", fmt.Errorf("load access pub key: %w", err)
	}
	parsedA, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return pubA, nil
	})
	PrintObject(parsedA, "Parsed A")
	if err == nil && parsedA.Valid {
		// still valid
		return accessToken, "", nil
	}

	// 2) If expired, validate refresh
	var verr *jwt.ValidationError
	PrintObject(errors.As(err, &verr), "Error Validate Or Refresh Access 1")
	PrintObject(jwt.ValidationErrorExpired, "Error Validate Or Refresh Access 2")
	if errors.As(err, &verr) && verr.Errors&jwt.ValidationErrorExpired != 0 {
		pubR, err := getRefreshPublicKey()
		if err != nil {
			return "", "", fmt.Errorf("load refresh pub key: %w", err)
		}
		PrintValue("Rezz 3", "Token")
		parsedR, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			}
			return pubR, nil
		})
		if err != nil || !parsedR.Valid {
			return "", "", fmt.Errorf("refresh token invalid: %w", err)
		}

		PrintValue("Rezz 4", "Token")

		// extract user ID
		claims := parsedR.Claims.(jwt.MapClaims)
		dataMap := claims["data"].(map[string]interface{})
		uid := dataMap["id"].(string)

		// issue new access
		newA, err := GenerateAccessToken(uid)
		if err != nil {
			return "", "", err
		}

		PrintValue("Rezz 5", "Token")

		// optionally rotate refresh
		exp := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Until(exp) < RefreshRotateT {
			newR, err := GenerateRefreshToken(uid)
			if err != nil {
				return "", "", err
			}
			return newA, newR, nil
		}
		return newA, "", nil
	}
	PrintValue("Rezz 7", "Token")

	return "", "", fmt.Errorf("access token invalid: %w", err)
}

func GetUserId(ctx *gin.Context) (*string, error) {
	token := ExtractToken(ctx)
	errId := ValidateToken(token)
	if errId != nil {
		return nil, *errId
	}

	userId, err := ExtractUserIDFromToken(token)
	if err != nil {
		return nil, err
	}
	PrintValue(userId, "User ID from token")
	return &userId, nil
}

func GetUserUUID(ctx *gin.Context) *uuid.UUID {
	token := ExtractToken(ctx)
	userId, err := ExtractUserIDFromToken(token)
	if err != nil {
		return nil
	}
	userUUID, err := uuid.FromString(userId)

	if err != nil {
		msg := "Failed parse to uuid from string"
		ErrorLog(err, 500, &msg)
	}
	return &userUUID
}

// func parseRSAPrivateKeyFromBase64(base64Key string) (*rsa.PrivateKey, error) {
// 	pemBytes, err := base64.StdEncoding.DecodeString(base64Key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
// }

// func parseRSAPublicKeyFromBase64(base64Key string) (*rsa.PublicKey, error) {
// 	pemBytes, err := base64.StdEncoding.DecodeString(base64Key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jwt.ParseRSAPublicKeyFromPEM(pemBytes)
// }
