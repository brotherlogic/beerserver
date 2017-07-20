package main

import (
	"log"
	"testing"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestCountBeers(t *testing.T) {
	mine := NewBeerCellar("testaddbydays")
	AddBeerByDays(mine, "1234", "01/01/16", "bomber", "14", "5")
	AddBeerByDays(mine, "2234", "01/01/16", "small", "14", "7")

	if CountBeers(mine, 1234) != 5 {
		t.Errorf("Wrong number of 1234: %v", mine)
	}

	if CountBeers(mine, 2234) != 7 {
		t.Errorf("Wrong number of 2234: %v", mine)
	}
}

func TestAddToEmpty(t *testing.T) {
	mine := NewBeerCellar("testempty")
	log.Printf("Built new cellar")
	c := AddBuilt(mine, &pb.Beer{Id: 234, Size: "bomber"})
	if c == nil {
		t.Errorf("Error adding beer to empty cellar")
	}
}

func TestRemoveFromFull(t *testing.T) {
	mine := NewBeerCellar("testremovefromfull")

	for i := 0; i < 20; i++ {
		AddBuilt(mine, &pb.Beer{Size: "bomber", Id: int64(i + 2), DrinkDate: int64(i + 3)})
	}

	if len(mine.Cellars[0].Beers) != 20 {
		t.Fatalf("Error in adding beers: %v", mine)
	}

	log.Printf("CELLAR: %v", mine)

	IDToRemove := mine.Cellars[0].Beers[0].Id
	RemoveBeer(mine, IDToRemove)

	if len(mine.Cellars[0].Beers) != 19 {
		t.Errorf("beer has not been removed: %v", mine)
	}

	log.Printf("NOWCELLAR: %v", mine)
}

func TestCountFreeSlots(t *testing.T) {
	mine := NewBeerCellar("testfreeslots")
	AddBeerByDays(mine, "1234", "01/01/06", "bomber", "14", "5")

	largeFree, smallFree := GetTotalFreeSlots(mine)
	if smallFree != 7*30 {
		t.Errorf("Wrong count of small free: %v (expeceted %v)", smallFree, 7*30)
	}
	if largeFree != 7*20+20-5 {
		t.Errorf("Wrong count of large free: %v (expeceted %v)", largeFree, 7*20+20-5)
	}
}

func TestAddByDate(t *testing.T) {
	mine := NewBeerCellar("testaddbydays")
	AddBeerByDays(mine, "1234", "01/01/16", "bomber", "14", "5")

	if Size(mine) != 5 {
		t.Errorf("Not enough beers added: %v\n", Size(mine))
	}

	if mine.Cellars[0].Beers[0].DrinkDate != 1451606400 {
		t.Errorf("Date on first entry is wrong: %v\n", mine.Cellars[0].Beers[0].DrinkDate)
	}

	if mine.Cellars[0].Beers[1].DrinkDate != 1452816000 {
		t.Errorf("Date on second entry is wrong: %v\n See cellar %v", mine.Cellars[0].Beers[1].DrinkDate, mine.Cellars[0])
	}
}

func TestListBeers(t *testing.T) {
	mine := NewBeerCellar("testlistbeers")
	AddBeer(mine, "1234", 7254, "bomber")
	AddBeer(mine, "1234", 7256, "bomber")

	beers := ListBeers(mine, 2, "bomber", 7255)
	if len(beers) != 1 || beers[0].Id != 1234 {
		t.Errorf("Returned beers are not correct: %v\n", beers)
	}
}

func TestAddByYear(t *testing.T) {
	mine := NewBeerCellar("testaddbydays")
	AddBeerByYears(mine, "1234", "01/01/16", "bomber", "1", "5")

	if Size(mine) != 5 {
		t.Errorf("Not enough beers added: %v\n", Size(mine))
	}

	if mine.Cellars[0].Beers[0].DrinkDate != 1451606400 {
		t.Errorf("Date on first entry is wrong: %v\n", mine.Cellars[0].Beers[0].DrinkDate)
	}

	if mine.Cellars[0].Beers[1].DrinkDate != 1483228800 {
		t.Errorf("Date on second entry is wrong: %v\n See cellar %v", mine.Cellars[0].Beers[1].DrinkDate, mine.Cellars[0])
	}

}

func TestAddToCellars(t *testing.T) {
	mine1 := NewBeerCellar("testaddcellar")
	AddBeer(mine1, "1234", 1233, "bomber")
	AddBeer(mine1, "1234", 1234, "bomber")

	if GetEmptyCellarCount(mine1) != 7 {
		t.Errorf("Cellar is not balanced: %v\n", mine1)
	}
}

func TestSyncAndSave(t *testing.T) {
	mine := NewBeerCellar("testsync")
	AddBeer(mine, "1234", 12342, "bomber")
	//This one should be pulled from the venue list
	AddBeer(mine, "1770600", 12343, "bomber")
	AddBeer(mine, "1234", 12344, "bomber")
	mine.SyncTime = 12345

	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	Sync(u, mine, fileFetcher{}, mainConverter{})
	if Size(mine) != 2 {
		t.Errorf("Sync remove has failed(%v): %v\n", Size(mine), mine)
	}

	if mine.SyncTime == 12345 {
		t.Errorf("Sync time was not updated")
	}
}

func TestRemoveFromCellar(t *testing.T) {
	mine := NewBeerCellar("testremovefromcellar")
	AddBeer(mine, "1234", 1234, "bomber")
	AddBeer(mine, "1235", 1235, "bomber")
	AddBeer(mine, "1234", 1236, "bomber")

	RemoveBeer(mine, 1234)

	if len(mine.Cellars[0].Beers) != 2 {
		t.Errorf("Beer has not been removed: %v\n", mine)
	}
}

func TestAddNoBeer(t *testing.T) {
	mine := NewBeerCellar("test")
	AddBeer(mine, "-1", 1234, "bomber")

	if Size(mine) != 0 {
		t.Errorf("Error on adding no beer - %v\n", Size(mine))
	}
}

func TestMin(t *testing.T) {
	if Min(3, 2) == 3 || Min(2, 3) == 3 {
		t.Errorf("Min is not returning the min\n")
	}
}

func TestGetNumberOfCellars(t *testing.T) {
	bc := NewBeerCellar("test")
	if GetNumberOfCellars(bc) != 8 {
		t.Errorf("Wrong number of cellars: %d\n", GetNumberOfCellars(bc))
	}
}

func TestAddToCellar(t *testing.T) {
	cellar := NewBeerCellar("test")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	AddBuilt(cellar, &beer1)
	AddBuilt(cellar, &beer2)
	AddBuilt(cellar, &beer3)

	if GetEmptyCellarCount(cellar) != 7 {
		t.Errorf("Too many cellars are not empty %d\n", GetEmptyCellarCount(cellar))
	}
}
