package helpers

import "golang.org/x/crypto/bcrypt"

func GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string,givenPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(givenPassword))
	return err
}