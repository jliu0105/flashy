package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//（Adevanced Encryption Standard ,AES）

// 16 byte: AES-128
// 24 byte: AES-192
// 32 byte: AES-256

// cannot disclose the key
var PwdKey = []byte("DIS**#KKKDJJSKDI")

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// copy padding number of []byte{byte(padding)}, combine into new slice and return that
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// reverse-filling, delete the added string
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("Encrpyted string Error！")
	} else {
		// get the length of padding
		unpadding := int(origData[length-1])
		// cut the clice, delete the added string and return the surface string
		return origData[:(length - unpadding)], nil
	}
}

// Ecrypt
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	// Create an encryption algorithm instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// get block size
	blockSize := block.BlockSize()
	// Fill the data to make the data length meet the demand
	origData = PKCS7Padding(origData, blockSize)
	//Adopt CBC encryption mode in AES encryption method
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//Execute ecrypt
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//achieve decryption
func AesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	// Create an encryption algorithm instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	// decryption
	blockMode.CryptBlocks(origData, cypted)
	// remove padding
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

//Encrypt base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEcrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// Decrypt
func DePwdCode(pwd string) ([]byte, error) {
	// Decrypt base64 string
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// Execute AES decrypt
	return AesDeCrypt(pwdByte, PwdKey)

}
