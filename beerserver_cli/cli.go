package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/beerserver/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	host, port, err := utils.Resolve("beerserver")
	if err != nil {
		log.Fatalf("Unable to reach beerserver: %v", err)
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}

	client := pb.NewBeerCellarServiceClient(conn)
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Minute)
	defer cancel2()

	switch os.Args[1] {
	case "list":
		listFlags := flag.NewFlagSet("List", flag.ExitOnError)
		var ondeck = listFlags.Bool("deck", false, "View on deck")
		var cellar = listFlags.Int("cellar", 1, "Cellar to view")
		if err := listFlags.Parse(os.Args[2:]); err == nil {
			ctx, cancel := utils.BuildContext("beerserver-cli", "beerserver-cli")
			defer cancel()
			list, err := client.ListBeers(ctx, &pb.ListBeerRequest{OnDeck: *ondeck})
			if err == nil {
				for i := 0; i < len(list.Beers); i++ {
					if *ondeck {
						fmt.Printf("%v. %v [%v] (%v)\n", i, list.Beers[i].Name, list.Beers[i].Id, time.Unix(list.Beers[i].DrinkDate, 0))
					}
					for _, b := range list.Beers {
						if int(b.Index) == i && int(b.InCellar) == *cellar {
							fmt.Printf("%v. %v (%v) - %v\n", i, b.Name, time.Unix(b.DrinkDate, 0), b.Id)
						}
					}
				}
			} else {
				fmt.Printf("Error getting beers: %v\n", err)
			}
		}
	case "add":
		addFlags := flag.NewFlagSet("Add", flag.ExitOnError)
		var id = addFlags.Int("id", -1, "Id of beer to add")
		var quantity = addFlags.Int("quantity", 0, "Number of beers to add")
		var size = addFlags.String("size", "", "Size of beer")
		if err := addFlags.Parse(os.Args[2:]); err == nil {
			if len(*size) != 0 {
				_, err := client.AddBeer(ctx2, &pb.AddBeerRequest{Beer: &pb.Beer{Id: int64(*id), Size: *size}, Quantity: int32(*quantity)})
				if err != nil {
					fmt.Printf("Error adding beer: %v\n", err)
				}
			}
		}
	}

}
