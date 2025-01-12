package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
)

type JwtInterface interface {
	Create(userID uuid.UUID, roleName string) (string, error)
	Decode(tokenString string, claims *Claims) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserID   uuid.UUID `json:"user_id"`
	RoleName string    `json:"role_name"`
	Email    string    `json:"email"`
}

type JwtManagerInterface interface {
	CreateManager(userID int, email string) (string, error)
	DecodeManager(tokenString string, claims *ClaimsManager) error
}

type ClaimsManager struct {
	jwt.RegisteredClaims
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

type JwtStruct struct {
	SecretKey   string
	ExpiredTime time.Duration
}

var Jwt = getJwt()

func getJwt() JwtInterface {
	return &JwtStruct{
		SecretKey:   env.AppEnv.JwtSecretKey,
		ExpiredTime: env.AppEnv.JwtExpTime,
	}
}

var JwtManager = getManagerJwt()

func getManagerJwt() JwtManagerInterface {
	return &JwtStruct{
		SecretKey:   env.AppEnv.JwtSecretKey,
		ExpiredTime: env.AppEnv.JwtExpTime,
	}
}

func (j *JwtStruct) Create(userID uuid.UUID, roleName string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gogo-manager",
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{"gogo-manager"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        userID.String(),
		},
		UserID:   userID,
		RoleName: roleName,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *JwtStruct) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func (j *JwtStruct) CreateManager(userID int, email string) (string, error) {
	claims := ClaimsManager{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gogo-manager",
			Subject:   email,
			Audience:  jwt.ClaimStrings{"gogo-manager"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(userID),
		},
		UserID: userID,
		Email: email,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *JwtStruct) DecodeManager(tokenString string, claims *ClaimsManager) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
