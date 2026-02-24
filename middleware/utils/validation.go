package utils

import (
	"fmt"
	"regexp"
	"unicode"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	// Simple email validation
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailPattern, email)
	if !matched {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePhone validates phone number format (simple check)
func ValidatePhone(phone string) error {
	if len(phone) < 7 || len(phone) > 15 {
		return fmt.Errorf("phone number must be between 7-15 digits")
	}
	// Check if contains only digits and optional +
	for _, ch := range phone {
		if !unicode.IsDigit(ch) && ch != '+' && ch != '-' {
			return fmt.Errorf("phone number contains invalid characters")
		}
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Check for complexity - at least 3 of: uppercase, lowercase, digit, special
	hasUpper, hasLower, hasDigit, hasSpecial := false, false, false, false

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case !unicode.IsLetter(ch) && !unicode.IsDigit(ch):
			hasSpecial = true
		}
	}

	// Count met criteria
	criteriaCount := 0
	if hasUpper {
		criteriaCount++
	}
	if hasLower {
		criteriaCount++
	}
	if hasDigit {
		criteriaCount++
	}
	if hasSpecial {
		criteriaCount++
	}

	if criteriaCount < 3 {
		return fmt.Errorf("password must contain at least 3 of: uppercase, lowercase, digit, special character")
	}

	return nil
}

// ValidateName validates user name
func ValidateName(name string) error {
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 100 {
		return fmt.Errorf("name must be less than 100 characters")
	}
	return nil
}

// ValidateLanguage validates language code
func ValidateLanguage(language string) error {
	validLanguages := map[string]bool{
		"km": true, // Khmer
		"en": true, // English
	}
	if !validLanguages[language] {
		return fmt.Errorf("invalid language. must be 'km' or 'en'")
	}
	return nil
}

// ValidateFarmName validates farm name
func ValidateFarmName(name string) error {
	if len(name) < 1 {
		return fmt.Errorf("farm name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("farm name must be less than 100 characters")
	}
	return nil
}

// ValidateDeviceName validates device name
func ValidateDeviceName(name string) error {
	if len(name) < 1 {
		return fmt.Errorf("device name cannot be empty")
	}
	if len(name) > 50 {
		return fmt.Errorf("device name must be less than 50 characters")
	}
	return nil
}

// ValidateRole validates user role
func ValidateRole(role string) error {
	validRoles := map[string]bool{
		"farmer": true,
		"viewer": true,
	}
	if !validRoles[role] {
		return fmt.Errorf("invalid role. must be 'farmer' or 'viewer'")
	}
	return nil
}
