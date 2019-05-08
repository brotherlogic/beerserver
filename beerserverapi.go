package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
)

//AddBeer adds a beer to the server
func (s *Server) AddBeer(ctx context.Context, req *pb.AddBeerRequest) (*pb.AddBeerResponse, error) {
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
		newBeer.Uid = time.Now().UnixNano()
		s.addBeerToCellar(newBeer)
		minDate = minDate.AddDate(1, 0, 0)
	}

	// Reorder all the cellars
	for _, cellar := range s.config.Cellar.Slots {
		s.reorderCellar(ctx, cellar)
	}

	s.save(ctx)
	return &pb.AddBeerResponse{}, nil
}

//DeleteBeer deletes a beer
func (s *Server) DeleteBeer(ctx context.Context, req *pb.DeleteBeerRequest) (*pb.DeleteBeerResponse, error) {
	if (req.Uid) == 0 {
		return nil, fmt.Errorf("Cannot specify zero id for delete request")
	}

	for _, c := range s.config.Cellar.Slots {
		for i, b := range c.Beers {
			if b.Uid == req.Uid {
				s.Log(fmt.Sprintf("Deleteing beer %v", b))
				c.Beers = append(c.Beers[:i], c.Beers[i+1:]...)
				s.save(ctx)
				return &pb.DeleteBeerResponse{}, nil
			}
		}
	}

	for i, b := range s.config.Cellar.OnDeck {
		if b.Uid == req.Uid {
			s.Log(fmt.Sprintf("Deleteing beer %v", b))
			s.config.Cellar.OnDeck = append(s.config.Cellar.OnDeck[:i], s.config.Cellar.OnDeck[i+1:]...)
			s.save(ctx)
			return &pb.DeleteBeerResponse{}, nil
		}
	}

	return nil, fmt.Errorf("Unable to locate beer with uid %v", req.Uid)
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
