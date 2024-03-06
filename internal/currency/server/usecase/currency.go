package usecase

import "github.com/diSpector/currency.git/pkg/currency/entities"

type CurrencyUseCase interface {
	GetCurrenciesByCodes(codes []string) (chan *entities.Currency, chan error)
	GetCurrenciesByCodesWoChans(codes []string) ([]*entities.Currency, error)
}
