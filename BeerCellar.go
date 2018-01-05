package main

import (
	"strconv"
	"strings"

	pb "github.com/brotherlogic/beerserver/proto"
)

import "time"

const (
	// NumCellars the number of cellars
	NumCellars = 4
)

// GetTotalFreeSlots Computes the number of free slots for each beer type
func GetTotalFreeSlots(cellar *pb.BeerCellar) (int, int) {
	sumLarge, sumSmall := 0, 0
	for _, cell := range cellar.Cellars {
		large, small := GetFreeSlots(cell)
		sumSmall += small
		sumLarge += large
	}
	return sumLarge, sumSmall
}

// Sync syncs with untappd
func Sync(u *Untappd, cellar *pb.BeerCellar) int {
	drunk := u.GetRecentDrinks(u.f, u.c, cellar.SyncTime)

	if len(drunk) == 0 {
		return 0
	}

	for _, val := range drunk {
		rb := RemoveBeer(cellar, val)

		if rb != nil {
			rb.DrinkDate = time.Now().Unix()
			cellar.Drunk = append(cellar.Drunk, rb)
		} else {
			cellar.Drunk = append(cellar.Drunk, &pb.Beer{Id: val, Size: "Out"})
		}
	}

	cellar.SyncTime = time.Now().Unix()

	return len(drunk)
}

// RemoveBeer removes a beer from the cellar
func RemoveBeer(cellar *pb.BeerCellar, id int64) *pb.Beer {
	cellarIndex := -1
	cellarCost := -1
	for i, c := range cellar.GetCellars() {
		cost := GetRemoveCost(c, id)
		if cost >= 0 {
			if cellarIndex < 0 || cost < cellarCost {
				cellarIndex = i
				cellarCost = cost
			}
		}
	}

	if cellarIndex >= 0 {
		return Remove(cellar.Cellars[cellarIndex], id)
	}

	return nil
}

// CountBeers returns the number of beers of a given id in the cellar
func CountBeers(cellar *pb.BeerCellar, id int64) int {
	sum := 0
	for _, v := range cellar.GetCellars() {
		sum += CountBeersInCellar(v, id)
	}
	return sum
}

// GetNumberOfCellars gets the number of cellar boxes in the cellar
func GetNumberOfCellars(cellar *pb.BeerCellar) int {
	return len(cellar.Cellars)
}

// Size gets the size of the cellar
func Size(cellar *pb.BeerCellar) int {
	size := 0
	for _, v := range cellar.GetCellars() {
		size += len(v.GetBeers())
	}
	return size
}

// AddBuilt Adds a beer to the cellar
func AddBuilt(cellar *pb.BeerCellar, beer *pb.Beer) *pb.Cellar {
	bestCellar := -1
	bestScore := -1

	for i, v := range cellar.GetCellars() {
		insertCount := ComputeInsertCost(v, beer)

		if insertCount >= 0 && (insertCount < bestScore || bestScore < 0) {
			bestScore = insertCount
			bestCellar = i
		}
	}

	AddBeerToCellar(cellar.Cellars[bestCellar], beer)
	return cellar.Cellars[bestCellar]
}

// GetEmptyCellarCount Gets the number of empty cellars
func GetEmptyCellarCount(cellar *pb.BeerCellar) int {
	count := 0
	for _, v := range cellar.Cellars {
		if len(v.Beers) == 0 {
			count++
		}
	}
	return count
}

// NewBeerCellar creates new beer cellar
func NewBeerCellar(name string) *pb.BeerCellar {
	bc := &pb.BeerCellar{
		Name:     name,
		SyncTime: time.Now().Unix(),
		Cellars:  make([]*pb.Cellar, NumCellars),
	}

	for i := range bc.Cellars {
		bc.Cellars[i] = &pb.Cellar{
			Beers: make([]*pb.Beer, 0),
		}
	}

	bc.Drunk = make([]*pb.Beer, 0)

	return bc
}

// AddBeerByDays adds beers by days to the cellar.
func AddBeerByDays(cellar *pb.BeerCellar, id string, date string, size string, days string, count string) {
	startDate, _ := time.Parse("02/01/06", date)
	countVal, _ := strconv.Atoi(count)
	daysVal, _ := strconv.Atoi(days)
	for i := 0; i < countVal; i++ {
		AddBeer(cellar, id, startDate.Unix(), size)
		startDate = startDate.AddDate(0, 0, daysVal)
	}
}

// AddBeerByYears adds beers by years to the cellar.
func AddBeerByYears(cellar *pb.BeerCellar, id string, date string, size string, years string, count string) {
	startDate, _ := time.Parse("02/01/06", date)
	countVal, _ := strconv.Atoi(count)
	yearsVal, _ := strconv.Atoi(years)
	for i := 0; i < countVal; i++ {
		AddBeer(cellar, id, startDate.Unix(), size)
		startDate = startDate.AddDate(yearsVal, 0, 0)
	}
}

// AddBeer adds the beer to the cellar
func AddBeer(cellar *pb.BeerCellar, id string, date int64, size string) {
	idNum, _ := strconv.Atoi(id)
	if idNum >= 0 {
		AddBuilt(cellar, &pb.Beer{Id: int64(idNum), DrinkDate: date, Size: size})
	}
}

func (s *Server) recacheBeer(b *pb.Beer) {
	if b.Name == "" || strings.Contains(b.Name, "their API limit") || b.Abv == 0 {
		b.Name, b.Abv = s.ut.GetBeerDetails(b.Id)
		s.saveCellar()
	}
}

// Min returns the min of the parameters
func Min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

// ListBeers lists the cellared beers of a given type
func ListBeers(cellar *pb.BeerCellar, num int, btype string, date int64) []*pb.Beer {
	retList := MergeCellars(btype, cellar.Cellars...)

	pointer := -1
	for i, v := range retList {
		if i < num && v.DrinkDate < date {
			pointer = i
		}
	}

	return retList[:pointer+1]
}
