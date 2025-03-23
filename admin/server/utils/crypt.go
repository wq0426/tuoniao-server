// @program:     flashbear
// @file:        crypt.go
// @author:      ac
// @create:      2024-10-20 14:56
// @description:
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// AesEncrypt encrypts the plaintext using the provided key with AES-256-CBC
func AesEncrypt(data string) (string, error) {
	key := "QP#es2#L0*4715nR&e3tWFT%156dT&3f"
	salt := "fvB23^#$IHK8#dpt"
	keyBytes := []byte(key)
	plainTextBytes := []byte(data + salt)
	if len(keyBytes) != 32 {
		return "", errors.New("key length must be 32 bytes")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	plaintext := pkcs7Padding(plainTextBytes, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, []byte(salt))
	mode.CryptBlocks(ciphertext, plaintext)

	return hex.EncodeToString(ciphertext) + Md5Encrypt(salt), nil
}

// AesDecrypt decrypts the ciphertext using the provided key with AES-256-CBC
func AesDecrypt(cipherText string) (string, error) {
	key := "QP#es2#L0*4715nR&e3tWFT%156dT&3f"
	salt := "fvB23^#$IHK8#dpt"
	// 使用md5加密salt， 并从cipherText中去掉salt的md5值
	saltMd5 := Md5Encrypt(salt)
	if len(cipherText) < 32 {
		return "", errors.New("ciphertext too short")
	}
	cipherTextLen := len(cipherText)
	cipherTextSalt := cipherText[cipherTextLen-len(saltMd5):]
	cipherTextData := cipherText[:cipherTextLen-len(cipherTextSalt)]
	if cipherTextSalt != saltMd5 {
		return "", errors.New("salt is not correct")
	}
	cipherText = cipherTextData
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	ciphertextBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertextBytes))
	mode := cipher.NewCBCDecrypter(block, []byte(salt))
	mode.CryptBlocks(plaintext, ciphertextBytes)

	plaintext = pkcs7UnPadding(plaintext)
	return string(plaintext), nil
}

func Md5Encrypt(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func CompareHashAndPassword(hashedPassword, password string) bool {
	// 解析出hash中的password
	decryptPassword, err := AesDecrypt(hashedPassword)
	if err != nil {
		return false
	}
	return decryptPassword[:len(password)] == password
}
