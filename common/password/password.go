

package password

import (
	"bytes"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

// GetPassword gets password from user input
func GetPassword() ([]byte, error) {
	fmt.Printf("Password:")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return passwd, nil
}

// GetConfirmedPassword gets double confirmed password from user input
func GetConfirmedPassword() ([]byte, error) {
	fmt.Printf("Password:")
	first, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if 0 == len(first) {
		fmt.Println("User have to input password.")
		os.Exit(1)
	}

	fmt.Printf("Re-enter Password:")
	second, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if 0 == len(second) {
		fmt.Println("User have to input password.")
		os.Exit(1)
	}

	if !bytes.Equal(first, second) {
		fmt.Println("Unmatched Password")
		os.Exit(1)
	}
	return first, nil
}

// GetPassword gets node's wallet password from command line or user input
func GetAccountPassword() ([]byte, error) {
	var passwd []byte
	var err error
	passwd, err = GetPassword()
	if err != nil {
		return nil, err
	}

	return passwd, nil
}
