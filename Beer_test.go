package main

import "testing"

func TestNewBeer(t *testing.T) {
	beer, _ := NewBeer("123~01/02/16~bomber")

	if beer.Id != 123 {
		t.Errorf("Beer id %d is not 123\n", beer.Id)
	}

	if beer.DrinkDate != 1454284800 {
		t.Errorf("Date %q is not 01/02/16\n", beer.DrinkDate)
	}
}

func TestBadBeerSize(t *testing.T) {
	_, err := NewBeer("123~01/01/16~madeupsize")

	if err == nil {
		t.Errorf("Beer size has not induced failure\n")
	}
}

func TestBadDate(t *testing.T) {
	_, err := NewBeer("123~01/15/16~bomber")

	if err == nil {
		t.Errorf("Beer with bad date has been parsed correctly\n")
	}
}
