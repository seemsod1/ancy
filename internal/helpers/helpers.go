package helpers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func GetFileFormat(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1] // Last part is the file format
	}
	return "" // No file format detected
}

func GenerateHashedFileName(title string) (string, error) {
	finalTitle, err := bcrypt.GenerateFromPassword([]byte(title), bcrypt.DefaultCost)
	if err != nil {

		return "", err
	}

	return fmt.Sprintf("%x", finalTitle), nil
}
