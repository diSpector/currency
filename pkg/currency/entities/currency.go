package entities

type Currency struct {
	Name string  `json:"name"`
	Code string  `json:"code"`
	Rate float64 `json:"exchange_rate"`
}
