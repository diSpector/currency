package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/diSpector/currency.git/internal/currency/server"
	"github.com/diSpector/currency.git/internal/currency/server/cache"
	"github.com/diSpector/currency.git/internal/currency/server/usecase"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	pb "github.com/diSpector/currency.git/pkg/currency/grpc"
)

const API_URL = `https://freetestapi.com/api/v1/currencies`
const PORT = `50060`

func main() {
	log.Println(`server run`)

	innerCache := cache.NewInnerCache()

	redisClient := redis.NewClient(&redis.Options{
		Addr: `localhost:6379`,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln(`err ping Redis:`, err)
	}
	defer redisClient.Close()

	redisCache := cache.NewRedisCache(redisClient, 2 * time.Minute)
	
	multiLayerCache := cache.NewMultiLayerCache(innerCache, redisCache)

	curUseCase := usecase.New(API_URL, multiLayerCache)

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
