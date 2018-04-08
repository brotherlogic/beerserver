package main

import (
	"context"
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

	if len(list.Beers) != 1 || list.Beers[0].DrinkDate != 1522446990 {
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
