package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/beerserver/proto"
)

func (s *Server) refreshStash(ctx context.Context) error {
	onDeck := make(map[int64]bool)

	rand.Seed(time.Now().UnixNano())

	found := false
	for _, beer := range s.config.GetCellar().GetOnDeck() {
		found = found || beer.GetSize() == "stash"
		onDeck[beer.Id] = true
	}

	var chosenBeer *pb.Beer
	chosenIndex := 0
	count := 0
	slot := 0
	for sn, c := range s.config.GetCellar().GetSlots() {
		if c.Accepts == "stash" {

			//Raise an issue on the stash size
			if len(c.Beers) < 20 {
				s.RaiseIssue(ctx, "Buy some beer!", fmt.Sprintf("The current size of the stash is %v", len(c.Beers)), false)
			}

			// Randomize the stash for the pull
			rand.Shuffle(len(c.Beers), func(i, j int) { c.Beers[i], c.Beers[j] = c.Beers[j], c.Beers[i] })
			for i, b := range c.Beers {
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

	if !found {
		s.RaiseIssue(ctx, "Adding to stash", fmt.Sprintf("Adding because %v", onDeck), false)
		err := s.printer.print(ctx, []string{fmt.Sprintf("%v\n", chosenBeer.Name)})
		if err == nil {
			chosenBeer.DrinkDate = time.Now().Unix()
			chosenBeer.OnDeck = true
			s.config.GetCellar().OnDeck = append(s.config.GetCellar().GetOnDeck(), chosenBeer)
			s.config.GetCellar().GetSlots()[slot].Beers = append(s.config.GetCellar().GetSlots()[slot].Beers[:chosenIndex], s.config.GetCellar().GetSlots()[slot].Beers[chosenIndex+1:]...)
		}
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

func (s *Server) addBeerToCellar(b *pb.Beer) error {
	bestIndex := int32(99)
	bestCellar := -1
	//Do we have room for this beer
	for i, c := range s.config.GetCellar().GetSlots() {
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

func (s *Server) syncDrunk(ctx context.Context, f httpResponseFetcher) {
	lastID := int32(0)
	for _, drunk := range s.config.GetDrunk() {
		if drunk.CheckinId > lastID {
			lastID = drunk.CheckinId
		}
	}
	ndrinks, err := s.ut.getLastBeers(f, mainConverter{}, mainUnmarshaller{}, lastID)
	s.lastErr = fmt.Sprintf("%v", err)
	s.addDrunks(ctx, err, ndrinks)
}

func (s *Server) addDrunks(ctx context.Context, err error, ndrinks []*pb.Beer) {
	if err == nil {
		s.config.Drunk = append(s.config.Drunk, ndrinks...)
		s.config.LastSync = time.Now().Unix()
		s.save(ctx)
	}
}

func (s *Server) moveToOnDeck(ctx context.Context, t time.Time) error {
	moved := []*pb.Beer{}

	for _, cellar := range s.config.GetCellar().GetSlots() {
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

	for _, cellar := range s.config.GetCellar().GetSlots() {
		if cellar.Accepts != "stash" {
			i := 0
			for i < len(cellar.Beers) {
				if cellar.Beers[i].DrinkDate < t.Unix() {
					cellar.Beers[i].OnDeck = true
					moved = append(moved, cellar.Beers[i])
					s.config.Cellar.OnDeck = append(s.config.Cellar.OnDeck, cellar.Beers[i])
					cellar.Beers = append(cellar.Beers[:i], cellar.Beers[i+1:]...)
				} else {
					i++
				}
			}
		}
	}

	return s.save(ctx)
}

func (s *Server) clearDeck(ctx context.Context) error {
	s.lastClean = time.Now()

	for _, bdr := range s.config.GetDrunk() {
		for i, bde := range s.config.GetCellar().GetOnDeck() {
			if bdr.Id == bde.Id && bdr.DrinkDate > bde.DrinkDate {
				s.config.Cellar.OnDeck = append(s.config.Cellar.OnDeck[:i], s.config.Cellar.OnDeck[i+1:]...)
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

func (s *Server) reorderCellar(ctx context.Context, cellar *pb.CellarSlot) {
	sort.SliceStable(cellar.Beers, func(i, j int) bool {
		return cellar.Beers[i].DrinkDate < cellar.Beers[j].DrinkDate
	})
	for i, beer := range cellar.Beers {
		beer.Index = int32(i)
	}
}
