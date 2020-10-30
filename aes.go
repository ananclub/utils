package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

type AesCrypt struct {
	Key       []byte
	Iv        []byte
	Blocker   cipher.Block
	Encrypter cipher.BlockMode
	Decrypter cipher.BlockMode
}

func NewAes(key, iv []byte) (aesCrypt *AesCrypt, err error) {

	Blocker, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	if len(iv) != Blocker.BlockSize() {
		err = errors.New("cipher: IV length must equal block size " + strconv.Itoa(Blocker.BlockSize()))
		return
	}
	Encrypter := cipher.NewCBCEncrypter(Blocker, iv)
	Decrypter := cipher.NewCBCDecrypter(Blocker, iv)

	aesCrypt = &AesCrypt{Key: key, Iv: iv, Blocker: Blocker, Encrypter: Encrypter, Decrypter: Decrypter}
	return
}

func (a *AesCrypt) EncryptBlock(data []byte) ([]byte, error) {

	cipherBytes := make([]byte, len(data))
	a.Encrypter.CryptBlocks(cipherBytes, data)
	return cipherBytes, nil
}

func (a *AesCrypt) DecryptBlock(src []byte) (data []byte, err error) {
	decrypted := make([]byte, len(src))

	a.Decrypter.CryptBlocks(decrypted, src)
	return decrypted, nil
}
func (a *AesCrypt) Encrypt(data []byte) ([]byte, error) {

	content := a.PKCS5Padding(data)
	cipherBytes := make([]byte, len(content))
	a.Encrypter.CryptBlocks(cipherBytes, content)
	return cipherBytes, nil
}

func (a *AesCrypt) Decrypt(src []byte) (data []byte, err error) {
	decrypted := make([]byte, len(src))

	a.Decrypter.CryptBlocks(decrypted, src)
	return a.PKCS5UnPadding(decrypted), nil
}

func (a *AesCrypt) PKCS5Padding(cipherText []byte) []byte {
	blockSize := a.Blocker.BlockSize()
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
func (a *AesCrypt) PKCS5PaddingWithBlockSize(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
func (a *AesCrypt) PKCS5UnPadding(ciphertext []byte) []byte {
	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}
