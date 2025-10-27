package helpers

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", nil
	}
	return string(hashedPassBytes), nil
}

func Compare(hashedPassBytes []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassBytes, []byte(password))
	return err == nil
}
