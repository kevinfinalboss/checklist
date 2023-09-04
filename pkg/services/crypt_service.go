package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var encodedSecretKey = "SMxP8/GDSrfMjrBxaalq7aIxknlUPwZb8TSFCy7uVZY="
var secretKey []byte

func init() {
	var err error
	secretKey, err = base64.StdEncoding.DecodeString(encodedSecretKey)
	if err != nil {
		fmt.Println("Error decoding secret key:", err)
	}
}

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(text))
	stream := cipher.NewCTR(block, secretKey[:block.BlockSize()])
	stream.XORKeyStream(ciphertext, []byte(text))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encryptedText string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	decodedText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	decryptedText := make([]byte, len(decodedText))
	stream := cipher.NewCTR(block, secretKey[:block.BlockSize()])
	stream.XORKeyStream(decryptedText, decodedText)

	return string(decryptedText), nil
}
