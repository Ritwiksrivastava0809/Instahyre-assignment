package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPasswordArgon2(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s:%s", saltStr, hashStr), nil
}

func VerifyPassword(storedPassword string, providedPassword string) error {
	parts := strings.Split(storedPassword, ":")
	if len(parts) != 2 {
		return errors.New("invalid stored password format")
	}
	saltStr := parts[0]
	storedHashStr := parts[1]

	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return fmt.Errorf("error decoding salt: %w", err)
	}

	hashedPassword := hashPasswordWithSalt(providedPassword, salt)

	if hashedPassword != storedHashStr {
		return errors.New("password does not match")
	}

	return nil
}

func hashPasswordWithSalt(password string, salt []byte) string {
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	hashStr := base64.RawStdEncoding.EncodeToString(hash)

	return hashStr
}
