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

	if len(s.config.Drunk) != 4099 {
		t.Fatalf("Drunk number is wrong: %v", len(s.config.Drunk))
	}

	s.syncDrunk(context.Background(), fileFetcher{})

	if len(s.config.Drunk) != 4099 {
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
	s.moveToOnDeck(context.Background(), time.Now())

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil {
		t.Fatalf("Error listing beers: %v", err)
	}
	if len(list.Beers) != 1 || !list.Beers[0].OnDeck {
		t.Errorf("List is not correct: %v", list)
	}
}

func TestMoveToOnDeckTwoCellars(t *testing.T) {
	s := InitTestServer(".testmovetoondeck", true)
	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 7936, Size: "bomber"}, Quantity: 2})
	_, err = s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 2324956, Size: "small"}, Quantity: 2})
	if err != nil {
		t.Fatalf("Error adding beer: %v", err)
	}
	time.Sleep(time.Second * 2)
	s.moveToOnDeck(context.Background(), time.Now())
	s.moveToOnDeck(context.Background(), time.Now())

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil {
		t.Fatalf("Error listing beers: %v", err)
	}
	if len(list.Beers) != 2 || !list.Beers[0].OnDeck {
		t.Errorf("List is not correct: %v", list)
	}

	log.Printf("HERE = %v", list.Beers)
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
	if len(s.config.Cellar.Slots[0].Beers) != 0 || len(s.config.Cellar.Slots[3].Beers) != 1 {
		t.Errorf("beer has been added to the wrong slot!: %v or %v or %v", s.config.Cellar.Slots[0], s.config.Cellar.Slots[1], s.config.Cellar.Slots[2])
	}
}

func TestBeerDateLatest(t *testing.T) {
	s := InitTestServer(".testbeerdatelatest", true)
	s.loadDrunk("loaddata/brotherlogic.json")

	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 5771, Size: "bomber"}, Quantity: 1})
	if err != nil {
		t.Fatalf("Error addding beer: %v", err)
	}

	//New beer should have drink date of 16th November 2018
	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})

	if len(list.Beers) != 1 {
		t.Fatalf("Wrong number of beers: %v", list)
	}

	date := time.Unix(list.Beers[0].DrinkDate, 0)
	if date.Year() != 2019 || date.Month() != time.April || date.Day() != 15 {
		t.Errorf("Wrong date returned: %v", date)
	}

}

func TestStashProcess(t *testing.T) {
	s := InitTestServer(".teststashprocess", true)
	_, err := s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "stash"}, Quantity: 10})
	_, err = s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 12345, Size: "stash"}, Quantity: 10})

	if err != nil {
		t.Fatalf("Error adding beer")
	}

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	if err != nil || len(list.Beers) != 20 {
		t.Fatalf("Bad initial add: %v, %v", list, err)
	}

	//Run the stash twice
	s.refreshStash(context.Background())
	s.refreshStash(context.Background())

	// Nineteen in the cellar
	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	if err != nil || len(list.Beers) != 19 {
		t.Fatalf("Bad post refresh list: %v, %v", len(list.Beers), err)
	}

	// One on deck
	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil || len(list.Beers) != 1 {
		t.Fatalf("Bad add to deck: %v, %v", list, err)
	}

}
