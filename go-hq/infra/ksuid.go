package infra

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"time"
)

func encodeBase62(b []byte) string {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var result bytes.Buffer
	var base = big.NewInt(62)
	var zero = big.NewInt(0)
	n := big.NewInt(0).SetBytes(b)

	for n.Cmp(zero) > 0 {
		mod := big.NewInt(0)
		n.DivMod(n, base, mod) // 计算n/base和n%base
		result.WriteString(string(base62Chars[mod.Int64()]))
	}

	// 由于上面的循环产生的字符串是反向的，所以需要将其反转回来
	encoded := result.Bytes()
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(encoded)
}

func GenerateKsuid() (string, error) {
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(timestampBytes, uint32(timestamp))

	// 随机部分，这里我们生成16字节的随机数据
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// 将时间戳和随机数据拼接
	ksuidBytes := append(timestampBytes, randomBytes...)

	return encodeBase62(ksuidBytes), nil
}
