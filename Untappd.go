package main

import "bufio"
import "encoding/json"
import "errors"
import "fmt"
import "log"
import "net/http"
import "os"
import "io/ioutil"
import "strconv"
import "strings"
import "time"

import pb "github.com/brotherlogic/beer/proto"

var untappdKey string
var untappdSecret string

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
type mainFetcher struct{}

func (fetcher mainFetcher) Fetch(url string) (*http.Response, error) {
	return http.Get(url)
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

// LoadCache - loads the cache from a given file
func LoadCache(folder string) {
	fileName := folder + "/untappd.metadata"

	file, err := os.Open(fileName)

	if err != nil {
		log.Printf("Error loading cache: %v - %v\n", fileName, err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			elems := strings.Split(line, "~")
			val, _ := strconv.Atoi(elems[0])
			beerMap[val] = elems[1]
		}
	}
}

// SaveCache - saves the cache to a given file
func SaveCache(folder string) {
	os.Mkdir(folder, 0777)
	fileName := folder + "/untappd.metadata"
	file, err := os.Create(fileName)

	defer file.Close()

	if err == nil {
		for k, v := range beerMap {
			fmt.Fprintf(file, "%v~%v\n", k, v)
		}
	} else {
		log.Printf("Failed opening file %v - %v\n", fileName, err)
	}
}

func getBeerPage(fetcher httpResponseFetcher, converter responseConverter, id int) string {
	url := "https://api.untappd.com/v4/beer/info/BID?client_id=CLIENTID&client_secret=CLIENTSECRET&compact=true"
	url = strings.Replace(url, "BID", strconv.Itoa(id), 1)
	url = strings.Replace(url, "CLIENTID", untappdKey, 1)
	url = strings.Replace(url, "CLIENTSECRET", untappdSecret, 1)

	response, err := fetcher.Fetch(url)

	if err != nil {
		log.Printf("Failed on getBeerPage: %q\n", err)
	} else {
		contents, err := converter.Convert(response)
		if err != nil {
			log.Printf("%q\n", err)
		} else {
			return string(contents)
		}
	}

	return "Failed to retrieve " + strconv.Itoa(id)
}

func getVenuePage(fetcher httpResponseFetcher, converter responseConverter, id int) string {
	url := "https://api.untappd.com/v4/venue/info/VID?client_id=CLIENTID&client_secret=CLIENTSECRET"
	url = strings.Replace(url, "VID", strconv.Itoa(id), 1)
	url = strings.Replace(url, "CLIENTID", untappdKey, 1)
	url = strings.Replace(url, "CLIENTSECRET", untappdSecret, 1)

	response, err := fetcher.Fetch(url)

	log.Printf("Getting venue page: %v\n", url)

	if err != nil {
		log.Printf("Failed on getVenuePage: %q\n", err)
	} else {
		contents, err := converter.Convert(response)
		if err != nil {
			log.Printf("%q\n", err)
		} else {
			return string(contents)
		}
	}

	return "Failed to retrieve " + strconv.Itoa(id)
}

func convertPageToName(page string, unmarshaller unmarshaller) string {
	var mapper map[string]interface{}
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		log.Printf("%q\n", err)
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

func convertPageToDrinks(page string, unmarshaller unmarshaller) ([]pb.Beer, error) {
	log.Printf("RUNNING\n")

	var mapper map[string]interface{}
	var values []pb.Beer
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		log.Printf("ERROR: %q\n", err)
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
		nbeer := pb.Beer{Id: beerID, DrinkDate: cdate.Unix()}
		values = append(values, nbeer)
	}

	return values, nil
}

// GetRecentDrinks Gets the most recent drinks from untappd
func GetRecentDrinks(fetcher httpResponseFetcher, converter responseConverter, date int64) []int64 {
	var unmarshaller unmarshaller = mainUnmarshaller{}

	var ret []int64

	text := getVenuePage(fetcher, converter, 2194560)
	drinks, _ := convertPageToDrinks(text, unmarshaller)

	for _, k := range drinks {
		if date < k.DrinkDate {
			ret = append(ret, k.Id)
		}
	}

	return ret
}

// GetBeerName Determines the name of the beer from the id
func GetBeerName(id int) string {

	//Check the cache
	if val, ok := beerMap[id]; ok {
		return val
	}

	var fetcher httpResponseFetcher = mainFetcher{}
	var converter responseConverter = mainConverter{}
	var unmarshaller unmarshaller = mainUnmarshaller{}
	text := getBeerPage(fetcher, converter, id)

	name := convertPageToName(text, unmarshaller)
	beerMap[id] = name
	return name
}
