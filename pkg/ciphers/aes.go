package ciphers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

func GenerateKeyAES() []byte {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		log.Fatal("Could not generate AES key")
	}

	return buf
}

func EncryptDataAES(data []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func DecryptDataAES(encrypted []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
        return nil, err
    }

	nonce, encrypted := encrypted[:nonceSize], encrypted[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
    if err != nil {
        return nil, err
    }

	return plaintext, nil
}
