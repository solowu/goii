package goii

import (
	"crypto/rc4"
	"encoding/base64"
)

func Rc4En(s string, k string) string {
	key := []byte(k)
	src := []byte(s)
	dst := make([]byte, len(src))
	if c, err := rc4.NewCipher(key); err != nil {
		return ""
	} else {
		c.XORKeyStream(dst, src)
		return base64.StdEncoding.EncodeToString(dst)
	}
}
func Rc4De(s string, k string) string {
	key := []byte(k)

	src, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	dst := make([]byte, len(src))
	if c, err := rc4.NewCipher(key); err != nil {
		return ""
	} else {
		c.XORKeyStream(dst, src)
		return string(dst)
	}
}
