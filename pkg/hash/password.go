package hash

import (
	"crypto/sha1"
	"fmt"
)

type SHA1Hash struct {
	salt string
}

func NewSHA1Hash(salt string) *SHA1Hash {
	return &SHA1Hash{
		salt: salt,
	}
}

func (h *SHA1Hash) Hash(password string) (string, error) {
	const op = "hash.password.Hash"

	hash := sha1.New()

	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
