package config

type Config struct {
	// Server configuration
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`

	// Logging configuration
	Logging struct {
		Level string `json:"level"`
	} `json:"logging"`

	// Database configuration
	Database struct {
		Type     string `json:"type"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`

	// Other configurations can be added here
}
