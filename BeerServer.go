package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
	pbgs "github.com/brotherlogic/goserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config *pb.Config
	ut     *Untappd
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
	}

	return nil
}

//GetState gets the state of the server
func (s *Server) GetState() []*pbgs.State {
	return []*pbgs.State{}
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
	s := &Server{&goserver.GoServer{}, &pb.Config{}, &Untappd{}}
	return s
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var init = flag.Bool("init", false, "Init the output")
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

	if *init {
		server.config.Cellar.Slots = append(server.config.Cellar.Slots, &pb.CellarSlot{Accepts: "small", NumSlots: 30})
		server.config.Cellar.Slots = append(server.config.Cellar.Slots, &pb.CellarSlot{Accepts: "bomber", NumSlots: 20})
		server.config.Cellar.Slots = append(server.config.Cellar.Slots, &pb.CellarSlot{Accepts: "bomber", NumSlots: 20})
		server.config.Cellar.Slots = append(server.config.Cellar.Slots, &pb.CellarSlot{Accepts: "bomber", NumSlots: 20})
		server.loadDrunk("loaddata/brotherlogic.json")
		server.save()
		panic("Saved!")
	}

	server.Serve()
}
