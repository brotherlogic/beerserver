package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/beerserver/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.ManualContext("buildserver-"+os.Args[1], time.Second*10)
	defer cancel()
	conn, err := utils.LFDialServer(ctx, "beerserver")

	defer conn.Close()

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}

	client := pb.NewBeerCellarServiceClient(conn)

	switch os.Args[1] {
	case "update":
		_, err := client.Update(ctx, &pb.UpdateRequest{})
		fmt.Printf("Updated: %v\n", err)
	case "consolidate":
		ctx, cancel := utils.BuildContext("beerserver-cli", "beerserver-cli")
		defer cancel()
		_, err := client.Consolidate(ctx, &pb.ConsolidateRequest{})
		fmt.Printf("Consolidated: %v\n", err)
	case "list":
		listFlags := flag.NewFlagSet("List", flag.ExitOnError)
		var ondeck = listFlags.Bool("deck", false, "View on deck")
		var cellar = listFlags.Int("cellar", 1, "Cellar to view")
		var summarise = listFlags.Bool("sum", false, "Summarize the view")
		if err := listFlags.Parse(os.Args[2:]); err == nil {
			list, err := client.ListBeers(ctx, &pb.ListBeerRequest{OnDeck: *ondeck})
			if err == nil {
				if *summarise {
					idcount := make(map[int64]int)
					for _, beer := range list.Beers {
						if beer.InCellar == int32(*cellar) {
							idcount[beer.GetId()]++
						}
					}
					for id, count := range idcount {
						for _, beer := range list.Beers {
							if beer.GetId() == id {
								fmt.Printf("%v - %v\n", count, beer.GetName())
								break
							}
						}
					}
				} else {
					for i := 0; i < len(list.Beers); i++ {
						if *ondeck {
							fmt.Printf("%v. %v [%v] (%v) - %v [%v]\n", i, list.Beers[i].Name, list.Beers[i].Id, time.Unix(list.Beers[i].DrinkDate, 0), list.Beers[i].Uid, list.Beers[i].Size)
						}
						for _, b := range list.Beers {
							if int(b.Index) == i && int(b.InCellar) == *cellar {
								fmt.Printf("%v. %v/%v (%v) - %v [%v]\n", b.Order, b.Name, b.BreweryId, time.Unix(b.DrinkDate, 0), b.Uid, b.Id)
							}
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
				_, err := client.AddBeer(ctx, &pb.AddBeerRequest{Beer: &pb.Beer{Id: int64(*id), Size: *size}, Quantity: int32(*quantity)})
				if err != nil {
					fmt.Printf("Error adding beer: %v\n", err)
				}
			}
		}
	case "delete":
		deleteFlags := flag.NewFlagSet("Delete", flag.ExitOnError)
		var id = deleteFlags.Int("id", -1, "Id of beer to delete")
		if err := deleteFlags.Parse(os.Args[2:]); err == nil {
			if *id != -1 {
				_, err := client.DeleteBeer(ctx, &pb.DeleteBeerRequest{Uid: int64(*id)})
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	}

}
