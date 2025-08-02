package config

// Config holds server configuration
type Config struct {
	Port                   int    `json:"port"`
	StorageDir             string `json:"storage_dir"`
	JWTSecret              string `json:"jwt_secret"`
	AccessTokenExpiryMin   int    `json:"access_token_expiry"`  // in minutes
	RefreshTokenExpiryHour int    `json:"refresh_token_expiry"` // in hours
}

func NewConfig(port int, storageDir string, jwtSecret string, accessTokenExpiryMin int, refreshTokenExpiryHour int) Config {
	return Config{
		Port:                   port,
		StorageDir:             storageDir,
		JWTSecret:              jwtSecret,
		AccessTokenExpiryMin:   accessTokenExpiryMin,
		RefreshTokenExpiryHour: refreshTokenExpiryHour,
	}
}
