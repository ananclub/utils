package utils

import "encoding/base64"

func Base64StdEncode(b []byte) (s string) {
	return base64.StdEncoding.EncodeToString(b)
}
func Base64StdDecode(s string) (b []byte, err error) {
	return base64.StdEncoding.DecodeString(s)
}
func Base64URLEncode(b []byte) (s string) {
	return base64.URLEncoding.EncodeToString(b)
}
func Base64URLDecode(s string) (b []byte, err error) {
	return base64.URLEncoding.DecodeString(s)
}
