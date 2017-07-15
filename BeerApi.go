package main

import (
	"errors"
	"log"
	"strconv"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
)

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
