package env

var JwtConf JwtConfig

type JwtConfig struct {
	Secret         string `env:"JWT_SECRET"`
	AdminJwtSecret string `env:"JWT_ADMIN_JWT"`
}
