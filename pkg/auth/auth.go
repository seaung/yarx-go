package auth

import "golang.org/x/crypto/bcrypt"

// 加密明文
func Encrypt(raw string) (string, error) {
	hkey, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)

	return string(hkey), err
}

// 比较明文与hash值
func CompareHash(hashPass, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(raw))
}
