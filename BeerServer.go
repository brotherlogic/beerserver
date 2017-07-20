package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	cellar    *pb.BeerCellar
	nameCache map[int64]string
	ut        *Untappd
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

func (s *Server) cacheName(id int64, name string) {
	s.nameCache[id] = name
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
	return nil
}

func (s *Server) saveCellar() {
	t := time.Now()
	log.Printf("Writing cellar")
	s.KSclient.Save(TOKEN, s.cellar)
	s.LogFunction("saveCellar", int32(time.Now().Sub(t).Nanoseconds()/1000000))
}

//AddBeer adds a beer to the cellar
func (s *Server) AddBeer(ctx context.Context, beer *pb.Beer) (*pb.Cellar, error) {
	log.Printf("CELLAR HERE = %v", s.cellar)
	cel := AddBuilt(s.cellar, beer)
	s.saveCellar()
	return cel, nil
}

//GetUntappd builds a untappd retriever
func GetUntappd(id, secret string) *Untappd {
	return &Untappd{untappdID: id, untappdSecret: secret, u: mainUnmarshaller{}, f: mainFetcher{}, c: mainConverter{}}
}

//Init builds a server
func Init() Server {
	s := Server{&goserver.GoServer{}, &pb.BeerCellar{}, make(map[int64]string), &Untappd{}}
	return s
}

func main() {
	var id = flag.String("id", "", "Untappd id")
	var secret = flag.String("secret", "", "Untappd secret")
	var quiet = flag.Bool("quiet", true, "Show all output")
	flag.Parse()

	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	server := Init()
	server.PrepServer()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)

	if len(*id) > 0 {
		server.KSclient.Save(UTTOKEN, &pb.Token{Id: *id, Secret: *secret})
		os.Exit(1)
	}

	tType := &pb.Token{}
	tResp, err := server.KSclient.Read(UTTOKEN, tType)

	if err != nil {
		log.Fatalf("Unable to read token: %v", err)
	}

	sToken := tResp.(*pb.Token)
	server.ut = GetUntappd(sToken.Id, sToken.Secret)

	bType := &pb.BeerCellar{}
	bResp, err := server.KSclient.Read(TOKEN, bType)

	if err != nil {
		log.Fatalf("Unable to read token: %v", err)
	}

	nCellar := bResp.(*pb.BeerCellar)
	if nCellar == nil || nCellar.Cellars == nil {
		nCellar = NewBeerCellar("mycellar")
	}
	server.cellar = nCellar

	log.Printf("INIT CELLAR 1 %v", server.cellar)
	server.Register = &server
	server.RegisterServer("beerserver", false)
	server.Serve()
}
