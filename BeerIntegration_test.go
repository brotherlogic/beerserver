package main

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestAddMultipleRetrieveMultiple(t *testing.T) {
	tm := time.Now()

	s := InitTestServer(".testaddmultipleretrievemulitple", true)
	s.config.Drunk = append(s.config.Drunk, &pb.Beer{Id: 1234, DrinkDate: tm.Unix() - 10, Size: "bomber"})
	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "bomber"}, Quantity: 3})

	if err != nil {
		t.Fatalf("Error when adding a beer: %v", err)
	}

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})

	if err != nil {
		t.Fatalf("Error when listing beers: %v", err)
	}

	if len(list.Beers) != 3 {
		t.Errorf("Wrong number of beers added: %v", list)
	}

}

func TestAddJR(t *testing.T) {
	s := InitTestServer(".testaddjr", true)
	s.loadDrunk("loaddata/brotherlogic.json")

	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 2407538, Size: "bomber"}, Quantity: 2})
	if err != nil {
		t.Fatalf("Error adding beer: %v", err)
	}

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	for _, b := range list.Beers {
		d := time.Unix(b.DrinkDate, 0)
		if d.Year() == 2018 {
			t.Errorf("Beer placement is wrong %v -> %v", b, d)
		}

		if b.Name != "Drake's Brewing Company - Barrel Aged Jolly Rodger (2017)" {
			t.Errorf("Beer name is wrong: %v", b.Name)
		}
	}
}

func TestAddMultipleWithCorrectDates(t *testing.T) {
	s := InitTestServer(".testaddmultiple", true)
	s.loadDrunk("loaddata/brotherlogic.json")

	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 2324956, Size: "bomber"}, Quantity: 1})

	if err != nil {
		t.Fatalf("Error when adding a beer: %v", err)
	}

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})

	if err != nil {
		t.Fatalf("Error when listing beers: %v", err)
	}

	if len(list.Beers) != 1 || list.Beers[0].DrinkDate != 1553982990 {
		t.Errorf("Beer added has errors: %v", list)
	}

	if len(list.Beers) != 1 {
		t.Errorf("Wrong number of beers added: %v", list)
	}

}

func TestResyncCorrectsSize(t *testing.T) {
	s := InitTestServer(".testresynccorrectssize", true)
	s.loadDrunk("loaddata/brotherlogic.json")

	if len(s.config.Drunk) != 3709 {
		t.Fatalf("Drunk number is wrong: %v", len(s.config.Drunk))
	}

	s.syncDrunk(fileFetcher{})

	if len(s.config.Drunk) != 3710 {
		t.Fatalf("Drunk number is still wrong %v", len(s.config.Drunk))
	}
}

func TestMoveToOnDeck(t *testing.T) {
	s := InitTestServer(".testmovetoondeck", true)
	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 7936, Size: "bomber"}, Quantity: 2})
	if err != nil {
		t.Fatalf("Error adding beer: %v", err)
	}
	time.Sleep(time.Second * 2)
	s.moveToOnDeck(time.Now())

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil {
		t.Fatalf("Error listing beers: %v", err)
	}
	if len(list.Beers) != 1 || !list.Beers[0].OnDeck {
		t.Errorf("List is not correct: %v", list)
	}
}

func TestOrderBeerCorrectly(t *testing.T) {
	s := InitTestServer(".testorderbeercorrectly", true)
	s.loadDrunk("loaddata/brotherlogic.json")

	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 2324956, Size: "bomber"}, Quantity: 1})
	_, err = s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 2428618, Size: "bomber"}, Quantity: 1})

	if err != nil {
		t.Fatalf("Error adding beer: %v", err)
	}

	drinkFirst := int64(0)
	drinkSecond := int64(0)
	for _, c := range s.config.Cellar.Slots {
		for _, b := range c.Beers {
			if b.Index == 0 {
				if drinkFirst != 0 {
					t.Fatalf("Trying to set two firsts: %v", c)
				}
				drinkFirst = b.DrinkDate
			}
			if b.Index == 1 {
				if drinkSecond != 0 {
					t.Fatalf("Trying to set two seconds: %v", c)
				}
				drinkSecond = b.DrinkDate
			}

		}
	}

	if drinkFirst > drinkSecond || drinkFirst == 0 || drinkSecond == 0 {
		t.Errorf("Ordering is incorrect!: %v", s.config.Cellar.Slots)
	}

	log.Printf("%v before %v", time.Unix(drinkFirst, 0), time.Unix(drinkSecond, 0))
}

func TestBeerGoesInRightCellar(t *testing.T) {
	s := InitTestServer(".testbeergoesinrightcellar", true)
	s.config.Cellar.Slots = append(s.config.Cellar.Slots, &pb.CellarSlot{Accepts: "bomber", NumSlots: 10, Beers: []*pb.Beer{&pb.Beer{Id: 2324956, Size: "bomber", DrinkDate: 1234}}})
	s.config.Cellar.Slots = append(s.config.Cellar.Slots, &pb.CellarSlot{Accepts: "bomber", NumSlots: 10, Beers: []*pb.Beer{&pb.Beer{Id: 2324956, Size: "bomber", DrinkDate: time.Now().Unix() + 100}}})

	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 2428618, Size: "bomber"}, Quantity: 1})

	if err != nil {
		t.Fatalf("Error adding the beer: %v", err)
	}

	//New beer should be in the third slot
	if len(s.config.Cellar.Slots[0].Beers) != 0 || len(s.config.Cellar.Slots[1].Beers) != 1 {
		t.Errorf("beer has been added to the wrong slot!: %v or %v or %v", s.config.Cellar.Slots[0], s.config.Cellar.Slots[1], s.config.Cellar.Slots[2])
	}
}
