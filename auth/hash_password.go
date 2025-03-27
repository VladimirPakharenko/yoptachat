package auth

import "golang.org/x/crypto/bcrypt"

// hashPassword хэширует пароль.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
