package main

import "errors"
import "strconv"
import "strings"
import "time"

import pb "github.com/brotherlogic/beer/proto"

// ByDate used to sort beer
type ByDate []*pb.Beer

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].GetDrinkDate() < a[j].GetDrinkDate() }

// NewBeer Builds a new beer
func NewBeer(line string) (pb.Beer, error) {
	elems := strings.Split(line, "~")

	if len(elems) != 3 {
		return pb.Beer{}, errors.New("Line is misspecified: " + line)
	}

	bid, _ := strconv.Atoi(elems[0])
	size := elems[2]

	// Ensure that the date parses correctly
	t, err := time.Parse("02/01/06", elems[1])

	// Ensure that the beer size is set
	if err == nil && size != "bomber" && size != "small" {
		err = errors.New(size + " is not a valid beer size")
	}

	beer := pb.Beer{Id: int64(bid), DrinkDate: t.Unix(), Size: size}
	return beer, err
}
