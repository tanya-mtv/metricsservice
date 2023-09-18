package hashsha

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CreateHash(key string, src []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	dst := hex.EncodeToString(h.Sum(nil))
	return dst
}
