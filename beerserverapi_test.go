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
