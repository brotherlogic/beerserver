package main

import (
	"fmt"
	"io/ioutil"

	pb "github.com/brotherlogic/beerserver/proto"
)

func findIndex(slot *pb.CellarSlot, date int64) int32 {
	return 0
}

func (s *Server) addBeerToCellar(b *pb.Beer) error {
	//Do we have room for this beer
	for _, c := range s.config.Cellar.Slots {
		if c.Accepts == b.Size {
			if int(c.NumSlots) > len(c.Beers) {
				index := findIndex(c, b.DrinkDate)
				for _, beer := range c.Beers {
					if beer.Index >= index {
						beer.Index++
					}
				}
				b.Index = index
				c.Beers = append(c.Beers, b)

				return nil
			}
		}
	}

	return fmt.Errorf("Unable to find space for this beer")
}

func (s *Server) loadDrunk(filestr string) {
	data, _ := ioutil.ReadFile(filestr)
	b := s.ut.convertDrinkListToBeers(string(data), mainUnmarshaller{})
	s.config.Drunk = b
}

func (s *Server) syncDrunk(f httpResponseFetcher) {
	lastID := s.config.Drunk[len(s.config.Drunk)-1].CheckinId
	ndrinks := s.ut.getLastBeers(f, mainConverter{}, mainUnmarshaller{}, lastID)
	s.config.Drunk = append(s.config.Drunk, ndrinks...)
}
