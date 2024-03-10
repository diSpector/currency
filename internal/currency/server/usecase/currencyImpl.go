package usecase

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/diSpector/currency.git/internal/currency/helpers"
	"github.com/diSpector/currency.git/internal/currency/server/cache"
	"github.com/diSpector/currency.git/pkg/currency/entities"
	"github.com/pkg/errors"
)

type CurrencyUseCaseImpl struct {
	url   string
	cache cache.Cache
}

func New(url string, cache cache.Cache) *CurrencyUseCaseImpl {
	return &CurrencyUseCaseImpl{
		url:   url,
		cache: cache,
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

	var notFoundInCache []string
	uniqCodes := helpers.GetUnique(codes)

	for i := range uniqCodes {
		v, err := s.cache.Get(uniqCodes[i])
		if err != nil {
			if errors.Is(err, cache.ErrNotFound) {
				notFoundInCache = append(notFoundInCache, uniqCodes[i])
			} else if !errors.Is(err, cache.ErrFoundInAbsent) {
				chErr <- err
				return
			}
		} else {
			if v != nil {
				ch <- v
			}
		}
	}

	if len(notFoundInCache) > 0 {
		log.Println(`request api for currencies`)
		cursFromApi, err := s.getAllCurrenciesFromApi(codes)
		if err != nil {
			chErr <- errors.Wrap(err, `err from api`)
			return
		}

		var curMap = make(map[string]struct{})

		for i := range notFoundInCache {
			var isFound = false
			for j := range cursFromApi {
				if notFoundInCache[i] == cursFromApi[j].Code {
					if _, ok := curMap[notFoundInCache[i]]; !ok {
						curMap[notFoundInCache[i]] = struct{}{}
						err := s.cache.Set(notFoundInCache[i], cursFromApi[j], 1*time.Minute)
						if err != nil {
							chErr <- err
							return
						}
						isFound = true
						ch <- cursFromApi[j]
					}
				}
			}
			// codes not found inside api response
			if !isFound {
				err := s.cache.Set(notFoundInCache[i], nil, 10*time.Minute)
				if err != nil {
					chErr <- err
					return
				}
			}
		}
	}

}

func (s *CurrencyUseCaseImpl) getAllCurrenciesFromApi(codes []string) ([]*entities.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
