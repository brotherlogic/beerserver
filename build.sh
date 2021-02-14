#protoc -I=./proto --go_out=plugins=grpc:./proto proto/beer.proto
#mv proto/github.com/brotherlogic/beerserver/* ./proto

/home/simon/.local/bin/protoc --proto_path ../../../ -I=./proto --go_out ./proto --go_opt=paths=source_relative --go-grpc_out ./proto --go-grpc_opt paths=source_relative --go-grpc_opt=require_unimplemented_servers=false proto/beer.proto

mv proto/github.com/brotherlogic/beerserver/proto/* ./proto