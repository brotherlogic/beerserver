package main

import (
	"context"
	"log"
	"os"

	pb "github.com/brotherlogic/beerserver/proto"
	keystoreclient "github.com/brotherlogic/keystore/client"
)

func doLog(str string) {
	log.Printf("STR " + str)
}

func InitTestServer(dir string, delete bool) *Server {
	s := Init()

	if delete {
		os.RemoveAll(dir)
	}
	s.KSclient = *keystoreclient.GetTestClient(dir)
	s.SkipLog = true
	s.SkipIssue = true
	s.printer = &prodPrinter{testing: true}
	s.ut = GetTestUntappd()
	//s.ut.l = doLog

	s.validateCellars(context.Background(), &pb.Config{Token: &pb.Token{}})

	return s
}
