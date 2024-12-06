// Utility functions for generating passwords and salts
package utils

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
)

// Securely generate passphrase using diceware module
func GeneratePassword() string {
	list, err := diceware.Generate(4)
	if err != nil {
		Logger.Error().Err(err).Msg("Error generating passphrase")

		return ""
	}

	return strings.Join(list, "-")
}

func GenerateSalt() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	const length = 16

	salt := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := range salt {
		index, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			Logger.Error().Err(err).Msg("Error generating salt")
			return ""
		}
		salt[i] = charset[index.Int64()]
	}

	return string(salt)
}
