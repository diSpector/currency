package main

import (
	"log"
	"net"

	"github.com/diSpector/currency.git/internal/currency/server"
	"github.com/diSpector/currency.git/internal/currency/server/usecase"
	"google.golang.org/grpc"

	pb "github.com/diSpector/currency.git/pkg/currency/grpc"
)

const API_URL = `https://freetestapi.com/api/v1/currencies`
const PORT = `50060`

func main() {
	log.Println(`server run`)

	curUseCase := usecase.New(API_URL)

	serv := server.NewServer(curUseCase)

	lis, err := net.Listen(`tcp`, `localhost:`+PORT)
	if err != nil {
		log.Fatalln(`err listen:`, err)
	}

	s := grpc.NewServer()
	pb.RegisterCurrencyApiServer(s, serv)

	log.Println(`grpc server is listening`)

	if err := s.Serve(lis); err != nil {
		log.Fatalln(`failed serve:`, err)
	}
}
