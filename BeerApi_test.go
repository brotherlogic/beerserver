package main

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	pb "github.com/brotherlogic/beerserver/proto"
	"github.com/brotherlogic/keystore/client"
)

func GetTestCellar() Server {
	s := Init()
	s.cellar = NewBeerCellar("testcellar")
	os.RemoveAll("testcellar")
	s.GoServer.KSclient = *keystoreclient.GetTestClient("testcellar")
	s.SkipLog = true

	s.ut = GetTestUntappd()

	return s
}

func TestGetDrunks(t *testing.T) {
	s := GetTestCellar()

	s.AddBeer(context.Background(), &pb.Beer{Id: 1770600, Size: "bomber", DrinkDate: 100})
	s.cellar.SyncTime = 0
	Sync(s.ut, s.cellar)

	drunk, err := s.GetDrunk(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Error getting drunks: %v", err)
	}

	if len(drunk.Beers) != 1 {
		t.Errorf("Beer list has come back empty: %v", drunk)
	}
}

func TestGetNameOutOfCache(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 7936, Size: "bomber", DrinkDate: 100})
	beer := s.cellar.Cellars[0].Beers[0]

	if beer.Name != "Firestone Walker Brewing Company - Parabola" || beer.Abv != 14 {
		t.Errorf("Wrong name returned: %v", beer)
	}
}

func TestGetCellarWithNames(t *testing.T) {
	s := GetTestCellar()
	s.cellar.Cellars = append(s.cellar.Cellars, &pb.Cellar{Beers: make([]*pb.Beer, 0), Name: "cellar1"})
	s.cellar.Cellars[0].Beers = append(s.cellar.Cellars[0].Beers, &pb.Beer{Id: 7936, Size: "bomber", DrinkDate: 100})

	cellar, err := s.GetCellar(context.Background(), &pb.Cellar{Name: "cellar1"})
	if err != nil {
		t.Errorf("Unable to pull cellar: %v", cellar)
	}

	if cellar.Beers[0].Name != "Firestone Walker Brewing Company - Parabola" || cellar.Beers[0].Abv != 14 {
		t.Errorf("Cellar pull hasn't refreshed anemes: %v", cellar.Beers[0])
	}
}

func TestRefreshBeerName(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 7936, Size: "bomber", DrinkDate: 100, Name: "This API key has reached their API limit for the hour. Please wait before making another call."})

	beer, err := s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})

	if err != nil {
		t.Fatalf("Error in Get beer: %v", err)
	}

	if strings.Contains(beer.Name, "API") {
		t.Errorf("Beer has not been renamed: %v", beer)
	}
}

func TestRefreshBeerNameFromCellar(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 7936, Size: "bomber", DrinkDate: 100, Name: "This API key has reached their API limit for the hour. Please wait before making another call."})

	cellar, err := s.GetCellar(context.Background(), &pb.Cellar{Name: "cellar1"})
	if err != nil {
		log.Fatalf("Error in cellar pull: %v", err)
	}

	if strings.Contains(cellar.GetBeers()[0].Name, "API") {
		t.Errorf("Error in pulling cellar: %v", cellar)
	}
}

func TestRemoveFromCellarTop(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 100})
	s.AddBeer(context.Background(), &pb.Beer{Id: 2, Size: "bomber", DrinkDate: 200})

	beer, err := s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
	if err != nil {
		t.Fatalf("Unable to get beer")
	}

	bID := beer.Id

	_, err = s.RemoveBeer(context.Background(), &pb.Beer{Id: bID, Size: "bomber"})
	if err != nil {
		t.Fatalf("Error in removing beer: %v", err)
	}

	beer, err = s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
	if err != nil {
		t.Fatalf("Unable to get beer second time around")
	}
	if beer.Id == bID {
		t.Errorf("Wrong beer returned second: %v", beer)
	}
}

func TestGetCellar(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 100})
	cellar, err := s.GetCellar(context.Background(), &pb.Cellar{Name: "cellar1"})

	if err != nil {
		t.Fatalf("Unable to retrieve cellar %v", err)
	}

	if len(cellar.Beers) != 1 || cellar.Beers[0].Id != 1 {
		t.Errorf("Error in retrieving beer: %v", cellar)
	}
}

func GetTestCellarFail() Server {
	s := Init()
	s.cellar = NewBeerCellar("testcellar")
	os.RemoveAll("testcellar")
	s.GoServer.KSclient = *keystoreclient.GetTestClient("testcellar")
	s.SkipLog = true

	return s
}

func TestGetCellarNotThere(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 100})
	cellar, err := s.GetCellar(context.Background(), &pb.Cellar{Name: "cellar9"})

	if err == nil {
		t.Fatalf("Cellar retrieved even though it's not there %v", cellar)
	}
}
