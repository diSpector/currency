package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/diSpector/currency.git/pkg/currency/entities"
	"github.com/pkg/errors"
)

type CurrencyUseCaseImpl struct {
	url string
}

func New(url string) *CurrencyUseCaseImpl {
	return &CurrencyUseCaseImpl{
		url: url,
	}
}

// GetCurrenciesByCodes with channels
func (s *CurrencyUseCaseImpl) GetCurrenciesByCodes(codes []string) (chan *entities.Currency, chan error) {
	ch := make(chan *entities.Currency, 10)
	chErr := make(chan error)

	go s.getCurrencies(codes, ch, chErr)

	return ch, chErr
}

func (s *CurrencyUseCaseImpl) GetCurrenciesByCodesWoChans(codes []string) ([]*entities.Currency, error) {
	cursFromApi, err := s.getAllCurrenciesFromApi(codes)
	if err != nil {
		return nil, err
	}

	var curMap = make(map[string]struct{})
	var res []*entities.Currency

	for i := range codes {
		for j := range cursFromApi {
			if codes[i] == cursFromApi[j].Code {
				if _, ok := curMap[codes[i]]; !ok {
					curMap[codes[i]] = struct{}{}
					res = append(res, cursFromApi[j])
				}
			}
		}
	}

	return res, nil
}

func (s *CurrencyUseCaseImpl) getCurrencies(codes []string, ch chan *entities.Currency, chErr chan error) {
	defer func() {
		close(ch)
		close(chErr)
	}()

	cursFromApi, err := s.getAllCurrenciesFromApi(codes)
	if err != nil {
		chErr <- errors.Wrap(err, `err from api`)
		return
	}

	var curMap = make(map[string]struct{})

	for i := range codes {
		for j := range cursFromApi {
			if codes[i] == cursFromApi[j].Code {
				if _, ok := curMap[codes[i]]; !ok {
					curMap[codes[i]] = struct{}{}
					ch <- cursFromApi[j]
				}
			}
		}
	}
}

func (s *CurrencyUseCaseImpl) getAllCurrenciesFromApi(codes []string) ([]*entities.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var res []*entities.Currency

	dec := json.NewDecoder(resp.Body)
	_, err = dec.Token()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for dec.More() {
		var cur entities.Currency
		err := dec.Decode(&cur)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		res = append(res, &cur)
	}

	_, err = dec.Token()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}
