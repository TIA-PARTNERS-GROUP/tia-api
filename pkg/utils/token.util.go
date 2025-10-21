package utils
import (
	"crypto/sha512"
	"encoding/hex"
)
func HashToken(token string) string {
	hasher := sha512.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}
