package messages

type CryptoCurrency struct {
	Name string `json:name`
}

type TopCryptoCurrencies struct {
	Currencies []CryptoCurrency `json:currencies`
}

func (t TopCryptoCurrencies) GetRoutingKey() string {
	return "crypto.top"
}
