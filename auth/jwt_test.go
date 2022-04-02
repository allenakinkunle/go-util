package auth_test

import (
	"github.com/allenakinkunle/go-util/auth"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	signingKey         = "signingKey"
	audience           = "audience"
	issuer             = "issuer"
	accessTokenExpiry  = time.Minute * 10
	refreshTokenExpiry = time.Minute * 20
	subject            = "identifier"
)

func TestValidJWT(t *testing.T) {
	tokenManager := createValidTokenManager(t, signingKey, audience, issuer, accessTokenExpiry, refreshTokenExpiry)

	tokens, err := tokenManager.NewJWT(subject)
	require.NoError(t, err)

	accessToken, _, err := tokenManager.ValidateJWT(tokens.AccessToken)
	require.NoError(t, err)
	claim, ok := accessToken.Claims.(*auth.CustomClaim)
	require.True(t, ok)
	verifyClaim(t, claim, "accessToken")

	refreshToken, _, err := tokenManager.ValidateJWT(tokens.RefreshToken)
	require.NoError(t, err)
	claim, ok = refreshToken.Claims.(*auth.CustomClaim)
	require.True(t, ok)
	verifyClaim(t, claim, "refreshToken")
}

func TestInvalidJWT(t *testing.T) {
	tests := []struct {
		name                string
		invalidTokenManager auth.ITokenManager
		expectedError       error
	}{
		{
			name:                "invalid signingKey",
			invalidTokenManager: createValidTokenManager(t, "invalid_signing_key", audience, issuer, accessTokenExpiry, refreshTokenExpiry),
			expectedError:       auth.ErrInvalidSigningKey,
		},
		{
			name:                "invalid audience",
			invalidTokenManager: createValidTokenManager(t, signingKey, "invalid_audience", issuer, accessTokenExpiry, refreshTokenExpiry),
			expectedError:       auth.ErrInvalidJWTAudience,
		},
		{
			name:                "invalid issuer",
			invalidTokenManager: createValidTokenManager(t, signingKey, audience, "invalid_issuer", accessTokenExpiry, refreshTokenExpiry),
			expectedError:       auth.ErrInvalidJWTIssuer,
		},
		{
			name:                "expired token",
			invalidTokenManager: createValidTokenManager(t, signingKey, audience, "invalid_issuer", time.Minute-10, time.Minute-10),
			expectedError:       auth.ErrInvalidJWTIssuer,
		},
	}

	for _, tt := range tests {
		validTokenManager := createValidTokenManager(t, signingKey, audience, issuer, accessTokenExpiry, refreshTokenExpiry)
		t.Run(tt.name, func(t *testing.T) {
			invalidToken, _ := tt.invalidTokenManager.NewJWT(subject)
			_, _, err := validTokenManager.ValidateJWT(invalidToken.AccessToken)
			require.Error(t, err)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestInvalidSignatureMethod(t *testing.T) {
	validTokenManager := createValidTokenManager(t, signingKey, audience, issuer, accessTokenExpiry, refreshTokenExpiry)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "ES384",
			token: "eyJhbGciOiJFUzM4NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdWRpZW5jZSIsImV4cCI6MTcwMDAwMDA1MDAsImlhdCI6MTcwMDAwMDAxMDAsImlzcyI6Imlzc3VlciIsInN1YiI6Im5hbm9pZCIsInRva2VuVHlwZSI6ImFjY2Vzc1Rva2VuIn0.Ctu0kfMH3QLrLSqSTpBdM_8Ak7qSDFJGdrasQWeJUCXreQr1_VgHe5avVaYtOAjhpg8nSomuh4_WKcrDz1jBxStKVxSooVINuLzMplsjuZQX7OcHIQIMCKmFrvgX2q2Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := validTokenManager.ValidateJWT(tt.token)
			require.Error(t, err)
			require.Equal(t, auth.ErrInvalidSigningMethod, err)
		})
	}
}

func createValidTokenManager(t testing.TB, signingKey, audience, issuer string, accessTokenExpiry, refreshTokenExpiry time.Duration) auth.ITokenManager {
	t.Helper()
	tokenManager, _ := auth.NewTokenManager(signingKey, audience, issuer, accessTokenExpiry, refreshTokenExpiry)
	return tokenManager
}

func verifyClaim(t testing.TB, claim *auth.CustomClaim, tokenType string) {
	t.Helper()
	require.NoError(t, claim.Valid())
	require.Equal(t, claim.TokenType, tokenType)
	require.Equal(t, claim.Issuer, issuer)
	require.Equal(t, claim.Audience, audience)
	require.Equal(t, claim.Subject, subject)
}
