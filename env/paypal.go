package env

var PaypalConf PaypalConfig

type PaypalConfig struct {
	PaypalClientId  string `env:"PAYPAL_CLIENT_ID"`
	PaypalSecret    string `env:"PAYPAL_SECRET"`
	PaypalReturnUrl string `env:"PAYPAL_RETURN_URL"`
}
