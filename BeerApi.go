package main

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
)

//AddBeer adds a beer to the cellar
func (s *Server) AddBeer(ctx context.Context, beer *pb.Beer) (*pb.Cellar, error) {
	log.Printf("CELLAR HERE = %v", s.cellar)
	if beer.Name == "" {
		beer.Name = s.ut.GetBeerName(beer.Id)
	}
	cel := AddBuilt(s.cellar, beer)
	s.saveCellar()
	return cel, nil
}

//GetDrunk gets a list of beers drunk
func (s *Server) GetDrunk(ctx context.Context, in *pb.Empty) (*pb.BeerList, error) {
	list := &pb.BeerList{Beers: s.cellar.Drunk}
	return list, nil
}

//GetBeer gets a beer from the cellar
func (s *Server) GetBeer(ctx context.Context, beer *pb.Beer) (*pb.Beer, error) {
	var beers []*pb.Beer
	for _, cellar := range s.cellar.GetCellars() {
		log.Printf("CELLAR: %v", cellar)
		for _, b := range cellar.GetBeers() {
			if b.DrinkDate < time.Now().Unix() && b.Size == beer.Size {
				beers = append(beers, b)
			}
		}
	}

	rBeer := beers[rand.Intn(len(beers))]
	s.recacheBeer(rBeer)
	return rBeer, nil
}

//RemoveBeer removes a beer from the cellar
func (s *Server) RemoveBeer(ctx context.Context, in *pb.Beer) (*pb.Beer, error) {
	beer := RemoveBeer(s.cellar, in.Id)
	s.saveCellar()
	return beer, nil
}

//GetCellar gets a single cellar
func (s *Server) GetCellar(ctx context.Context, in *pb.Cellar) (*pb.Cellar, error) {
	needWrite := false
	for i, cellar := range s.cellar.Cellars {
		log.Printf("CELLAR NAME = %v", cellar.Name)

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
				if beer.Name == "" {
					s.recacheBeer(beer)
				}
			}

			return cellar, nil
		}
	}

	if needWrite {
		s.saveCellar()
	}

	return nil, errors.New("Unable to locate cellar")
}
