package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"io/ioutil"
	"k071123/internal/services/user_service/domain/models"
	"k071123/internal/services/user_service/domain/services"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type fullClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func GenerateAuthToken(userUUID string, role permissions.Role, cfg services.Config) (string, error) {
	expHour, err := strconv.Atoi(cfg.AccessExpire())
	if err != nil {
		log.Printf("wrong expire time: %v", cfg.AccessExpire())
		return "", errs.ErrBadRequest
	}

	claims := fullClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userUUID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: string(role),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	privateKey, err := ParsePrivate(cfg.PrivatePemPath())
	if err != nil {
		log.Printf("failed to parse RSA private key: %v", err)
		return "", errs.NewErrorWithDetails(errs.ErrInternalServerError, fmt.Sprintf("failed to parse RSA private key: %s", err))
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("failed to sign token: %v", err.Error())
		return "", errs.NewErrorWithDetails(errs.ErrInternalServerError, fmt.Sprintf("failed to generate auth token: %s", err))
	}
	return tokenString, nil
}

// GenerateRefreshToken генерирует refresh_token для пользователя.
func GenerateRefreshToken(userUUID uuid.UUID, cfg services.Config) (*models.RefreshToken, error) {
	expHour, err := strconv.Atoi(cfg.RefreshExpire())
	if err != nil {
		return nil, errs.ErrBadRequest
	}
	expirationTime := time.Now().Add(time.Duration(expHour) * time.Hour)
	refreshToken := &models.RefreshToken{
		RefreshTokenUUID: uuid.New(),
		UserUUID:         userUUID,
		ExpiresAt:        expirationTime,
	}
	return refreshToken, nil
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func ReadPem(path string) (*pem.Block, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading key file: %w", err)
	}
	block, _ := pem.Decode(contents)
	if block == nil {
		return nil, fmt.Errorf("failed decoding key file as pem")
	}
	return block, nil
}

func ParsePublic(path string) (*rsa.PublicKey, error) {
	pubKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %v", err)
	}

	return pubKey, nil
}

func ParsePrivate(path string) (*rsa.PrivateKey, error) {
	block, err := ReadPem(path)
	if err != nil {
		return nil, err
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed decoding private key: %w", err)
	}
	return private.(*rsa.PrivateKey), nil
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func VerifyRefreshToken(tokenStr string, config services.Config) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AccessSecret()), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
