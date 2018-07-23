package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
	pbt "github.com/brotherlogic/tracer/proto"
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
		b.InCellar = int32(bestCellar)
		c.Beers = append(c.Beers, b)

		return nil

	}

	return fmt.Errorf("Unable to find space for this beer")
}

func (s *Server) loadDrunk(filestr string) {
	data, _ := ioutil.ReadFile(filestr)
	b := s.ut.convertDrinkListToBeers(string(data), mainUnmarshaller{})
	for _, beer := range b {
		found := false
		for _, d := range s.config.Drunk {
			if d.CheckinId == beer.CheckinId {
				found = true
			}
		}

		if !found {
			s.config.Drunk = append(s.config.Drunk, beer)
		}
	}
}

func (s *Server) syncDrunk(f httpResponseFetcher) {
	lastID := int32(0)
	for _, drunk := range s.config.Drunk {
		if drunk.CheckinId > lastID {
			lastID = drunk.CheckinId
		}
	}
	ndrinks, err := s.ut.getLastBeers(f, mainConverter{}, mainUnmarshaller{}, lastID)
	if err == nil {
		s.config.Drunk = append(s.config.Drunk, ndrinks...)
		s.config.LastSync = time.Now().Unix()
		s.save()
	} else {
		s.Log(fmt.Sprintf("Sync error with %v: %v", lastID, err))
	}
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

	s.save()
}

func (s *Server) clearDeck(ctx context.Context) {
	s.LogTrace(ctx, "ClearDeck", time.Now(), pbt.Milestone_START_FUNCTION)

	s.lastClean = time.Now()

	for _, bdr := range s.config.Drunk {
		for i, bde := range s.config.Cellar.OnDeck {
			if bdr.Id == bde.Id && bdr.Id == 2500585 {
				s.Log(fmt.Sprintf("FOUND MATCH %v and %v also %v and %v", bdr.Id, bde.Id, time.Unix(bdr.DrinkDate, 0), time.Unix(bde.DrinkDate, 0)))
			}
			if bdr.Id == bde.Id && bdr.DrinkDate > bde.DrinkDate {
				s.config.Cellar.OnDeck = append(s.config.Cellar.OnDeck[:i], s.config.Cellar.OnDeck[i+1:]...)
			}
		}
	}

	s.LogTrace(ctx, "ClearDeck", time.Now(), pbt.Milestone_END_FUNCTION)
}
