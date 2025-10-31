package ingest

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Event struct {
	Message    string       `json:"message"`
	Level      string       `json:"level"`
	StackTrace []StackFrame `json:"stacktrace"`
}
