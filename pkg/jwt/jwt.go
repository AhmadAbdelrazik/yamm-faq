package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtService struct {
	key []byte
}

// New creates a new instance of JWT Service. JWT Service is used for
// generating session tokens for authentication.
func New(key string) *JwtService {
	return &JwtService{
		key: []byte(key),
	}
}

type UserClaims struct {
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}

// NewUserClaims Generates user claims to be used in generating JWT Tokens
func NewUserClaims(userId int, userType string) UserClaims {
	return UserClaims{
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yamm-faq-api",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			ID:        fmt.Sprint(userId),
		},
	}
}

// GenerateToken generate JWT Token with specific user ID and Type.
func (j *JwtService) GenerateToken(claims UserClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(j.key)
}

// VerifyToken
func (j *JwtService) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		return j.key, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("Unable to use user claims")
	}

	return claims, nil
}
