package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func unpad(buf []byte) []byte {
	padding := int(buf[len(buf)-1])
	return buf[:len(buf)-padding]
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
	var (
		plaintext []byte
		iv        []byte
		block     cipher.Block
		mode      cipher.BlockMode
		err       error
	)
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("Invalid ciphertext length: to short")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("Invalid ciphertext length: not a multiple of blocksize")

	}
	iv = ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	mode = cipher.NewCBCDecrypter(block, iv)
	plaintext = make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = unpad(plaintext)

	return plaintext, nil
}
