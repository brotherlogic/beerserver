package main

import (
	"fmt"
	"io/ioutil"
	"time"

	pb "github.com/brotherlogic/beerserver/proto"
)

func findIndex(slot *pb.CellarSlot, date int64) int32 {
	bestIndex := int32(len(slot.Beers))
	for _, b := range slot.Beers {
		if date < b.DrinkDate && b.Index < bestIndex {
			bestIndex = b.Index
		}
	}
	return bestIndex
}

func (s *Server) addBeerToCellar(b *pb.Beer) error {
	bestIndex := int32(99)
	bestCellar := -1
	//Do we have room for this beer
	for i, c := range s.config.Cellar.Slots {
		if c.Accepts == b.Size {
			if len(c.Beers) == 0 {
				bestCellar = i
			} else {
				if int(c.NumSlots) > len(c.Beers) {
					index := findIndex(c, b.DrinkDate)
					if index < bestIndex {
						bestIndex = index
						bestCellar = i
					}
				}
			}
		}
	}

	if bestCellar >= 0 {
		// We reset when we have to add to a new cellar
		if bestIndex == 99 {
			bestIndex = 0
		}
		c := s.config.Cellar.Slots[bestCellar]
		for _, beer := range c.Beers {
			if beer.Index >= bestIndex {
				beer.Index++
			}
		}
		b.Index = bestIndex
		c.Beers = append(c.Beers, b)

		return nil

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

func (s *Server) moveToOnDeck(t time.Time) {
	for _, cellar := range s.config.Cellar.Slots {
		i := 0
		for i < len(cellar.Beers) {
			if cellar.Beers[i].DrinkDate < t.Unix() {
				cellar.Beers[i].OnDeck = true
				s.config.Cellar.OnDeck = append(s.config.Cellar.OnDeck, cellar.Beers[i])
				cellar.Beers = append(cellar.Beers[:i], cellar.Beers[i+1:]...)
			} else {
				i++
			}
		}
	}
}
