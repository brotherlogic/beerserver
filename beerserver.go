package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"
	pbgs "github.com/brotherlogic/goserver/proto"
	pbp "github.com/brotherlogic/printer/proto"
	rpb "github.com/brotherlogic/reminders/proto"
)

var (
	pstashSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "beerserver_stash",
		Help: "The size of the beer stash",
	})
	psmallSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "beerserver_small",
		Help: "The size of the small beers",
	})
	pbomberSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "beerserver_bombers",
		Help: "The size of the bomber beers",
	})
	pdrunk = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "beerserver_drunk",
		Help: "The number of beers drunk",
	})
)

func (s *Server) monitor(ctx context.Context, config *pb.Config) error {
	for _, c := range config.GetCellar().GetSlots() {
		if c.Accepts == "small" {
			psmallSize.Set(float64(len(c.Beers)))
		}
		if c.Accepts == "bomber" {
			pbomberSize.Set(float64(len(c.Beers)))
		}
		if c.Accepts == "stash" {
			pstashSize.Set(float64(len(c.Beers)))
		}
	}
	return nil
}

type printer interface {
	print(ctx context.Context, lines []string) error
}

type prodPrinter struct {
	testing bool
	fail    bool
	count   int64
	dial    func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodPrinter) print(ctx context.Context, lines []string) error {
	if p.fail {
		return fmt.Errorf("Built fail")
	}
	if p.testing {
		return nil
	}
	conn, err := p.dial(ctx, "printer")
	if err != nil {
		return err
	}
	defer conn.Close()

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
	ut        *Untappd
	lastClean time.Time
	printer   printer
	lastErr   string
}

//Init builds a server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&Untappd{},
		time.Unix(1, 0),
		&prodPrinter{},
		"",
	}
	s.printer = &prodPrinter{dial: s.FDialServer}
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
	rpb.RegisterReminderReceiverServer(server, s)
}

// ReportHealth Determines if the server is healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (s *Server) validateCellars(ctx context.Context, config *pb.Config) error {
	for _, c := range config.GetCellar().GetSlots() {
		for _, b := range c.Beers {
			if b.Uid == 0 {
				b.Uid = time.Now().UnixNano()
				s.Log(fmt.Sprintf("Updating UID for %v", b))
				time.Sleep(time.Millisecond * 10)
			}
		}
	}

	for _, b := range config.GetCellar().GetOnDeck() {
		if b.Uid == 0 {
			b.Uid = time.Now().UnixNano()
			time.Sleep(time.Millisecond * 10)
		}
	}

	//Create the cellars if we need to
	if config.GetCellar() == nil {
		config.Cellar = &pb.Cellar{Slots: []*pb.CellarSlot{}, OnDeck: []*pb.Beer{}}
	}
	if len(config.GetCellar().GetSlots()) == 0 {
		config.GetCellar().Slots = append(config.GetCellar().Slots, &pb.CellarSlot{Accepts: "small", NumSlots: 30})
		config.GetCellar().Slots = append(config.GetCellar().Slots, &pb.CellarSlot{Accepts: "bomber", NumSlots: 20})
		config.GetCellar().Slots = append(config.GetCellar().Slots, &pb.CellarSlot{Accepts: "stash", NumSlots: 9999})
		config.GetCellar().Slots = append(config.GetCellar().Slots, &pb.CellarSlot{Accepts: "wedding", NumSlots: 24})
	}

	return s.save(ctx, config)
}

// Mote promotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

func (s *Server) load(ctx context.Context) (*pb.Config, error) {
	bType := &pb.Config{}
	bResp, _, err := s.KSclient.Read(ctx, TOKEN, bType)

	if err != nil {
		return nil, err
	}

	config := bResp.(*pb.Config)
	s.ut = GetUntappd(config.GetToken().GetId(), config.GetToken().GetSecret())
	s.ut.l = s.Log

	pdrunk.Set(float64(len(config.GetDrunk())))

	ndrunk := make([]*pb.Beer, 0)
	seen := make(map[int64]bool)
	for _, b := range config.GetDrunk() {
		if !seen[b.GetDrinkDate()] {
			seen[b.GetDrinkDate()] = true
			ndrunk = append(ndrunk, b)
		}
	}

	if len(ndrunk) != len(config.GetDrunk()) {
		s.RaiseIssue("Cleaning drunk list!", fmt.Sprintf("Removing %v beers from the drunk list, there are %v remaining", len(config.GetDrunk())-len(ndrunk), len(ndrunk)))
		if len(ndrunk) > 0 {
			config.Drunk = ndrunk
		}
	}

	return config, nil
}

//GetState gets the state of the server
func (s *Server) GetState() []*pbgs.State {
	return []*pbgs.State{}
}

func (s *Server) loadDrunk(filestr string, config *pb.Config) {
	data, _ := ioutil.ReadFile(filestr)
	b, _ := s.ut.convertDrinkListToBeers(string(data), mainUnmarshaller{})
	for _, beer := range b {
		found := false
		for _, d := range config.Drunk {
			if d.CheckinId == beer.CheckinId {
				found = true
			}
		}

		if !found {
			config.Drunk = append(config.Drunk, beer)
		}
	}
}

func (s *Server) checkSync(ctx context.Context, config *pb.Config) error {
	if time.Now().Sub(time.Unix(config.LastSync, 0)) > time.Hour*24*7 {
		s.RaiseIssue("BeerServer Sync Issue", fmt.Sprintf("Last Sync was %v (%v)", time.Unix(config.LastSync, 0), s.lastErr))
	}

	drunkDate := int64(0)
	for _, drunk := range config.Drunk {
		if drunk.GetDrinkDate() > drunkDate {
			if drunk.CheckinId > 0 {
				drunkDate = drunk.GetDrinkDate()
			}
		}
	}

	if time.Now().Sub(time.Unix(drunkDate, 0)) > time.Hour*24*7 {
		s.RaiseIssue("BeerServer Sync Issue", fmt.Sprintf("Last Syncd beer was %v", time.Unix(drunkDate, 0)))
	}

	return nil
}

func (s *Server) save(ctx context.Context, config *pb.Config) error {
	return s.KSclient.Save(ctx, TOKEN, config)
}

//GetUntappd builds a untappd retriever
func GetUntappd(id, secret string) *Untappd {
	return &Untappd{untappdID: id, untappdSecret: secret, u: mainUnmarshaller{}, f: mainFetcher{}, c: mainConverter{}}
}

func (s *Server) refreshBreweryID(ctx context.Context, config *pb.Config) error {
	for _, cellar := range config.GetCellar().GetSlots() {
		for _, b := range cellar.GetBeers() {
			if b.GetBreweryId() == 0 {
				tb := s.ut.GetBeerDetails(b.Id)
				if tb.GetBreweryId() > 0 {
					b.BreweryId = tb.GetBreweryId()
					s.Log(fmt.Sprintf("%v -> %v has no brewery id", tb.GetBreweryId(), b))
					s.RaiseIssue("Brewery ID is missing", fmt.Sprintf("%v is missing the brewery ID", b))
				}
			}
		}
	}
	return s.save(ctx, config)
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	server := Init()
	server.PrepServer()

	server.Register = server
	err := server.RegisterServerV2("beerserver", false, true)
	if err != nil {
		return
	}

	server.Serve()
}
