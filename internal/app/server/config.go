package server


// Config ...
type Config struct {
	BindAddr string 	`toml:"bind_addr"`
	LogLevel string 	`toml:"log_level"`
	AuthFile string 	`toml:"auth_file"`
	SessionKey string 	`toml:"session_key"`
	NewUserKey string	`toml:"new_user_key"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		AuthFile: "users.json",
	}
}
