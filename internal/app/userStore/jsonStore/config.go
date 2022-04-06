package jsonStore

// Config ...
type Config struct {
	AuthFile   string `toml:"auth_file"`
	NewUserKey string `toml:"new_user_key"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		AuthFile:   ":users.json",
		NewUserKey: "[new_user_key]",
	}
}
