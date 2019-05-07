package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
	pbgs "github.com/brotherlogic/goserver/proto"
	pbp "github.com/brotherlogic/printer/proto"
)

type printer interface {
	print(ctx context.Context, lines []string) error
}

type prodPrinter struct {
	testing bool
	count   int64
}

func (p *prodPrinter) print(ctx context.Context, lines []string) error {
	if p.testing {
		return nil
	}
	host, port, err := utils.Resolve("printer")
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := pbp.NewPrintServiceClient(conn)
	_, err = client.Print(ctx, &pbp.PrintRequest{Lines: lines, Origin: "beerserver"})

	if err == nil {
		p.count++
	}

	return err

}

//Server main server type
type Server struct {
	*goserver.GoServer
	config    *pb.Config
	ut        *Untappd
	lastClean time.Time
	printer   printer
}

//Init builds a server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&pb.Config{},
		&Untappd{},
		time.Unix(1, 0),
		&prodPrinter{},
	}
	return s
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

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (s *Server) validateCellars(ctx context.Context) {
	for _, c := range s.config.Cellar.Slots {
		for _, b := range c.Beers {
			if b.Uid == 0 {
				b.Uid = time.Now().UnixNano()
				s.Log(fmt.Sprintf("Updating UID for %v", b))
				time.Sleep(time.Millisecond * 10)
			}
		}
	}
}

// Mote promotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		// Read the cellar
		bType := &pb.Config{}
		bResp, _, err := s.KSclient.Read(ctx, TOKEN, bType)

		if err != nil {
			return err
		}

		s.config = bResp.(*pb.Config)
		s.ut = GetUntappd(s.config.Token.Id, s.config.Token.Secret)
		s.Log(fmt.Sprintf("FOUND %v and %v", s.config.Token.Id, s.config.Token.Secret))
		s.ut.l = s.Log
		s.validateCellars(ctx)
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
	missingUID := int64(0)
	if s.config != nil {
		for _, c := range s.config.Cellar.Slots {
			for _, b := range c.Beers {
				if b.Uid == 0 {
					missingUID++
				}
			}
		}
	}

	return []*pbgs.State{
		&pbgs.State{Key: "missing_uid", Value: missingUID},
		&pbgs.State{Key: "prints", Value: s.printer.(*prodPrinter).count},
		&pbgs.State{Key: "lastddate", TimeValue: drunkDate},
		&pbgs.State{Key: "lastdrunk", Text: lastDrunk},
		&pbgs.State{Key: "lastsync", TimeValue: s.config.LastSync},
		&pbgs.State{Key: "last_clean", TimeValue: s.lastClean.Unix()},
		&pbgs.State{Key: "druuuunk", Value: int64(len(s.config.Drunk))},
		&pbgs.State{Key: "drnnnnnk", Value: missing},
	}
}

func (s *Server) loadDrunk(filestr string) {
	data, _ := ioutil.ReadFile(filestr)
	b, _ := s.ut.convertDrinkListToBeers(string(data), mainUnmarshaller{})
	for _, beer := range b {
		found := false
		for _, d := range s.config.Drunk {
			if d.CheckinId == beer.CheckinId {
				found = true
			}
		}

		if !found {
			s.config.Drunk = append(s.config.Drunk, beer)
		}
	}
}

func (s *Server) checkSync(ctx context.Context) error {
	if time.Now().Sub(time.Unix(s.config.LastSync, 0)) > time.Hour*24*7 {
		s.RaiseIssue(ctx, "BeerServer Sync Issue", fmt.Sprintf("Last Sync was %v", time.Unix(s.config.LastSync, 0)), false)
	}

	drunkDate := int64(0)
	for _, drunk := range s.config.Drunk {
		if drunk.GetDrinkDate() > drunkDate {
			if drunk.CheckinId > 0 {
				drunkDate = drunk.GetDrinkDate()
			}
		}
	}

	if time.Now().Sub(time.Unix(drunkDate, 0)) > time.Hour*24*7 {
		s.RaiseIssue(ctx, "BeerServer Sync Issue", fmt.Sprintf("Last Syncd beer was %v", time.Unix(drunkDate, 0)), false)
	}

	return nil
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, TOKEN, s.config)
}

//GetUntappd builds a untappd retriever
func GetUntappd(id, secret string) *Untappd {
	return &Untappd{untappdID: id, untappdSecret: secret, u: mainUnmarshaller{}, f: mainFetcher{}, c: mainConverter{}}
}

func (s *Server) doSync(ctx context.Context) error {
	s.syncDrunk(ctx, mainFetcher{})
	return nil
}

func (s *Server) doMove(ctx context.Context) error {
	s.moveToOnDeck(ctx, time.Now())
	return nil
}

func (s *Server) checkCellars(ctx context.Context) error {
	for i, slot := range s.config.Cellar.Slots {
		if s.cellarOutOfOrder(ctx, slot) {
			s.RaiseIssue(ctx, "Cellar not ordered", fmt.Sprintf("Cellar %v is not ordered correctly", i), false)
		}
	}
	return nil
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
		server.Mote(context.Background(), true)
		server.loadDrunk("loaddata/brotherlogic.json")
		server.save(context.Background())
		log.Fatalf("UPDATED: %v", server.config.Drunk)
	}

	server.RegisterRepeatingTask(server.checkCellars, "check_cellars", time.Minute)
	server.RegisterRepeatingTask(server.doSync, "do_sync", time.Hour)
	server.RegisterRepeatingTask(server.doMove, "do_move", time.Hour)
	server.RegisterRepeatingTask(server.clearDeck, "clear_deck", time.Minute*5)
	server.RegisterRepeatingTask(server.checkSync, "check_sync", time.Hour)

	server.Serve()
}
