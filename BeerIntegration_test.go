package main

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestConsolidate(t *testing.T) {
	s := InitTestServer(".testconsolidate", true)

	config := &pb.Config{Cellar: &pb.Cellar{
		Slots: []*pb.CellarSlot{
			&pb.CellarSlot{
				Accepts:  "small",
				NumSlots: 5,
				Beers: []*pb.Beer{
					&pb.Beer{Id: 1},
				},
			},
			&pb.CellarSlot{
				Accepts:  "bomber",
				NumSlots: 5,
				Beers: []*pb.Beer{
					&pb.Beer{Id: 2},
				},
			},
			&pb.CellarSlot{
				Accepts:  "small",
				NumSlots: 10,
				Beers: []*pb.Beer{
					&pb.Beer{Id: 3},
					&pb.Beer{Id: 3},
					&pb.Beer{Id: 3},
					&pb.Beer{Id: 3},
					&pb.Beer{Id: 3},
				},
			},
			&pb.CellarSlot{
				Accepts:  "bomber",
				NumSlots: 10,
				Beers: []*pb.Beer{
					&pb.Beer{Id: 4},
					&pb.Beer{Id: 4},
				},
			},
		},
	}}
	s.save(context.Background(), config)

	r, _ := s.Consolidate(context.Background(), &pb.ConsolidateRequest{})
	log.Printf("%v", r)

	/*	if len(r.GetConfig().GetCellar().GetSlots()) != 3 {
			t.Fatalf("Wrong number of slots %v", len(r.GetConfig().GetCellar().GetSlots()))
		}

		if len(r.GetConfig().GetCellar().GetSlots()[0].GetBeers()) != 5 {
			t.Errorf("Wrong number of beers overall %v", len(r.GetConfig().GetCellar().GetSlots()[0].GetBeers()))
		}*/
}

func TestAddMultipleRetrieveMultiple(t *testing.T) {
	tm := time.Now()

	s := InitTestServer(".testaddmultipleretrievemulitple", true)
	config := &pb.Config{}
	config.Drunk = append(config.Drunk, &pb.Beer{Id: 1234, DrinkDate: tm.Unix() - 10, Size: "bomber"})
	s.save(context.Background(), config)

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

func TestBeerDateLatest(t *testing.T) {
	s := InitTestServer(".testbeerdatelatest", true)

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
	_, err = s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 12345, Size: "stash", BreweryId: 333511}, Quantity: 3})

	if err != nil {
		t.Fatalf("Error adding beer")
	}

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	if err != nil || len(list.Beers) != 13 {
		t.Fatalf("Bad initial add: %v, %v", list, err)
	}

	//Run the stash twice
	s.Update(context.Background(), &pb.UpdateRequest{})

	// Nineteen in the cellar
	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	if err != nil || len(list.Beers) != 12 {
		t.Fatalf("Bad post refresh list: %v, %v", len(list.Beers), err)
	}

	// One on deck
	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil || len(list.Beers) != 1 {
		t.Fatalf("Bad add to deck: %v, %v", list, err)
	}

}
