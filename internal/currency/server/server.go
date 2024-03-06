package server

import (
	"time"

	"github.com/diSpector/currency.git/internal/currency/server/usecase"
	pb "github.com/diSpector/currency.git/pkg/currency/grpc"
	"github.com/pkg/errors"
)

type Server struct {
	use usecase.CurrencyUseCase
	pb.UnimplementedCurrencyApiServer
}

func NewServer(use usecase.CurrencyUseCase) *Server {
	return &Server{
		use: use,
	}
}

func (s *Server) GetCurrency(req *pb.CurrencyRequest, stream pb.CurrencyApi_GetCurrencyServer) error {
	cursCh, errCh := s.use.GetCurrenciesByCodes(req.Name)

	for {
		select {
		case cur, ok := <-cursCh:
			if !ok {
				return nil
			} else {
				stream.Send(&pb.CurrencyResponse{
					Name: cur.Name,
					Code: cur.Code,
					Rate: cur.Rate,
				})
			}
		case err := <-errCh:
			if err != nil {
				return err
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *Server) GetCurrencyWoChans(req *pb.CurrencyRequest, stream pb.CurrencyApi_GetCurrencyServer) error {
	if req == nil || len(req.Name) == 0 {
		return errors.New(`currencies slice is empty`)
	}

	curs, err := s.use.GetCurrenciesByCodesWoChans(req.Name)
	if err != nil {
		return err
	}

	for i := range curs {
		stream.Send(&pb.CurrencyResponse{
			Name: curs[i].Name,
			Code: curs[i].Code,
			Rate: curs[i].Rate,
		})
	}

	return nil
}
