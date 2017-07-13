package main

import (
	"context"
	"net/http"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	cellar *pb.BeerCellar
}

type mainFetcher struct{}

func (fetcher mainFetcher) Fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

// DoRegister Registers this server
func (s *Server) DoRegister(server *grpc.Server) {
	//Nothing to register
}

// ReportHealth Determines if the server is healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes this server
func (s *Server) Mote(master bool) error {
	return nil
}

//AddBeer adds a beer to the cellar
func (s *Server) AddBeer(ctx context.Context, beer *pb.Beer) (*pb.Cellar, error) {
	cel := AddBuilt(s.cellar, beer)
	return cel, nil
}

//GetBeer gets a beer from the cellar
func (s *Server) GetBeer(ctx context.Context, beer *pb.Beer) (*pb.Beer, error) {
	var bestBeer *pb.Beer

	for _, cellar := range s.cellar.GetCellars() {
		b1 := cellar.GetBeers()[0]
		if b1 != nil && b1.Size == beer.Size {
			if bestBeer == nil {
				bestBeer = b1
			} else {
				if bestBeer.DrinkDate > b1.DrinkDate {
					bestBeer = b1
				}
			}
		}
	}

	return bestBeer, nil
}

//Init builds a server
func Init() Server {
	s := Server{&goserver.GoServer{}, &pb.BeerCellar{}}
	s.Register = &s
	return s
}

func main() {
	server := Init()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)
	server.PrepServer()
	server.RegisterServer("beerserver", false)
	server.Serve()
}
