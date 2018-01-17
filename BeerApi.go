package main

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
)

//GetToDrink gets beers in the to drink list
func (s *Server) GetToDrink(ctx context.Context, req *pb.GetToDrinkRequest) (*pb.GetToDrinkResponse, error) {
	return &pb.GetToDrinkResponse{Beers: s.cellar.GetTodrink().GetBeers()}, nil
}

//AddBeer adds a beer to the cellar
func (s *Server) AddBeer(ctx context.Context, beer *pb.Beer) (*pb.Cellar, error) {
	if beer.Name == "" || beer.GetAbv() == 0 {
		beer.Name, beer.Abv = s.ut.GetBeerDetails(beer.Id)
	}
	var cel *pb.Cellar
	if beer.Staged {
		s.cellar.GetTodrink().Beers = append(s.cellar.GetTodrink().GetBeers(), beer)
		cel = &pb.Cellar{}
	} else {
		cel = AddBuilt(s.cellar, beer)
	}
	s.saveCellar()
	return cel, nil
}

//GetDrunk gets a list of beers drunk
func (s *Server) GetDrunk(ctx context.Context, in *pb.Empty) (*pb.BeerList, error) {
	t := time.Now()
	list := &pb.BeerList{Beers: s.cellar.Drunk}

	for _, beer := range list.GetBeers() {
		s.recacheBeer(beer)
	}
	s.saveCellar()

	s.LogFunction("GetDrunk", t)
	return list, nil
}

func authGet(ld *pb.Beer, r *pb.Beer) bool {
	t := time.Now().Add(time.Hour * -24)
	return ld.GetDrinkDate()-t.Unix() < 0 || r.Size == "bomber"
}

//GetBeer gets a beer from the cellar
func (s *Server) GetBeer(ctx context.Context, beer *pb.Beer) (*pb.Beer, error) {

	if len(s.cellar.Drunk) > 0 && !authGet(s.cellar.Drunk[len(s.cellar.Drunk)-1], beer) {
		return nil, errors.New("Unauthorized get for " + beer.Size)
	}

	var beers []*pb.Beer
	for _, cellar := range s.cellar.GetCellars() {
		for _, b := range cellar.GetBeers() {
			if b.Size == beer.Size && b.Staged {
				return b, nil
			}

			if b.DrinkDate < time.Now().Unix() && b.Size == beer.Size {
				beers = append(beers, b)
			}
		}
	}

	rBeer := beers[rand.Intn(len(beers))]
	s.recacheBeer(rBeer)
	rBeer.Staged = true
	s.saveCellar()
	return rBeer, nil
}

//RemoveBeer removes a beer from the cellar
func (s *Server) RemoveBeer(ctx context.Context, in *pb.Beer) (*pb.Beer, error) {
	beer := RemoveBeer(s.cellar, in.Id)
	s.saveCellar()
	return beer, nil
}

//GetServerState gets the state of the system
func (s *Server) GetServerState(ctx context.Context, in *pb.Empty) (*pb.ServerState, error) {
	t := time.Now()
	r := &pb.ServerState{}
	fb, fs := GetTotalFreeSlots(s.cellar)
	r.FreeSmall = int32(fs)
	r.FreeBombers = int32(fb)

	s.LogFunction("GetState", t)
	return r, nil
}

//GetCellar gets a single cellar
func (s *Server) GetCellar(ctx context.Context, in *pb.Cellar) (*pb.Cellar, error) {
	needWrite := false
	for i, cellar := range s.cellar.Cellars {

		if len(cellar.Name) == 0 {
			cellar.Name = "cellar" + strconv.Itoa(i+1)
			needWrite = true
		}

		if cellar.Name == in.Name {
			if needWrite {
				s.saveCellar()
				needWrite = false
			}

			//Ensure the beers are named
			for _, beer := range cellar.Beers {
				s.recacheBeer(beer)
			}

			return cellar, nil
		}
	}

	if needWrite {
		s.saveCellar()
	}

	return nil, errors.New("Unable to locate cellar")
}
