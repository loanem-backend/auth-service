package bcryptx

import "golang.org/x/crypto/bcrypt"

func Hash(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func Validate(raw, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}
