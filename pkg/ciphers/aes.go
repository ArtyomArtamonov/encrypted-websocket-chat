package ciphers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

func GenerateKey() []byte {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		log.Fatal("Could not generate AES key")
	}

	return buf
}

func EncryptData(data []byte, key []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("Could not create new AES cipher")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal("Could not initialize GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        log.Fatal(err)
    }

	return gcm.Seal(nonce, nonce, data, nil)
}

func DecryptData(encrypted []byte, key []byte) ([]byte, error) {
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
