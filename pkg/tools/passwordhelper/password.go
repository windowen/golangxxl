package passwordhelper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func GenPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	hashPassword := string(hash)
	return hashPassword
}

func VerifyPassword(hashPassword, plaintextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plaintextPassword))
	return err == nil
}

func ValidatePassword(ps string, min, max int) error {
	if len(ps) < min {
		return fmt.Errorf("password len is < %d", min)
	}
	if len(ps) > max {
		return fmt.Errorf("password len is > %d", max)
	}
	num := `[0-9]{1}`
	a_z := `(?i)[a-z]{1}`
	//A_Z := `[A-Z]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("password need a_z :%v", err)
	}
	//if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
	//	return fmt.Errorf("password need A_Z :%v", err)
	//}

	return nil
}
