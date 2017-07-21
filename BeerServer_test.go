package main

import (
	"os"

	"github.com/brotherlogic/keystore/client"
)

func NewTestBeerServer(dir string, delete bool) Server {
	s := Init()
	if delete {
		os.RemoveAll(dir)
	}
	s.KSclient = *keystoreclient.GetTestClient(dir)

	return s
}
