### Generate structs for grpc
```sh
protoc api/proto/currency.proto --go_out=./ --go-grpc_out=./
```

### Redis
> https://redis.io/

### Redis + Docker
```sh
docker pull redis
docker run --name new-redis -d  -p 6379:6379 redis
docker exec -it new-redis bash
redis-cli
```