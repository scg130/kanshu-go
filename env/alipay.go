package env

var AliPayConf AliPayConfig

type AliPayConfig struct {
	AppId         string `env:"ALIPAY_APPID"`
	RsaPrivateKey string `env:"ALIPAY_RSA_PRIVATE_KEY"`
	RsaPublicKey  string `env:"ALIPAY_RSA_PUBLIC_KEY"`
	NotifyUrl     string `env:"ALIPAY_NOTIFY_URL"`
}
