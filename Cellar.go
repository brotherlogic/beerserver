package main

import (
	"fmt"
	"log"
	"sort"

	pb "github.com/brotherlogic/beerserver/proto"
)

import "math"

// Printer prints to various places
type Printer interface {
	// Println prints a line
	Println(string)
}

// StdOutPrint a printer which prints to standard out
type StdOutPrint struct{}

// Println prints a line to standard out
func (stdprinter *StdOutPrint) Println(output string) {
	fmt.Printf("%v\n", output)
}

// NewCellar builds a new cellar
func NewCellar(cname string) pb.Cellar {
	return pb.Cellar{Name: cname}
}

// GetFreeSlots Computes the number of free slots for each beer type
func GetFreeSlots(cellar *pb.Cellar) (int, int) {
	if len(cellar.Beers) == 0 {
		return 20, 30
	}

	if cellar.Beers[0].Size == "bomber" {
		return 20 - len(cellar.Beers), 0
	}

	return 0, 30 - len(cellar.Beers)
}

// Remove removes a beer from the cellar
func Remove(cellar *pb.Cellar, id int64) *pb.Beer {
	removeIndex := -1
	for i, v := range cellar.Beers {
		if v.Id == id && removeIndex < 0 {
			removeIndex = i
		}

		if v.Id == id && v.Staged {
			removeIndex = i
		}
	}

	beer := cellar.Beers[removeIndex]
	cellar.Beers = append(cellar.Beers[:removeIndex], cellar.Beers[removeIndex+1:]...)
	return beer
}

// GetRemoveCost computes the cost of removing a beer from the cellar
func GetRemoveCost(cellar *pb.Cellar, id int64) int {
	for _, beer := range cellar.Beers {
		if beer.Id == id && beer.Staged {
			return 0
		}
	}

	for i, beer := range cellar.Beers {
		if beer.Id == id {
			return i + 1
		}
	}

	return -1
}

// CountBeersInCellar counts the number of beers in the cellar
func CountBeersInCellar(cellar *pb.Cellar, id int64) int {
	sum := 0
	for _, v := range cellar.Beers {
		if v.Id == id {
			sum++
		}
	}
	return sum
}

// GetNext Removes the next beer from the cellar
func GetNext(cellar *pb.Cellar) *pb.Beer {
	beer := cellar.Beers[0]
	cellar.Beers = cellar.Beers[1:]
	return beer
}

func getInsertPoint(cellar *pb.Cellar, beer *pb.Beer) int {
	insertPoint := -1
	for i := 0; i < len(cellar.Beers); i++ {
		if beer.DrinkDate < cellar.Beers[i].DrinkDate {
			insertPoint = i
			break
		}
	}

	if insertPoint == -1 {
		insertPoint = len(cellar.Beers)
	}
	return insertPoint
}

// ComputeInsertCost Determines the cost of inserting beer into the cellar
func ComputeInsertCost(cellar *pb.Cellar, beer *pb.Beer) int {
	//Insert cost of an empty cellar should be high
	if len(cellar.Beers) == 0 {
		return int(math.MaxInt16)
	}

	//Don't mix sizes
	if cellar.Beers[0].Size != beer.Size {
		return -1
	}

	// Ensure that cellars don't overflow
	if cellar.Beers[0].Size == "small" && len(cellar.Beers) >= 30 {
		return -1
	} else if cellar.Beers[0].Size == "bomber" && len(cellar.Beers) >= 20 {
		return -1
	}

	insertPoint := getInsertPoint(cellar, beer)

	return insertPoint
}

// AddBeerToCellar adds a beer to the cellar
func AddBeerToCellar(cellar *pb.Cellar, beer *pb.Beer) {
	log.Printf("Adding %v to %v", beer, len(cellar.Beers))

	insertPoint := getInsertPoint(cellar, beer)

	log.Printf("Found insert point: %v", insertPoint)
	log.Printf("WAS %v", cellar.Beers)
	nb := make([]*pb.Beer, 0)
	nb = append(nb, cellar.Beers[0:insertPoint]...)
	log.Printf("NB %v", nb)
	nb = append(nb, beer)
	log.Printf("MB %v", nb)
	nb = append(nb, cellar.Beers[insertPoint:]...)
	log.Printf("OB %v", cellar.Beers)
	cellar.Beers = nb
	log.Printf("NOW %v", cellar.Beers)
}

// MergeCellars combines N cellars into a list of beers
func MergeCellars(bsize string, cellars ...*pb.Cellar) []*pb.Beer {
	var retArr []*pb.Beer
	for _, cellar := range cellars {
		if len(cellar.Beers) > 0 && cellar.Beers[0].Size == bsize {
			retArr = append(retArr, cellar.Beers...)
		}
	}

	sort.Sort(ByDate(retArr))

	return retArr
}
