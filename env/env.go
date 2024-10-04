package env

import (
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

const TokenExpire = 604800

func init() {
	var err error
	if err = env.Parse(&AliPayConf); err != nil {
		logrus.Println("parse alipay config error")
	}

	if err = env.Parse(&JwtConf); err != nil {
		logrus.Println("parse jwt config error")
	}

	if err = env.Parse(&PaypalConf); err != nil {
		logrus.Println("parse paypal config error")
	}
	if err = env.Parse(&AppConf); err != nil {
		logrus.Println("parse app config error")
	}
}
