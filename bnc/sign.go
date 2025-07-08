package bnc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignByHmacSHA256(payload, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(payload))
	return h.Sum(nil)
}

func SignByHmacSHA256ToHex(payload, key string) string {
	return hex.EncodeToString(SignByHmacSHA256(payload, key))
}
