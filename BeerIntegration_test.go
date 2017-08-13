package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestRemoveStagingBeer(t *testing.T) {
	s := GetTestCellar()

	for i := 0; i < 10; i++ {
		s.AddBeer(context.Background(), &pb.Beer{Id: 1770600, Size: "bomber", DrinkDate: 100})
	}
	//Get and remove each - they should all be staged for removal
	for i := 0; i < 10; i++ {
		beer, err := s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
		if err != nil {
			t.Fatalf("Error in getting beer: %v", err)
		}
		if beer.Id != 1770600 || !beer.Staged {
			t.Fatalf("Got beer is incorrect: %v", beer)
		}

		beer, err = s.RemoveBeer(context.Background(), &pb.Beer{Id: 1770600, Staged: true})
		if err != nil {
			t.Fatalf("Error in removing beer: %v", err)
		}
		if beer.Id != 1770600 || !beer.Staged {
			t.Fatalf("Removed beer is incorrect: %v", beer)
		}
	}
}

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
