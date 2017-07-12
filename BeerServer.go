package main

import (
	"context"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beer/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	cellar *pb.BeerCellar
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
	server.RegisterServer("beer", false)
	server.Serve()
}
