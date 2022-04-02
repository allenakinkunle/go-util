package auth

import "golang.org/x/crypto/bcrypt"

type IPasswordHasher interface {
	Hash(password []byte) ([]byte, error)
	VerifyPassword(correctPassword, password []byte) error
}

type bcryptHasher struct {
	hashedPassword string
}

func NewBCryptHasher() *bcryptHasher {
	return &bcryptHasher{}
}

func (b bcryptHasher) Hash(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (b bcryptHasher) VerifyPassword(correctPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(correctPassword, password)
}

// mockHasher Mock of IPasswordHasher
type mockHasher struct{}

func NewMockHasher() mockHasher {
	return mockHasher{}
}

func (mockHasher) Hash(password []byte) ([]byte, error) {
	return []byte("hashed_password"), nil
}

func (h mockHasher) VerifyPassword(correctPassword, password []byte) error {
	return nil
}
