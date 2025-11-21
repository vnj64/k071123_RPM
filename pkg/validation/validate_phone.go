package validation

import (
	"errors"
	"regexp"
)

func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return errors.New("phone cannot be empty")
	}
	phoneRegex := regexp.MustCompile(`^7\d{10}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("phone must be in format 7********** (exactly 8 digits, starting with 7)")
	}
	return nil
}
