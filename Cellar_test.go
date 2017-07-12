package main

import "log"
import "math"
import "testing"

type LineCounterPrinter struct {
	lcount int
}

func (lineCounter *LineCounterPrinter) Println(output string) {
	lineCounter.lcount = lineCounter.lcount + 1
}

func (lineCounter LineCounterPrinter) GetCount() int {
	return lineCounter.lcount
}

func TestCellarFree(t *testing.T) {
	mine1 := NewCellar("freetest1")
	mine2 := NewCellar("freetest2")
	mine3 := NewCellar("freetest3")

	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1235~01/03/16~small")

	AddBeerToCellar(&mine1, &beer1)
	AddBeerToCellar(&mine1, &beer2)
	AddBeerToCellar(&mine2, &beer3)

	large1, small1 := GetFreeSlots(&mine1)
	if large1 != 20-2 || small1 != 0 {
		t.Errorf("Problem computing free slots of bomber cellar: %v, %v, %v", large1, small1, mine1)
	}
	large2, small2 := GetFreeSlots(&mine2)
	if large2 != 0 || small2 != 30-1 {
		t.Errorf("Problem computing free slots of bomber cellar: %v, %v, %v", large2, small2, mine2)
	}
	large3, small3 := GetFreeSlots(&mine3)
	if large3 != 20 || small3 != 30 {
		t.Errorf("Problem computing free slots of bomber cellar: %v, %v, %v", large3, small3, mine3)
	}
}

func TestRemoveBeer(t *testing.T) {
	mine := NewCellar("testremovecostcellar")
	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1235~01/02/16~small")
	beer3, _ := NewBeer("1236~01/03/16~small")
	AddBeerToCellar(&mine, &beer1)
	AddBeerToCellar(&mine, &beer2)
	AddBeerToCellar(&mine, &beer3)

	Remove(&mine, 1235)

	if len(mine.Beers) != 2 {
		t.Errorf("Beer has not been removed: %v\n", mine)
	}
}

func TestRemoveCost(t *testing.T) {
	mine := NewCellar("testremovecostcellar")
	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1235~01/02/16~small")
	beer3, _ := NewBeer("1236~01/03/16~small")
	AddBeerToCellar(&mine, &beer1)
	AddBeerToCellar(&mine, &beer2)
	AddBeerToCellar(&mine, &beer3)
	if GetRemoveCost(&mine, 1237) >= 0 {
		t.Errorf("Removing non cellared beer is not less than zero: %v\n", GetRemoveCost(&mine, 1237))
	}

	if GetRemoveCost(&mine, 1234) != 0 {
		t.Errorf("Remove cost is wrong (0): %v\n", GetRemoveCost(&mine, 1234))
	}
	if GetRemoveCost(&mine, 1235) != 1 {
		t.Errorf("Remove cost is wrong (1): %v\n", GetRemoveCost(&mine, 1234))
	}
	if GetRemoveCost(&mine, 1236) != 2 {
		t.Errorf("Remove cost is wrong (2): %v\n", GetRemoveCost(&mine, 1234))
	}

}

func TestMidInsert(t *testing.T) {
	mine := NewCellar("testinsertcellar")
	beer1, _ := NewBeer("1234~01/01/16~small")
	beer2, _ := NewBeer("1235~01/02/16~small")
	beer3, _ := NewBeer("1236~01/03/16~small")
	beer4, _ := NewBeer("1237~01/04/16~small")

	AddBeerToCellar(&mine, &beer1)
	AddBeerToCellar(&mine, &beer2)
	AddBeerToCellar(&mine, &beer4)

	log.Printf("1. %v\n", mine)

	AddBeerToCellar(&mine, &beer3)

	log.Printf("2. %v\n", mine)

	if CountBeersInCellar(&mine, 1237) != 1 {
		t.Errorf("Problem with cellar: %v (%v)\n", mine, CountBeersInCellar(&mine, 1237))
	}
}

func TestProblemsOfClobbering(t *testing.T) {
	mine1 := NewCellar("testprobcellar")
	beer1, _ := NewBeer("938229~01/03/18~small")
	AddBeerToCellar(&mine1, &beer1)
	log.Printf("%v\n", mine1)
	if CountBeersInCellar(&mine1, 938229) != 1 {
		t.Errorf("Problem with cellar 1. (%v) %v\n", CountBeersInCellar(&mine1, 938229), mine1)
	}
	beer2, _ := NewBeer("938229~01/03/17~small")
	AddBeerToCellar(&mine1, &beer2)
	log.Printf("%v\n", mine1)
	if CountBeersInCellar(&mine1, 938229) != 2 {
		t.Errorf("Problem with cellar 2. (%v) %v\n", CountBeersInCellar(&mine1, 938229), mine1)
	}
	beer3, _ := NewBeer("938229~01/03/16~small")
	AddBeerToCellar(&mine1, &beer3)
	log.Printf("%v\n", mine1)
	if CountBeersInCellar(&mine1, 938229) != 3 {
		t.Errorf("Problem with cellar 3. (%v) %v\n", CountBeersInCellar(&mine1, 938229), mine1)
	}
	beer4, _ := NewBeer("768356~01/09/18~small")
	AddBeerToCellar(&mine1, &beer4)
	log.Printf("%v\n", mine1)
	beer5, _ := NewBeer("768356~01/09/18~small")
	AddBeerToCellar(&mine1, &beer5)
	log.Printf("%v\n", mine1)
	beer6, _ := NewBeer("768356~01/09/17~small")
	AddBeerToCellar(&mine1, &beer6)
	log.Printf("%v\n", mine1)
	beer7, _ := NewBeer("938229~01/03/17~small")
	AddBeerToCellar(&mine1, &beer7)
	log.Printf("%v\n", mine1)
	if CountBeersInCellar(&mine1, 938229) != 4 {
		t.Errorf("Problem with cellar 4. (%v) %v\n", CountBeersInCellar(&mine1, 938229), mine1)
	}
	beer8, _ := NewBeer("552346~01/03/17~small")
	AddBeerToCellar(&mine1, &beer8)
	log.Printf("%v\n", mine1)

	if len(mine1.Beers) != 8 {
		t.Errorf("Not the right number of beers, 8 but %v\n", len(mine1.Beers))
	}
}

func TestMergeCellar(t *testing.T) {
	cellar1 := NewCellar("cellar1")
	cellar2 := NewCellar("cellar2")
	cellar3 := NewCellar("cellar3")

	beer1, _ := NewBeer("1234~10/01/15~bomber")
	AddBeerToCellar(&cellar1, &beer1)
	beer2, _ := NewBeer("1235~01/01/15~bomber")
	AddBeerToCellar(&cellar2, &beer2)
	beer3, _ := NewBeer("1236~05/01/15~bomber")
	AddBeerToCellar(&cellar3, &beer3)

	merged := MergeCellars("bomber", &cellar1, &cellar2, &cellar3)
	merged2 := MergeCellars("small", &cellar1, &cellar2, &cellar3)

	log.Printf("GOT MERGED: %v", merged)
	log.Printf("GOT MERGED2: %v", merged2)

	if merged[0].Id != 1235 {
		t.Errorf("Merged list is ordered incorrectly %v\n", merged[0])
	}

	if len(merged2) != 0 {
		t.Errorf("Merged small list is non-empty: %v\n", merged2)
	}
}

func TestPrint(t *testing.T) {
	printer := StdOutPrint{}

	//Line below should not fail
	printer.Println("Made up")
}

func TestCellarOverflowSmall(t *testing.T) {
	cellar := NewCellar("testing_cellar")

	//Should be able to add 30 small bottles before insert cost is < 0
	for i := 0; i < 30; i++ {
		beer1, _ := NewBeer("1234~01/01/15~small")
		cost := ComputeInsertCost(&cellar, &beer1)
		if cost < 0 {
			t.Errorf("Inserting %v into %v is too costly\n", beer1, cellar)
		}
		AddBeerToCellar(&cellar, &beer1)
	}

	beer2, _ := NewBeer("1234~01/01/15~small")
	cost := ComputeInsertCost(&cellar, &beer2)
	if cost >= 0 {
		t.Errorf("Inserting %v into %v is not prohibited (%v)\n", beer2, cellar, cost)
	}
}

func TestCellarOverflowBomber(t *testing.T) {
	cellar := NewCellar("testing_cellar")

	//Should be able to add 30 small bottles before insert cost is < 0
	for i := 0; i < 20; i++ {
		beer1, _ := NewBeer("1234~01/01/15~bomber")
		cost := ComputeInsertCost(&cellar, &beer1)
		if cost < 0 {
			t.Errorf("Inserting %v into %v is too costly\n", beer1, cellar)
		}
		AddBeerToCellar(&cellar, &beer1)
	}

	beer2, _ := NewBeer("1234~01/01/15~bomber")
	cost := ComputeInsertCost(&cellar, &beer2)
	if cost >= 0 {
		t.Errorf("Inserting %v into %v is not prohibited (%v)\n", beer2, cellar, cost)
	}
}

func TestComputeMixSizes(t *testing.T) {
	cellar := NewCellar("testing_cellar")
	beer1, err1 := NewBeer("1234~01/01/16~bomber")
	beer2, err2 := NewBeer("1235~01/01/16~small")

	if err1 != nil || err2 != nil {
		t.Errorf("Parse issue %v,%v\n", err1, err2)
	}

	AddBeerToCellar(&cellar, &beer1)

	cost := ComputeInsertCost(&cellar, &beer2)
	if cost >= 0 {
		t.Errorf("Beer sizes have been mixed: %v\n", cost)
	}
}

func TestComputeEmptyCellarCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer, _ := NewBeer("1234~01/01/16")

	ic := ComputeInsertCost(&cellar, &beer)

	if ic != math.MaxInt16 {
		t.Errorf("Insert cost should be really high here; in fact it's %d\n", ic)
	}
}

func TestComputeBackCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")

	AddBeerToCellar(&cellar, &beer1)
	ic := ComputeInsertCost(&cellar, &beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestComputeMiddleCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	AddBeerToCellar(&cellar, &beer1)
	AddBeerToCellar(&cellar, &beer3)
	ic := ComputeInsertCost(&cellar, &beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestCellarInsert(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	AddBeerToCellar(&cellar, &beer1)
	AddBeerToCellar(&cellar, &beer2)
	AddBeerToCellar(&cellar, &beer3)

	beert1 := GetNext(&cellar)
	beert2 := GetNext(&cellar)
	beert3 := GetNext(&cellar)

	if beert1 != &beer1 {
		t.Errorf("Beer1 has been inserted in the wrong position\n")
	}
	if beert2 != &beer2 {
		t.Errorf("Beer2 has been inserted in the wrong position\n")
	}
	if beert3 != &beer3 {
		t.Errorf("Beer3 has been inserted in the wrong position\n")
	}
}
