package main

import (
	"context"
	"testing"

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

	s.Consolidate(context.Background(), &pb.ConsolidateRequest{})

	/*	if len(r.GetConfig().GetCellar().GetSlots()) != 3 {
			t.Fatalf("Wrong number of slots %v", len(r.GetConfig().GetCellar().GetSlots()))
		}

		if len(r.GetConfig().GetCellar().GetSlots()[0].GetBeers()) != 5 {
			t.Errorf("Wrong number of beers overall %v", len(r.GetConfig().GetCellar().GetSlots()[0].GetBeers()))
		}*/
}

func TestAddMultipleRetrieveMultiple(t *testing.T) {
	s := InitTestServer(".testaddmultipleretrievemulitple", true)

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
	if err != nil || len(list.Beers) != 13 {
		t.Fatalf("Bad post refresh list: %v, %v", len(list.Beers), err)
	}

	// One on deck
	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil || len(list.Beers) != 0 {
		t.Fatalf("Bad add to deck: %v, %v", list, err)
	}

}
