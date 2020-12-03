package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
)

func (s *Server) runExtract(cellar *pb.Cellar, ti time.Time) {
	for _, cellarSlot := range cellar.GetSlots() {
		for s.canExtractBeer(cellarSlot, cellar.GetOnDeck(), ti) {
			beer := s.extractBeer(cellarSlot, ti)

			//There's actually nothing extractable
			if beer == nil {
				break
			}

			//Remove the beer
			nbeers := []*pb.Beer{}
			for _, b := range cellarSlot.GetBeers() {
				if b.GetUid() != beer.GetUid() {
					nbeers = append(nbeers, b)
				}
			}
			cellarSlot.Beers = nbeers

			if cellarSlot.GetMoveTo() == "" {
				cellar.OnDeck = append(cellar.OnDeck, beer)
			} else {
				for _, cs := range cellar.GetSlots() {
					if cs.GetAccepts() == cellarSlot.GetMoveTo() {
						cs.Beers = append(cs.Beers, beer)
					}
				}
			}
		}
	}
}

func (s *Server) extractBeer(cellar *pb.CellarSlot, ti time.Time) *pb.Beer {
	switch cellar.GetExtractionAlgorithm() {
	case pb.Extraction_ON_DATE:
		for _, b := range cellar.GetBeers() {
			if b.GetDrinkDate() < ti.Unix() {
				return b
			}
		}
	case pb.Extraction_STASH_REMOVE, pb.Extraction_HOMEBREW_REMOVE:
		rand.Shuffle(len(cellar.Beers), func(i, j int) { cellar.Beers[i], cellar.Beers[j] = cellar.Beers[j], cellar.Beers[i] })
		if len(cellar.GetBeers()) > 0 {
			return cellar.GetBeers()[0]
		}
	}

	return nil
}

func (s *Server) canExtractBeer(cellar *pb.CellarSlot, deck []*pb.Beer, ti time.Time) bool {
	switch cellar.GetExtractionAlgorithm() {
	case pb.Extraction_ON_DATE:
		return true
	case pb.Extraction_STASH_REMOVE:
		if ti.Weekday() == time.Monday || ti.Weekday() == time.Wednesday || ti.Weekday() == time.Friday {
			if ti.Hour() < 12 {
				for _, b := range deck {
					if b.GetSize() == "stash" {
						return false
					}
				}
				return true
			}
		}
	case pb.Extraction_HOMEBREW_REMOVE:
		if ti.Weekday() == time.Saturday {
			if ti.Hour() < 12 {
				count := 0
				for _, b := range deck {
					if b.GetSize() == "homebrew" {
						count++
					}
				}
				return count < 2
			}
		}
	}

	return false
}

func (s *Server) refreshStash(ctx context.Context, config *pb.Config) error {
	onDeck := make(map[int64]bool)

	rand.Seed(time.Now().UnixNano())

	found := false
	for _, beer := range config.GetCellar().GetOnDeck() {
		found = found || beer.GetSize() == "stash"
		s.Log(fmt.Sprintf("Found %v in the stash", beer))
		onDeck[beer.Id] = true
	}

	if !found {
		s.Log(fmt.Sprintf("Not adding to stash because: %v", onDeck))
	}

	var chosenBeer *pb.Beer
	chosenIndex := 0
	count := 0
	slot := 0
	opCount := 0
	gtCount := 0
	hsCount := 0
	stashSize := 100
	for sn, c := range config.GetCellar().GetSlots() {
		if c.Accepts == "stash" {

			//Raise an issue on the stash size
			stashSize = len(c.Beers)

			// Randomize the stash for the pull
			rand.Shuffle(len(c.Beers), func(i, j int) { c.Beers[i], c.Beers[j] = c.Beers[j], c.Beers[i] })
			for i, b := range c.Beers {
				if b.GetBreweryId() == 333511 {
					opCount++
				}
				if b.GetBreweryId() == 36950 {
					gtCount++
				}
				if b.GetBreweryId() == 233405 {
					hsCount++
				}
				if !onDeck[b.Id] {
					chosenBeer = b
					chosenIndex = i
					slot = sn
				} else {
					count++
				}
			}
		}
	}

	if opCount < 4 {
		s.RaiseIssue("Buy Original Pattern", fmt.Sprintf("Stocks are low (%v), go buy some Original Patter", opCount))
	}
	if gtCount < 4 {
		s.RaiseIssue("Buy Ghost Town", fmt.Sprintf("Stocks are low (%v), go buy some Ghost Town", opCount))
	}
	if hsCount < 4 {
		s.RaiseIssue("Buy Humble Sea", fmt.Sprintf("Stocks are low (%v), go buy some Humble Sea", opCount))
	}

	if stashSize < 20 {
		s.RaiseIssue("Buy some beer!", fmt.Sprintf("The current size of the stash is %v", stashSize))
	}

	if !found {
		err := s.printer.print(ctx, []string{fmt.Sprintf("%v\n", chosenBeer.Name)})
		if err == nil {
			chosenBeer.DrinkDate = time.Now().Unix()
			chosenBeer.OnDeck = true
			config.GetCellar().OnDeck = append(config.GetCellar().GetOnDeck(), chosenBeer)
			config.GetCellar().GetSlots()[slot].Beers = append(config.GetCellar().GetSlots()[slot].Beers[:chosenIndex], config.GetCellar().GetSlots()[slot].Beers[chosenIndex+1:]...)
		}
		return err
	}
	return nil
}

func findIndex(slot *pb.CellarSlot, date int64) int32 {
	bestIndex := int32(len(slot.Beers))
	for _, b := range slot.Beers {
		if date < b.DrinkDate && b.Index < bestIndex {
			bestIndex = b.Index
		}
	}
	return bestIndex
}

func (s *Server) addBeerToCellar(b *pb.Beer, config *pb.Config) error {
	bestIndex := int32(99)
	bestCellar := -1
	//Do we have room for this beer
	for i, c := range config.GetCellar().GetSlots() {
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
		c := config.Cellar.Slots[bestCellar]
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

func (s *Server) syncDrunk(ctx context.Context, config *pb.Config) error {
	ndrinks, err := s.ut.getLastBeers(mainFetcher{}, mainConverter{}, mainUnmarshaller{})
	if err != nil {
		return err
	}
	s.lastErr = fmt.Sprintf("%v", err)
	return s.addDrunks(ctx, config, ndrinks)
}

func (s *Server) addDrunks(ctx context.Context, config *pb.Config, ndrinks []*pb.Beer) error {
	for _, toadd := range ndrinks {
		found := false
		for _, drunk := range config.Drunk {
			if drunk.GetDrinkDate() == toadd.GetDrinkDate() {
				found = true
			}
		}

		if !found {
			config.Drunk = append(config.Drunk, toadd)
		}
	}

	config.LastSync = time.Now().Unix()
	return s.save(ctx, config)
}

func (s *Server) moveToOnDeck(ctx context.Context, config *pb.Config, t time.Time) error {
	moved := []*pb.Beer{}

	for _, cellar := range config.GetCellar().GetSlots() {
		if cellar.Accepts != "stash" {
			for _, b := range cellar.Beers {
				if b.DrinkDate < t.Unix() {
					moved = append(moved, b)
				}
			}
		}
	}

	if len(moved) > 0 {
		moveStr := []string{}
		for _, b := range moved {
			moveStr = append(moveStr, fmt.Sprintf("%v\n", b.Name))
		}
		err := s.printer.print(ctx, moveStr)
		if err != nil {
			return err
		}
	}

	for _, cellar := range config.GetCellar().GetSlots() {
		if cellar.Accepts != "stash" {
			i := 0
			for i < len(cellar.Beers) {
				if cellar.Beers[i].DrinkDate < t.Unix() {
					cellar.Beers[i].OnDeck = true
					moved = append(moved, cellar.Beers[i])
					config.Cellar.OnDeck = append(config.Cellar.OnDeck, cellar.Beers[i])
					cellar.Beers = append(cellar.Beers[:i], cellar.Beers[i+1:]...)
				} else {
					i++
				}
			}
		}
	}

	return nil
}

func (s *Server) cellarOutOfOrder(ctx context.Context, cellar *pb.CellarSlot) bool {
	for i := range cellar.Beers {
		for j := range cellar.Beers {
			if j > i {
				if (cellar.Beers[i].Index < cellar.Beers[j].Index && cellar.Beers[i].DrinkDate > cellar.Beers[j].DrinkDate) || (cellar.Beers[i].Index > cellar.Beers[j].Index && cellar.Beers[i].DrinkDate < cellar.Beers[j].DrinkDate) {

					return true
				}
			}
		}
	}

	return false
}

func (s *Server) reorderCellar(ctx context.Context, cellar *pb.CellarSlot, cellarNumber int32) {
	sort.SliceStable(cellar.Beers, func(i, j int) bool {
		return cellar.Beers[i].DrinkDate < cellar.Beers[j].DrinkDate
	})
	for i, beer := range cellar.Beers {
		beer.Index = int32(i)
		beer.InCellar = cellarNumber
	}
}

func (s *Server) clearDeck(config *pb.Config) {
	s.lastClean = time.Now()

	for _, bdr := range config.GetDrunk() {
		for i, bde := range config.GetCellar().GetOnDeck() {
			if bdr.Id == bde.Id && bdr.DrinkDate > bde.DrinkDate {
				config.Cellar.OnDeck = append(config.Cellar.OnDeck[:i], config.Cellar.OnDeck[i+1:]...)
			}
		}
	}
}
