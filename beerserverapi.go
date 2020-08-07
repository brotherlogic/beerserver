package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
	rpb "github.com/brotherlogic/reminders/proto"
)

//AddBeer adds a beer to the server
func (s *Server) AddBeer(ctx context.Context, req *pb.AddBeerRequest) (*pb.AddBeerResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	b := s.ut.GetBeerDetails(req.Beer.Id)
	b.Size = req.Beer.Size

	minTime := int64(0)
	before := false
	for _, b := range config.Drunk {
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
		newBeer.BreweryId = b.BreweryId
		s.addBeerToCellar(newBeer, config)
		minDate = minDate.AddDate(1, 0, 0)
	}

	// Reorder all the cellars
	for i, cellar := range config.GetCellar().GetSlots() {
		s.reorderCellar(ctx, cellar, int32(i))
	}

	return &pb.AddBeerResponse{}, s.save(ctx, config)
}

//DeleteBeer deletes a beer
func (s *Server) DeleteBeer(ctx context.Context, req *pb.DeleteBeerRequest) (*pb.DeleteBeerResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	if (req.Uid) == 0 {
		return nil, fmt.Errorf("Cannot specify zero id for delete request")
	}

	for _, c := range config.GetCellar().GetSlots() {
		for i, b := range c.Beers {
			if b.Uid == req.Uid {
				s.Log(fmt.Sprintf("Deleteing beer %v", b))
				c.Beers = append(c.Beers[:i], c.Beers[i+1:]...)
				return &pb.DeleteBeerResponse{}, s.save(ctx, config)
			}
		}
	}

	for i, b := range config.GetCellar().GetOnDeck() {
		if b.Uid == req.Uid {
			s.Log(fmt.Sprintf("Deleteing beer %v", b))
			config.Cellar.OnDeck = append(config.Cellar.OnDeck[:i], config.Cellar.OnDeck[i+1:]...)
			return &pb.DeleteBeerResponse{}, s.save(ctx, config)
		}
	}

	return nil, fmt.Errorf("Unable to locate beer with uid %v", req.Uid)
}

//ListBeers gets the beers in the cellar
func (s *Server) ListBeers(ctx context.Context, req *pb.ListBeerRequest) (*pb.ListBeerResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	beers := make([]*pb.Beer, 0)

	if req.OnDeck {
		for _, b := range config.GetCellar().GetOnDeck() {
			beers = append(beers, b)
		}
	} else {

		for _, c := range config.GetCellar().GetSlots() {
			for _, b := range c.Beers {
				beers = append(beers, b)
			}
		}
	}

	return &pb.ListBeerResponse{Beers: beers}, nil
}

func (s *Server) singleConsolidate(ctx context.Context, config *pb.Config) error {
	for i, co := range config.Cellar.Slots {
		if int(co.NumSlots) > len(co.Beers) {
			for j, cs := range config.Cellar.Slots {
				if j > i && cs.Accepts == co.Accepts && len(cs.Beers) > 0 {
					beer := cs.Beers[0]
					cs.Beers = cs.Beers[1:]
					co.Beers = append(co.Beers, beer)
					s.reorderCellar(ctx, co, int32(i))

					return nil
				}
			}
		}
	}

	return fmt.Errorf("No consolidation needed")
}

// Consolidate the cellars into a smaller number
func (s *Server) Consolidate(ctx context.Context, req *pb.ConsolidateRequest) (*pb.ConsolidateResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	//Move beers around until we can move no more
	for s.singleConsolidate(ctx, config) == nil {
		// pass
	}

	//Clean out the existing cellars
	newCellarSlots := make([]*pb.CellarSlot, 0)
	for _, c := range config.Cellar.Slots {
		if len(c.Beers) > 0 {
			newCellarSlots = append(newCellarSlots, c)
		}
	}

	config.Cellar.Slots = newCellarSlots

	return &pb.ConsolidateResponse{}, s.save(ctx, config)
}

// Update runs a full update
func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	err = s.syncDrunk(ctx, config)
	if err != nil {
		return nil, err
	}

	s.clearDeck(config)
	err = s.moveToOnDeck(ctx, config, time.Now())
	if err != nil {
		return nil, err
	}

	if time.Now().Hour() < 14 && (time.Now().Weekday() == time.Monday || time.Now().Weekday() == time.Wednesday || time.Now().Weekday() == time.Friday) {
		err = s.refreshStash(ctx, config)
		if err != nil {
			return nil, err
		}
	}

	return &pb.UpdateResponse{}, s.save(ctx, config)
}

// Receive a request ot update
func (s *Server) Receive(ctx context.Context, req *rpb.ReceiveRequest) (*rpb.ReceiveResponse, error) {
	_, err := s.Update(ctx, &pb.UpdateRequest{})
	return &rpb.ReceiveResponse{}, err
}
