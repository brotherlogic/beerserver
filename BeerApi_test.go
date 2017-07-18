package main

import (
	"context"
	"log"
	"os"
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

func TestGetName(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 100})
	s.cacheName(1, "madeupname")

	beer, err := s.GetName(context.Background(), &pb.Beer{Id: 1})
	if err != nil {
		t.Fatalf("Error in getting name: %v", err)
	}

	if beer.Name != "madeupname" {
		t.Errorf("Wrong name returned: %v", beer)
	}
}

func TestGetNameOutOfCache(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 7936, Size: "bomber", DrinkDate: 100})

	beer, err := s.GetName(context.Background(), &pb.Beer{Id: 7936})
	if err != nil {
		t.Fatalf("Error in getting name: %v", err)
	}

	if beer.Name != "Firestone Walker Brewing Company - Parabola" {
		t.Errorf("Wrong name returned: %v", beer)
	}
}

func TestGetBeerBeRandom(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 100})
	s.AddBeer(context.Background(), &pb.Beer{Id: 2, Size: "bomber", DrinkDate: 200})

	idSum := int64(0)
	for i := 0; i < 10; i++ {
		beer, err := s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
		if err != nil {
			t.Fatalf("Error getting beer: %v", err)
		}
		log.Printf("Got beer: %v", beer)
		idSum += beer.Id
	}

	if idSum == 10 || idSum == 20 {
		t.Errorf("All the beers have the same id!")
	}
}

func TestRemoveFromCellarTop(t *testing.T) {
	s := GetTestCellar()
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 100})
	s.AddBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber", DrinkDate: 200})

	beer, err := s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
	if err != nil {
		t.Fatalf("Unable to get beer")
	}
	if beer.DrinkDate != 100 {
		t.Errorf("Wrong beer returned: %v", beer)
	}

	_, err = s.RemoveBeer(context.Background(), &pb.Beer{Id: 1, Size: "bomber"})
	if err != nil {
		t.Fatalf("Error in removing beer: %v", err)
	}

	beer, err = s.GetBeer(context.Background(), &pb.Beer{Size: "bomber"})
	if err != nil {
		t.Fatalf("Unable to get beer second time around")
	}
	if beer.DrinkDate != 200 {
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
