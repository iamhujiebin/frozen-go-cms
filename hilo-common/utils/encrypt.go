package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

func EncryptionData(data interface{}, aeskey []byte) string {
	str, _ := json.Marshal(data)
	pass := str
	xpass, err := aesEncrypt(pass, aeskey)
	if err != nil {
		return ""
	}
	pass64 := base64.StdEncoding.EncodeToString(xpass)
	return pass64
}

func Encryption(str string, aeskey []byte) string {
	pass := []byte(str)
	xpass, err := aesEncrypt(pass, aeskey)
	if err != nil {
		return ""
	}
	pass64 := base64.StdEncoding.EncodeToString(xpass)
	return pass64
}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Decryption(str string, aeskey []byte) (string, error) {
	pass64, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	xpass, err := aesDecrypt(pass64, aeskey)
	if err != nil {
		return "", err
	}
	return string(xpass), nil
}

func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
