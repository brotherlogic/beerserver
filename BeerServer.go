package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
	pbgs "github.com/brotherlogic/goserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config    *pb.Config
	ut        *Untappd
	lastClean time.Time
}

const (
	//TOKEN Where we store the cellar
	TOKEN = "/github.com/brotherlogic/beerserver/cellar"

	//UTTOKEN where we store the untapped token
	UTTOKEN = "/github.com/brotherlogic/beerserver/untappd"
)

type mainFetcher struct{}

func (fetcher mainFetcher) Fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

// DoRegister Registers this server
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterBeerCellarServiceServer(server, s)
}

// ReportHealth Determines if the server is healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes this server
func (s *Server) Mote(master bool) error {
	if master {
		// Read the cellar
		bType := &pb.Config{}
		bResp, _, err := s.KSclient.Read(TOKEN, bType)

		if err != nil {
			return err
		}

		s.config = bResp.(*pb.Config)
		s.ut = GetUntappd(s.config.Token.Id, s.config.Token.Secret)
	}

	return nil
}

//GetState gets the state of the server
func (s *Server) GetState() []*pbgs.State {
	drunkDate := int64(0)
	lastDrunk := ""
	missing := int64(0)
	for _, drunk := range s.config.Drunk {
		if drunk.GetDrinkDate() > drunkDate {
			if drunk.CheckinId > 0 {
				drunkDate = drunk.GetDrinkDate()
				lastDrunk = fmt.Sprintf("%v", drunk.CheckinId)
			} else {
				missing++
			}
		}
	}
	return []*pbgs.State{
		&pbgs.State{Key: "lastddate", TimeValue: drunkDate},
		&pbgs.State{Key: "lastdrunk", Text: lastDrunk},
		&pbgs.State{Key: "lastsync", TimeValue: s.config.LastSync},
		&pbgs.State{Key: "last_clean", TimeValue: s.lastClean.Unix()},
		&pbgs.State{Key: "druuuunk", Value: int64(len(s.config.Drunk))},
		&pbgs.State{Key: "drnnnnnk", Value: missing},
	}
}

func (s *Server) checkSync(ctx context.Context) {
	if time.Now().Sub(time.Unix(s.config.LastSync, 0)) > time.Hour*24*7 {
		s.RaiseIssue(ctx, "BeerServer Sync Issue", fmt.Sprintf("Last Sync was %v", time.Unix(s.config.LastSync, 0)))
	}
}

func (s *Server) save() {
	s.KSclient.Save(TOKEN, s.config)
}

//GetUntappd builds a untappd retriever
func GetUntappd(id, secret string) *Untappd {
	return &Untappd{untappdID: id, untappdSecret: secret, u: mainUnmarshaller{}, f: mainFetcher{}, c: mainConverter{}}
}

//Init builds a server
func Init() *Server {
	s := &Server{&goserver.GoServer{}, &pb.Config{}, &Untappd{}, time.Unix(1, 0)}
	return s
}

func (s *Server) doSync(ctx context.Context) {
	s.syncDrunk(mainFetcher{})
}

func (s *Server) doMove(ctx context.Context) {
	s.moveToOnDeck(time.Now())
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var updateDrunk = flag.Bool("update", false, "Update the drunk")
	flag.Parse()

	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	server := Init()
	server.PrepServer()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)

	server.Register = server
	server.RegisterServer("beerserver", false)

	if *updateDrunk {
		server.Mote(true)
		server.loadDrunk("loaddata/brotherlogic.json")
		server.save()
		log.Fatalf("UPDATED: %v", server.config.Drunk)
	}

	server.RegisterRepeatingTask(server.doSync, time.Hour)
	server.RegisterRepeatingTask(server.doMove, time.Hour)
	server.RegisterRepeatingTask(server.clearDeck, time.Minute*5)
	server.RegisterRepeatingTask(server.checkSync, time.Hour)

	server.Serve()
}
