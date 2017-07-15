package main

import (
	"context"
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

	return s
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
