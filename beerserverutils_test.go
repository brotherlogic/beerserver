package main

import (
	"context"
	"testing"
	"time"

	pb "github.com/brotherlogic/beerserver/proto"
)

func TestRemoveDrunkThings(t *testing.T) {
	s := InitTestServer(".testremovedrunkthings", true)
	ts := time.Now().Unix()
	s.config.Cellar.OnDeck = append(s.config.Cellar.OnDeck, &pb.Beer{Id: 123, DrinkDate: ts - 2})
	s.config.Drunk = append(s.config.Drunk, &pb.Beer{Id: 123, DrinkDate: ts})

	s.clearDeck(context.Background())

	if len(s.config.Cellar.OnDeck) == 1 {
		t.Errorf("Failure to clear the decks: %v", s.config)
	}
}

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
	s.reorderCellar(context.Background(), slot)
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
	s.reorderCellar(context.Background(), slot)
	if s.cellarOutOfOrder(context.Background(), slot) {
		t.Errorf("Problem in ordering")
	}

	if slot.Beers[0].Index != 0 || slot.Beers[0].Id != 124 {
		t.Errorf("Cellar is not ordered right still: %v", slot)
	}

}
