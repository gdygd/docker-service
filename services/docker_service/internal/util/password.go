package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword는 bcrypt를 사용하여 비밀번호를 해시화
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword는 비밀번호와 해시를 비교
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
