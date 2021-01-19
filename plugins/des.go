package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"strings"
)

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go

func padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

func Encrypt(in string) (string, error) {
	s := strings.Split(in, ",")
	if len(s) != 2 {
		return "", errors.New("error param")
	}
	println(s[1])
	block, _ := des.NewCipher([]byte(s[1]))
	src := padding([]byte(s[0]), block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, []byte(s[1]))
	blockmode.CryptBlocks(src, src)
	//b, _ := hex.DecodeString(src)
	encodedStr := hex.EncodeToString(src)
	return encodedStr, nil
}
func Decrypt(in string) (string, error) {
	s := strings.Split(in, ",")
	if len(s) != 2 {
		return "", errors.New("error param")
	}
	src, _ := hex.DecodeString(s[0])
	key := []byte(s[1])
	block, _ := des.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return string(src), nil
}
func main() {
	r, e1 := Encrypt("123456,xuegx123")
	println(r)
	println(e1)
	s, e := Decrypt(r + ",xuegx123")
	println(s)
	println(e)

	//s, e = Decrypt("40F20D4E761C9C227049FFADDCC9ED96E2C309B905A242C1087B567EB94FFA9054235316438A47B0CCB1F27CA120F9F414A50A3D5D5CA350,11111111")
	//println(s)
	//println(e)

}
