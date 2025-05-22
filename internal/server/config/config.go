package config

// Config holds server configuration
type Config struct {
	Port       int
	StorageDir string
}

// not to be hard coded
func NewConfig(port int, storageDir string) Config {
	return Config{
		Port:       port,
		StorageDir: storageDir,
	}
}
