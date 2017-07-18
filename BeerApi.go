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

	return beers[rand.Intn(len(beers))], nil
}

//RemoveBeer removes a beer from the cellar
func (s *Server) RemoveBeer(ctx context.Context, in *pb.Beer) (*pb.Beer, error) {
	beer := RemoveBeer(s.cellar, in.Id)
	s.saveCellar()
	return beer, nil
}

//GetName gets the name of the beer
func (s *Server) GetName(ctx context.Context, in *pb.Beer) (*pb.Beer, error) {
	if val, ok := s.nameCache[in.Id]; ok {
		in.Name = val
		return in, nil
	}

	name := s.ut.GetBeerName(in.Id)
	s.nameCache[in.Id] = name

	in.Name = name
	return in, nil
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
			return cellar, nil
		}
	}

	if needWrite {
		s.saveCellar()
	}

	return nil, errors.New("Unable to locate cellar")
}
