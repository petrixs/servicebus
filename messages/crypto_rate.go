package messages

type CryptoCurrencyRate struct {
	Currency string  `json:currency`
	Rate     float64 `json:rate`
}

func (c CryptoCurrencyRate) GetRoutingKey() string {

	return "crypto.rate"
}
