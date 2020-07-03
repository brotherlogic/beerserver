package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestAddBeerToFullCellar(t *testing.T) {
	s := InitTestServer(".testaddbeertofullcellar", true)
	s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "bomber"}, Quantity: 100})
	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(list.Beers) != 10 {
		t.Errorf("Wrong number of beers: %v", len(list.Beers))
	}
}

func TestDeleteBeer(t *testing.T) {
	s := InitTestServer(".testdeletebeer", true)

	s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "bomber"}, Quantity: 1})
	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})

	if err != nil {
		t.Fatalf("Error in listing beers: %v", err)
	}

	if len(list.GetBeers()) == 0 {
		t.Fatalf("Bad return: %v", list)
	}

	uid := list.GetBeers()[0].GetUid()
	_, err = s.DeleteBeer(context.Background(), &pb.DeleteBeerRequest{Uid: uid})
	if err != nil {
		t.Fatalf("Error in deleting beer: %v", err)
	}

	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	if err != nil {
		t.Fatalf("Error in listing beers: %v", err)
	}

	if len(list.Beers) != 0 {
		t.Errorf("Beer was not deleted: %v", list.Beers)
	}
}

func TestDeleteBeerOnDeck(t *testing.T) {
	s := InitTestServer(".testdeletebeer", true)

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})

	if err != nil {
		t.Fatalf("Error in listing beers: %v", err)
	}

	if len(list.GetBeers()) == 0 {
		t.Fatalf("Bad list: %v", list)
	}

	uid := list.GetBeers()[0].GetUid()
	_, err = s.DeleteBeer(context.Background(), &pb.DeleteBeerRequest{Uid: uid})
	if err != nil {
		t.Fatalf("Error in deleting beer: %v", err)
	}

	list, err = s.ListBeers(context.Background(), &pb.ListBeerRequest{OnDeck: true})
	if err != nil {
		t.Fatalf("Error in listing beers: %v", err)
	}

	if len(list.Beers) != 0 {
		t.Errorf("Beer was not deleted: %v", list.Beers)
	}
}

func TestDeleteBeerWithEmptyRequest(t *testing.T) {
	s := InitTestServer(".testdeletebeer", true)

	s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "bomber"}, Quantity: 1})
	_, err := s.DeleteBeer(context.Background(), &pb.DeleteBeerRequest{})
	if err == nil {
		t.Errorf("Empty request did not fail")
	}
}

func TestDeleteBeerWithWrongRequest(t *testing.T) {
	s := InitTestServer(".testdeletebeer", true)

	s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "bomber"}, Quantity: 1})
	_, err := s.DeleteBeer(context.Background(), &pb.DeleteBeerRequest{Uid: 1})
	if err == nil {
		t.Errorf("Bad request did not fail")
	}
}

func TestUidClash(t *testing.T) {
	s := InitTestServer(".testdeletebeer", true)

	s.AddBeer(context.Background(), &pb.AddBeerRequest{Beer: &pb.Beer{Id: 1234, Size: "bomber"}, Quantity: 2})

	list, err := s.ListBeers(context.Background(), &pb.ListBeerRequest{})
	if err != nil {
		t.Errorf("Unable to list beers: %v", err)
	}

	if len(list.Beers) != 2 {
		t.Fatalf("Wrong number of beers: %v", list.Beers)
	}

	if list.Beers[0].Uid == list.Beers[1].Uid {
		t.Errorf("UID clash! %v", list.Beers)
	}
}
