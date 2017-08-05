package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestStagingThing(t *testing.T) {
	s := GetTestCellar()

	s.AddBeer(context.Background(), &pb.Beer{Id: 1770600, Size: "bomber", DrinkDate: 100})
	beer, err := s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
	if err != nil || beer.Id != 1770600 {
		t.Fatalf("Beer has come back bad: %v,%v", beer, err)
	}
	s.AddBeer(context.Background(), &pb.Beer{Id: 123, Size: "bomber", DrinkDate: 100})

	//Beer returned should *always* be the staged beer
	for i := 0; i < 200; i++ {
		beer, err = s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
		if err != nil || beer.Id != 1770600 || !beer.Staged {
			t.Fatalf("Beer has come back bad when staged: %v,%v", beer, err)
		}
	}

	s.cellar.SyncTime = 0
	Sync(s.ut, s.cellar)

	beer, err = s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
	if err != nil || beer.Id != 123 {
		t.Errorf("Beer has come back bad second time: %v,%v", beer, err)
	}
}
