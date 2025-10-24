package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", nil
	}
	return string(hashedPassBytes), nil
}

func comparePassword(hashedPassBytes []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassBytes, []byte(password))
	return err == nil
}
