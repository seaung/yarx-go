package auth

import "golang.org/x/crypto/bcrypt"

// 使用bcrypt加密纯文本内容
func Encrypt(src string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// 比较密文和明文是否相同
func Compare(raw, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(raw), []byte(password))
}
