package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	pb "github.com/brotherlogic/beerserver/proto"
)

//Untappd holds the untappd details
type Untappd struct {
	untappdID     string
	untappdSecret string
	u             unmarshaller
	f             httpResponseFetcher
	c             responseConverter
}

var beerMap map[int]string

// Match a match when doing a search
type Match struct {
	id   int
	name string
}

type unmarshaller interface {
	Unmarshal([]byte, *map[string]interface{}) error
}
type mainUnmarshaller struct{}

func (unmarshaller mainUnmarshaller) Unmarshal(inp []byte, resp *map[string]interface{}) error {
	return json.Unmarshal(inp, resp)
}

type responseConverter interface {
	Convert(*http.Response) ([]byte, error)
}

type mainConverter struct{}

func (converter mainConverter) Convert(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

type httpResponseFetcher interface {
	Fetch(url string) (*http.Response, error)
}

// Search finds beers in the cache that match blah
func Search(namePart string) []Match {
	ret := make([]Match, 0, len(beerMap))
	for k, v := range beerMap {
		if strings.Contains(v, namePart) {
			ret = append(ret, Match{id: k, name: v})
		}
	}

	return ret
}

func cacheBeer(id int, name string) {
	beerMap[id] = name
}

func init() {
	beerMap = make(map[int]string)
}

func (u *Untappd) getBeerPage(fetcher httpResponseFetcher, converter responseConverter, id int) string {
	url := "https://api.untappd.com/v4/beer/info/BID?client_id=CLIENTID&client_secret=CLIENTSECRET&compact=true"
	url = strings.Replace(url, "BID", strconv.Itoa(id), 1)
	url = strings.Replace(url, "CLIENTID", u.untappdID, 1)
	url = strings.Replace(url, "CLIENTSECRET", u.untappdSecret, 1)

	response, err := fetcher.Fetch(url)

	if err == nil {
		contents, err := converter.Convert(response)
		if err == nil {
			return string(contents)
		}
	}

	return "Failed to retrieve " + strconv.Itoa(id)
}

func (u *Untappd) getVenuePage(fetcher httpResponseFetcher, converter responseConverter, id int) string {
	url := "https://api.untappd.com/v4/venue/info/VID?client_id=CLIENTID&client_secret=CLIENTSECRET"
	url = strings.Replace(url, "VID", strconv.Itoa(id), 1)
	url = strings.Replace(url, "CLIENTID", u.untappdID, 1)
	url = strings.Replace(url, "CLIENTSECRET", u.untappdSecret, 1)

	response, err := fetcher.Fetch(url)

	if err == nil {
		contents, err := converter.Convert(response)
		if err == nil {
			return string(contents)
		}
	}

	return "Failed to retrieve " + strconv.Itoa(id)
}

func (u *Untappd) getUserPage(fetcher httpResponseFetcher, converter responseConverter, username string, minID int) (string, error) {
	url := "https://api.untappd.com/v4/user/checkins/USERNAME?client_id=CLIENTID&client_secret=CLIENTSECRET&min_id=MINID"
	url = strings.Replace(url, "USERNAME", username, 1)
	url = strings.Replace(url, "CLIENTID", u.untappdID, 1)
	url = strings.Replace(url, "CLIENTSECRET", u.untappdSecret, 1)
	url = strings.Replace(url, "MINID", strconv.Itoa(minID), 1)

	response, err := fetcher.Fetch(url)
	if err != nil {
		return "", err
	}

	contents, _ := converter.Convert(response)
	return string(contents), nil
}

func convertPageToName(page string, unmarshaller unmarshaller) string {
	var mapper map[string]interface{}
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		return "Failed to unmarshal"
	}

	meta := mapper["meta"].(map[string]interface{})
	metaCode := int(meta["code"].(float64))
	if metaCode != 200 {
		return meta["error_detail"].(string)
	}

	response := mapper["response"].(map[string]interface{})
	beer := response["beer"].(map[string]interface{})
	brewery := beer["brewery"].(map[string]interface{})
	return brewery["brewery_name"].(string) + " - " + beer["beer_name"].(string)
}

func convertPageToABV(page string, unmarshaller unmarshaller) float32 {
	var mapper map[string]interface{}
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		return -1
	}

	meta := mapper["meta"].(map[string]interface{})
	metaCode := int(meta["code"].(float64))
	if metaCode != 200 {
		return -1
	}

	response := mapper["response"].(map[string]interface{})
	beer := response["beer"].(map[string]interface{})
	abv := beer["beer_abv"].(float64)

	return float32(abv)
}

func (u *Untappd) convertPageToDrinks(page string, unmarshaller unmarshaller) ([]*pb.Beer, error) {
	var mapper map[string]interface{}
	var values []*pb.Beer
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		return values, err
	}

	meta := mapper["meta"].(map[string]interface{})
	metaCode := int(meta["code"].(float64))
	if metaCode != 200 {
		return values, errors.New("Couldn't retrieve drinks")
	}

	response := mapper["response"].(map[string]interface{})
	venue := response["venue"].(map[string]interface{})
	checkins := venue["checkins"].(map[string]interface{})
	items := checkins["items"].([]interface{})

	for _, v := range items {
		beer := v.(map[string]interface{})["beer"].(map[string]interface{})
		beerID := int64(beer["bid"].(float64))
		date := string(v.(map[string]interface{})["created_at"].(string))
		cdate, _ := time.Parse(time.RFC1123Z, date)
		nbeer := &pb.Beer{Id: beerID, DrinkDate: cdate.Unix()}
		values = append(values, nbeer)
	}

	return values, nil
}

func (u *Untappd) convertUserPageToDrinks(page string, unmarshaller unmarshaller) ([]*pb.Beer, error) {
	var mapper map[string]interface{}
	var values []*pb.Beer
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		return values, err
	}

	meta := mapper["meta"].(map[string]interface{})
	metaCode := int(meta["code"].(float64))
	if metaCode != 200 {
		return values, errors.New("Couldn't retrieve drinks")
	}

	response := mapper["response"].(map[string]interface{})
	items := response["items"].([]interface{})

	for _, v := range items {
		beer := v.(map[string]interface{})["beer"].(map[string]interface{})
		beerID := int64(beer["bid"].(float64))
		date := string(v.(map[string]interface{})["created_at"].(string))
		cdate, _ := time.Parse(time.RFC1123Z, date)
		nbeer := &pb.Beer{Id: beerID, DrinkDate: cdate.Unix()}
		values = append(values, nbeer)
	}

	return values, nil
}

// GetRecentDrinks Gets the most recent drinks from untappd
func (u *Untappd) GetRecentDrinks(fetcher httpResponseFetcher, converter responseConverter, date int64) []int64 {
	var unmarshaller unmarshaller = mainUnmarshaller{}

	var ret []int64

	text := u.getVenuePage(fetcher, converter, 2194560)
	drinks, _ := u.convertPageToDrinks(text, unmarshaller)

	for _, k := range drinks {
		if date < k.DrinkDate {
			ret = append(ret, k.Id)
		}
	}

	return ret
}

// GetBeerDetails Determines the name of the beer from the id
func (u *Untappd) GetBeerDetails(id int64) *pb.Beer {
	text := u.getBeerPage(u.f, u.c, int(id))
	name := convertPageToName(text, u.u)
	abv := convertPageToABV(text, u.u)
	return &pb.Beer{Name: name, Abv: abv, Id: id}
}

func (u *Untappd) convertDrinkListToBeers(page string, unmarshaller unmarshaller) []*pb.Beer {
	var resp []interface{}
	json.Unmarshal([]byte(page), &resp)

	beers := make([]*pb.Beer, len(resp))
	for i, c := range resp {
		m := c.(map[string]interface{})
		id := m["beer_url"].(string)
		elems := strings.Split(id, "/")
		val, _ := strconv.Atoi(elems[len(elems)-1])
		cid := m["checkin_url"].(string)
		celems := strings.Split(cid, "/")
		cval, _ := strconv.Atoi(celems[len(elems)-1])
		created := m["created_at"].(string)
		layout := "2006-01-02 15:04:05"
		t, _ := time.Parse(layout, created)
		beers[i] = &pb.Beer{Name: m["beer_name"].(string), Id: int64(val), DrinkDate: t.Unix(), CheckinId: int32(cval)}
	}

	return beers
}

func (u *Untappd) getLastBeers(f httpResponseFetcher, c responseConverter, un unmarshaller, lastID int32) []*pb.Beer {
	page, err := u.getUserPage(f, c, "brotherlogic", int(lastID))
	if err != nil {
		return []*pb.Beer{}
	}
	list, _ := u.convertUserPageToDrinks(page, un)
	return list
}
