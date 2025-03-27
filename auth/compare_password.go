package auth

import "golang.org/x/crypto/bcrypt"

// comparePasswords сравнивает хэшированный пароль с обычным.
func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
