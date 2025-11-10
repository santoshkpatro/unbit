package worker

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/santoshkpatro/unbit/internal/models"
)

func ComputeFingerprint(properties models.Properties) string {
	// Implement fingerprint computation logic here
	data := properties.Message
	for _, f := range properties.Stacktrace {
		data += fmt.Sprintf("%s:%s:%d", f.Function, f.File, f.Line)
	}
	hash := sha1.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
}
