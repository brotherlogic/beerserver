package main

import (
	"context"
	"log"
	"os"

	pb "github.com/brotherlogic/beerserver/proto"
	"github.com/brotherlogic/keystore/client"
)

func doLog(str string) {
	log.Printf(str)
}

func InitTestServer(dir string, delete bool) *Server {
	s := Init()

	if delete {
		os.RemoveAll(dir)
	}
	s.KSclient = *keystoreclient.GetTestClient(dir)
	s.GoServer.KSclient.Save(context.Background(), TOKEN, &pb.Config{Token: &pb.Token{}})
	s.SkipLog = true
	s.SkipIssue = true
	s.printer = &prodPrinter{testing: true}
	s.ut = GetTestUntappd()
	s.ut.l = doLog

	return s
}
