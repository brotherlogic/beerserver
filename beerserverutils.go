package main

import (
	"fmt"
	"log"

	pb "github.com/brotherlogic/beerserver/proto"
)

func findIndex(slot *pb.CellarSlot, date int64) int32 {
	return 0
}

func (s *Server) addBeerToCellar(b *pb.Beer) error {
	//Do we have room for this beer
	for _, c := range s.config.Cellar.Slots {
		log.Printf("WHAT %v and %v", c.Accepts, b.Size)
		if c.Accepts == b.Size {
			log.Printf("HERE %v and %v", c.NumSlots, c.Beers)
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
