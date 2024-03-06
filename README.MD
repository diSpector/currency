### Generate structs for grpc
```sh
protoc api/proto/currency.proto --go_out=./ --go-grpc_out=./
```