package main

import (
	"os"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/beerserver/proto"
)

func InitTestServer(dir string, delete bool) *Server {
	s := Init()
	s.config.Cellar = &pb.Cellar{Slots: []*pb.CellarSlot{&pb.CellarSlot{Accepts: "bomber", NumSlots: 10}}}
	if delete {
		os.RemoveAll(dir)
	}
	s.KSclient = *keystoreclient.GetTestClient(dir)
	s.ut = GetTestUntappd()

	return s
}
