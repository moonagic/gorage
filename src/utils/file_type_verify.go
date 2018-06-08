package utils

import (
	"gorage/src/config"
	"strings"
)

// VerifyFileType Checkout file type
func VerifyFileType(fileName string) bool {
	validFileTypes := config.GetTypes()

	types := strings.Split(validFileTypes, ",")
	for _, value := range types {
		if strings.HasSuffix(fileName, "."+strings.Replace(value, " ", "", -1)) {
			return true
		}
	}
	return false
}
