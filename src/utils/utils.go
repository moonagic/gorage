package utils

import (
	"os"
	"fmt"
	"errors"
)

func CheckoutDir(dir string) error {
	if _, err := os.Stat(dir); err == nil {

	} else {
		fmt.Println("Dir not exists, try to create...", dir)
		err := os.MkdirAll(dir, 0711)
		if err != nil {
			fmt.Println("Error creating directory", dir)
			fmt.Println("err:", err)
			return errors.New("ERROR CREATING DIRECTORY")
		}
	}
	return nil
}