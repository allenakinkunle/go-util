package auth

import "github.com/golang-jwt/jwt"

type mockTokenManager struct{}

func NewMockTokenManager() mockTokenManager {
	return mockTokenManager{}
}

func (mockTokenManager) NewJWT(subject string) (*Tokens, error) {
	return &Tokens{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil
}

func (m mockTokenManager) ValidateJWT(token string) (*jwt.Token, string, error) {
	return nil, "nanoid", nil
}

func (m mockTokenManager) GenerateRandomToken() string {
	return "randomToken"
}
