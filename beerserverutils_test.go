package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestCellarOrder(t *testing.T) {
	s := InitTestServer(".testtry", true)
	slot := &pb.CellarSlot{Beers: []*pb.Beer{
		&pb.Beer{Id: 123, DrinkDate: 5, Index: 1},
		&pb.Beer{Id: 123, DrinkDate: 10, Index: 2},
	}}
	if s.cellarOutOfOrder(context.Background(), slot) {
		t.Errorf("Problem in ordering")
	}

}

func TestCellarOrderSort(t *testing.T) {
	s := InitTestServer(".testtry", true)
	slot := &pb.CellarSlot{Beers: []*pb.Beer{
		&pb.Beer{Id: 123, DrinkDate: 5, Index: 1},
		&pb.Beer{Id: 123, DrinkDate: 10, Index: 2},
	}}
	s.reorderCellar(context.Background(), slot, int32(0))
	if s.cellarOutOfOrder(context.Background(), slot) {
		t.Errorf("Problem in ordering")
	}

}

func TestCellarOrderBad(t *testing.T) {
	s := InitTestServer(".testtry", true)
	slot := &pb.CellarSlot{Beers: []*pb.Beer{
		&pb.Beer{Id: 123, DrinkDate: 10, Index: 1},
		&pb.Beer{Id: 124, DrinkDate: 5, Index: 2},
	}}
	if !s.cellarOutOfOrder(context.Background(), slot) {
		t.Errorf("Problem in ordering")
	}

}

func TestCellarOrderBadSort(t *testing.T) {
	s := InitTestServer(".testtry", true)
	slot := &pb.CellarSlot{Beers: []*pb.Beer{
		&pb.Beer{Id: 123, DrinkDate: 10, Index: 10},
		&pb.Beer{Id: 124, DrinkDate: 5, Index: 25},
	}}
	s.reorderCellar(context.Background(), slot, int32(0))
	if s.cellarOutOfOrder(context.Background(), slot) {
		t.Errorf("Problem in ordering")
	}

	if slot.Beers[0].Index != 0 || slot.Beers[0].Id != 124 {
		t.Errorf("Cellar is not ordered right still: %v", slot)
	}

}

func TestAddDrinks(t *testing.T) {
	s := InitTestServer(".testadddrunks", true)
	s.addDrunks(context.Background(), &pb.Config{}, []*pb.Beer{})
}
