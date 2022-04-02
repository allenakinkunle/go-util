package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	ErrJWTSigningKeyNotSet  = errors.New("jwt signing key not set")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrExpiredJWTToken      = errors.New("expired JWT token")
	ErrInvalidJWTAudience   = errors.New("invalid jwt audience")
	ErrInvalidJWTIssuer     = errors.New("invalid jwt issuer")
	ErrInvalidSigningKey    = errors.New("invalid jwt signature")
)

type ITokenManager interface {
	NewJWT(subject string) (*Tokens, error)
	ValidateJWT(token string) (*jwt.Token, string, error)
	GenerateRandomToken() string
}

type CustomClaim struct {
	*jwt.StandardClaims
	TokenType string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type tokenManager struct {
	SigningKey               string
	Audience                 string
	Issuer                   string
	AccessTokenTimeToExpiry  time.Duration
	RefreshTokenTimeToExpiry time.Duration
}

func NewTokenManager(signingKey string, audience string, issuer string,
	accessTokenTimeToExpiry, refreshTokenTimeToExpiry time.Duration) (*tokenManager, error) {
	if signingKey == "" {
		return nil, ErrJWTSigningKeyNotSet
	}

	return &tokenManager{
		SigningKey:               signingKey,
		Audience:                 audience,
		Issuer:                   issuer,
		AccessTokenTimeToExpiry:  accessTokenTimeToExpiry,
		RefreshTokenTimeToExpiry: refreshTokenTimeToExpiry,
	}, nil
}

func (t *tokenManager) NewJWT(subject string) (*Tokens, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaim{
		TokenType: "accessToken",
		StandardClaims: &jwt.StandardClaims{
			Audience:  t.Audience,
			Issuer:    t.Issuer,
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().UTC().Add(t.AccessTokenTimeToExpiry).Unix(),
			Subject:   subject,
		},
	})
	accessTokenString, err := accessToken.SignedString([]byte(t.SigningKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaim{
		TokenType: "refreshToken",
		StandardClaims: &jwt.StandardClaims{
			Audience:  t.Audience,
			Issuer:    t.Issuer,
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().UTC().Add(t.RefreshTokenTimeToExpiry).Unix(),
			Subject:   subject,
		},
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(t.SigningKey))
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	return tokens, nil
}

func (t *tokenManager) ValidateJWT(jwtToken string) (*jwt.Token, string, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		if isAudienceOK := token.Claims.(*CustomClaim).VerifyAudience(t.Audience, false); !isAudienceOK {
			return nil, ErrInvalidJWTAudience
		}

		if isIssuerOK := token.Claims.(*CustomClaim).VerifyIssuer(t.Issuer, false); !isIssuerOK {
			return nil, ErrInvalidJWTIssuer
		}

		return []byte(t.SigningKey), nil
	})

	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorSignatureInvalid:
			return nil, "", ErrInvalidSigningKey
		case jwt.ValidationErrorExpired:
			return nil, "", ErrExpiredJWTToken
		}
		if err.(*jwt.ValidationError).Inner != nil {
			return nil, "", err.(*jwt.ValidationError).Inner
		}
		return nil, "", err
	}

	subject := token.Claims.(*CustomClaim).Subject

	return token, subject, nil
}

func (t *tokenManager) GenerateRandomToken() string {
	token := make([]byte, 5)
	rand.Read(token)

	hasher := sha1.New()
	hasher.Write(token)
	tokenString := fmt.Sprintf("%x", hasher.Sum(nil))

	return tokenString
}
