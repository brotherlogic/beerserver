package main

import (
	"log"
	"os"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/beerserver/proto"
)

func doLog(str string) {
	log.Printf(str)
}

func InitTestServer(dir string, delete bool) *Server {
	s := Init()
	s.config.Cellar = &pb.Cellar{Slots: []*pb.CellarSlot{&pb.CellarSlot{Accepts: "bomber", NumSlots: 10}, &pb.CellarSlot{Accepts: "small", NumSlots: 10}}}
	if delete {
		os.RemoveAll(dir)
	}
	s.KSclient = *keystoreclient.GetTestClient(dir)
	s.SkipLog = true
	s.ut = GetTestUntappd()
	s.ut.l = doLog

	return s
}
