package model

type NamedLogger struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type NamedLoggerList struct {
	Loggers []*NamedLogger `json:"loggers"`
}
