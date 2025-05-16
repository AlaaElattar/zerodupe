package config

// Config holds server configuration
type Config struct {
	Port       int
	StorageDir string
}

func NewConfig() Config {
	return Config{
		Port:       8080,
		StorageDir: "server/storage",
	}
}
