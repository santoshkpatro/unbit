package worker

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/santoshkpatro/unbit/internal/models"
)

func ComputeFingerprint(event models.Event) string {
	// Implement fingerprint computation logic here
	data := event.Message
	for _, f := range event.StackTrace {
		data += fmt.Sprintf("%s:%s:%d", f.Function, f.File, f.Line)
	}
	hash := sha1.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
}
