package worker

import (
	"encoding/json"

	"github.com/santoshkpatro/unbit/internal/models"
)

func StackTraceToJSON(stacktrace []models.StackFrame) string {
	jsonData, err := json.Marshal(stacktrace)
	if err != nil {
		return "[]"
	}
	return string(jsonData)
}
