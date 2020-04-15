protoc -I=./proto --go_out=plugins=grpc:./proto proto/beer.proto
mv proto/github.com/brotherlogic/beerserver/* ./proto
