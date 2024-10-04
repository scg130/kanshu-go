package env

var AppConf AppConfig

type AppConfig struct {
	Domain   string `env:"DOMAIN"`
	FrontUrl string `env:"FRONT_URL"`
}
