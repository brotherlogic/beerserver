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
