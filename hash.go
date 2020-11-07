package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func Md5String(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
func Md5File(file io.Reader) string {
	h := md5.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}
func Sha1String(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha1File(file io.Reader) string {
	h := sha1.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}
func Sha256String(source string) string {
	h := sha256.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha256File(file io.Reader) string {
	h := sha256.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}
func Sha512String(source string) string {
	h := sha512.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha512File(file io.Reader) string {
	h := sha512.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256(message string, secret string) (str string, err error) {

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	if _, err = h.Write([]byte(message)); err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
func HmacSha1(message string, secret string) (str string, err error) {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	if _, err = h.Write([]byte(message)); err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
