package config

// Config holds server configuration
type Config struct {
	Port               int    `json:"port"`
	StorageDir         string `json:"storage_dir"`
	JWTSecret          string `json:"jwt_secret"`
	AccessTokenExpiry  int    `json:"access_token_expiry"`  // in minutes
	RefreshTokenExpiry int    `json:"refresh_token_expiry"` // in hours
}

func NewConfig(port int, storageDir string, jwtSecret string, accessTokenExpiry int, refreshTokenExpiry int) Config {
	return Config{
		Port:               port,
		StorageDir:         storageDir,
		JWTSecret:          jwtSecret,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}
}
