package main

import (
	"log"
	"strconv"
	"strings"

	pb "github.com/brotherlogic/beerserver/proto"
)

import "time"

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
func Sync(u *Untappd, cellar *pb.BeerCellar, fetcher httpResponseFetcher, converter responseConverter) {
	drunk := u.GetRecentDrinks(fetcher, converter, cellar.SyncTime)
	log.Printf("Found these: %v\n", drunk)
	for _, val := range drunk {
		log.Printf("Removing %v from cellar\n", val)
		RemoveBeer(cellar, val)
	}

	cellar.SyncTime = time.Now().Unix()
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
		log.Printf("Removing %v from %v\n", id, cellarIndex)
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

		log.Printf("Adding beer to cellar %v: %v\n", i, insertCount)

		if insertCount >= 0 && (insertCount < bestScore || bestScore < 0) {
			bestScore = insertCount
			bestCellar = i
		}
	}

	log.Printf("Adding to %v given %v", bestCellar, cellar.Cellars)
	AddBeerToCellar(cellar.Cellars[bestCellar], beer)
	log.Printf("DONE Adding to %v given %v", bestCellar, cellar.Cellars)
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
		Cellars:  make([]*pb.Cellar, 8),
	}

	for i := range bc.Cellars {
		bc.Cellars[i] = &pb.Cellar{
			Beers: make([]*pb.Beer, 0),
		}
	}

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
	if b.Name == "" || strings.Contains(b.Name, "their API limit") {
		b.Name = s.ut.GetBeerName(b.Id)
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
	log.Printf("Cellar looks like %v\n", cellar.Cellars)
	retList := MergeCellars(btype, cellar.Cellars...)

	pointer := -1
	for i, v := range retList {
		if i < num && v.DrinkDate < date {
			pointer = i
		} else {
			log.Printf("%v, %v and %v %v\n", i, num, v.DrinkDate, date)
		}
	}

	log.Printf("RETURNING %v", retList)
	return retList[:pointer+1]
}
