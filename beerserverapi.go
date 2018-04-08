package main

import (
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
	"github.com/golang/protobuf/proto"
)

//AddBeer adds a beer to the server
func (s *Server) AddBeer(ctx context.Context, req *pb.AddBeerRequest) (*pb.AddBeerResponse, error) {
	b := s.ut.GetBeerDetails(req.Beer.Id)
	b.Size = req.Beer.Size

	minTime := time.Now().Unix()
	before := false
	for _, b := range s.config.Drunk {
		if b.Id == req.Beer.Id && b.DrinkDate < minTime {
			minTime = b.DrinkDate
			before = true
		}
	}

	minDate := time.Unix(minTime, 0)
	if before {
		minDate.AddDate(1, 0, 0)
	}

	for i := 0; i < int(req.Quantity); i++ {
		newBeer := proto.Clone(req.Beer).(*pb.Beer)
		newBeer.DrinkDate = minDate.Unix()
		s.addBeerToCellar(newBeer)
		minDate = minDate.AddDate(1, 0, 0)
	}

	s.save()
	return &pb.AddBeerResponse{}, nil
}

//ListBeers gets the beers in the cellar
func (s *Server) ListBeers(ctx context.Context, req *pb.ListBeerRequest) (*pb.ListBeerResponse, error) {
	beers := make([]*pb.Beer, 0)

	if req.OnDeck {
		for _, b := range s.config.Cellar.OnDeck {
			beers = append(beers, b)
		}
	} else {

		for _, c := range s.config.Cellar.Slots {
			for _, b := range c.Beers {
				beers = append(beers, b)
			}
		}
	}

	return &pb.ListBeerResponse{Beers: beers}, nil
}
