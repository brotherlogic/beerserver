package main

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"

	pb "github.com/brotherlogic/beerserver/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

//AddBeer adds a beer to the server
func (s *Server) AddBeer(ctx context.Context, req *pb.AddBeerRequest) (*pb.AddBeerResponse, error) {
	t := time.Now()
	b := s.ut.GetBeerDetails(req.Beer.Id)
	b.Size = req.Beer.Size

	minTime := int64(0)
	before := false
	for _, b := range s.config.Drunk {
		if b.Id == req.Beer.Id && b.DrinkDate > minTime {
			minTime = b.DrinkDate
			before = true
		}
	}
	if minTime == 0 {
		minTime = time.Now().Unix()
	}

	minDate := time.Unix(minTime, 0)
	if before {
		minDate = minDate.AddDate(1, 0, 0)
	}

	for i := 0; i < int(req.Quantity); i++ {
		newBeer := proto.Clone(req.Beer).(*pb.Beer)
		newBeer.DrinkDate = minDate.Unix()
		newBeer.Name = b.Name
		s.addBeerToCellar(newBeer)
		minDate = minDate.AddDate(1, 0, 0)
	}

	s.save()
	s.LogFunction("AddBeer", t)
	return &pb.AddBeerResponse{}, nil
}

//ListBeers gets the beers in the cellar
func (s *Server) ListBeers(ctx context.Context, req *pb.ListBeerRequest) (*pb.ListBeerResponse, error) {
	s.LogTrace(ctx, "ListBeers", time.Now(), pbt.Milestone_START_FUNCTION)
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

	s.LogTrace(ctx, "ListBeers", time.Now(), pbt.Milestone_END_FUNCTION)
	return &pb.ListBeerResponse{Beers: beers}, nil
}
