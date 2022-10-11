package utils

import (
	"crypto/rand"
	"math/big"
)

const ShortURLCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func GenerateRandomShortStringInSIze(maxSize int) (string, error) {
	ret := make([]byte, maxSize)
    charsetLen := big.NewInt(int64(len(ShortURLCharset)))
	for i := 0; i < maxSize; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		ret[i] = ShortURLCharset[num.Int64()]
	}

	return string(ret), nil
}
