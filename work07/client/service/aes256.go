package service

import "github.com/mervick/aes-everywhere/go/aes256"

var PwdKey = []byte("LIANGDONG")

func EncryptPasswords(pwd []byte) ([]byte, error) { //加密密码
	encrypted := aes256.Encrypt(string(pwd), string(PwdKey))

	return []byte(encrypted), nil
}

func DecryptThePassword(pwd []byte) ([]byte, error) { //解密密码
	decrypted := aes256.Decrypt(string(pwd), string(PwdKey))
	return []byte(decrypted), nil
}
