package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
)

type stubFailUnmarshaller struct{}

func (unmarshaller stubFailUnmarshaller) Unmarshal(int []byte, resp *map[string]interface{}) error {
	return errors.New("Built to fail")
}

type stubPassFetcher struct{}

func (fetcher stubPassFetcher) Fetch(url string) (*http.Response, error) {
	var resp = &http.Response{}
	return resp, nil
}

type stubFailConverter struct{}

func (converter stubFailConverter) Convert(response *http.Response) ([]byte, error) {
	return make([]byte, 0), errors.New("Built to fail")
}

type fileFetcher struct {
	fail bool
}

func (fetcher fileFetcher) Fetch(url string) (*http.Response, error) {
	if fetcher.fail {
		return nil, fmt.Errorf("Built to fail")
	}
	strippedURL := strings.Replace(strings.Replace(url[24:], "?", "_", -1), "&", "_", -1)
	data, err := os.Open("testdata/" + strippedURL)
	if err != nil {
		return nil, err
	}
	response := &http.Response{}
	response.Body = data
	return response, nil
}

type stubFailFetcher struct{}

func (fetcher stubFailFetcher) Fetch(url string) (*http.Response, error) {
	var resp = &http.Response{}
	var err = errors.New("Built to fail")
	return resp, err
}

type blankVenuePageFetcher struct{}

func (fetcher blankVenuePageFetcher) Fetch(url string) (*http.Response, error) {
	return &http.Response{}, nil
}

func GetTestUntappd() *Untappd {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret", f: fileFetcher{}, c: mainConverter{}, u: mainUnmarshaller{}}
	return u
}

func TestGetBeerDetails(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret", f: fileFetcher{}, c: mainConverter{}, u: mainUnmarshaller{}, l: doLog}
	beer := u.GetBeerDetails(7936)
	if beer.Name != "Firestone Walker Brewing Company - Parabola" || beer.Abv != 14.5 || beer.BreweryId != 524 {
		t.Errorf("Beer has come back wrong %v", beer)
	}
}

func TestSearch(t *testing.T) {
	beerMap = make(map[int]string)
	cacheBeer(1234, "Testing Beer")
	cacheBeer(1235, "Made up Thing")

	matches := Search("eer")
	if len(matches) != 1 {
		t.Errorf("Wrong number of matches returned :%v should have been 1, given %v", len(matches), beerMap)
	}
}

func TestGetBeerPage(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 7936)
	if !strings.Contains(beerPage, "Firestone") {
		t.Errorf("Beer page is not being retrieved\n%q\n", beerPage)
	}
}

func TestGetBeerPageConvertFail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 7936)
	beers, err := u.convertPageToDrinks(beerPage, stubFailUnmarshaller{})
	if err == nil {
		t.Errorf("Bad unmarshal has not failed: %v", beers)
	}
}

func TestGetBeerPage200Fail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 0)
	beers, err := u.convertPageToDrinks(beerPage, mainUnmarshaller{})
	if err == nil {
		t.Errorf("Bad unmarshal has not failed: %v", beers)
	}
}

func TestGetBeerName200Fail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 0)
	name, id := convertPageToName(beerPage, mainUnmarshaller{})
	if strings.Contains(name, "Firestone") {
		t.Errorf("Get name worked: %v", name)
	}
	if id != -1 {
		t.Errorf("Get id passed: %v", id)
	}
}

func TestGetBeerAbv200Fail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 0)
	abv := convertPageToABV(beerPage, mainUnmarshaller{})
	if abv != -1 {
		t.Errorf("Get name worked: %v", abv)
	}
}

func TestGetBeerNameBadUnmarshal(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 0)
	name, id := convertPageToName(beerPage, stubFailUnmarshaller{})
	if strings.Contains(name, "Firestone") {
		t.Errorf("Get name worked: %v", name)
	}

	if id != -1 {
		t.Errorf("Pah: %v", id)
	}
}

func TestGetVenuePage(t *testing.T) {
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	drinks := u.getVenuePage(fetcher, converter, 2194560)
	if !strings.Contains(drinks, "Three Burners") {
		t.Errorf("Venue Page is not being retrieved correctly\n%v\n", drinks)
	}
}

func TestGetRecentDrinks(t *testing.T) {
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	drinks := u.GetRecentDrinks(fetcher, converter, 1234)

	if len(drinks) != 15 {
		t.Errorf("Not enough drinks processed %v\n", len(drinks))
	}

	found := false
	for _, v := range drinks {
		if v == 2097330 {
			found = true
		}
	}

	if !found {
		t.Errorf("Beer drunk was not found: %v\n", drinks)
	}
}

func TestGetBeerPageFailHttp(t *testing.T) {
	var fetcher httpResponseFetcher = stubFailFetcher{}
	var converter = mainConverter{}
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	beerPage := u.getBeerPage(fetcher, converter, 7936)
	if !strings.Contains(beerPage, "Failed to retrieve") {
		t.Errorf("Beer page retrieve did not fail\n%q\n", beerPage)
	}
}

func TestGetVenuePageFailHttp(t *testing.T) {
	var fetcher httpResponseFetcher = stubFailFetcher{}
	var converter = mainConverter{}
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	drinks := u.getVenuePage(fetcher, converter, 2194560)
	if !strings.Contains(drinks, "Failed to retrieve") {
		t.Errorf("Beer page retrieve did not fail\n%v\n", drinks)
	}
}

func TestGetBeerPageConvertHttp(t *testing.T) {
	var fetcher httpResponseFetcher = stubPassFetcher{}
	var converter = stubFailConverter{}
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	beerPage := u.getBeerPage(fetcher, converter, 7936)
	if !strings.Contains(beerPage, "Failed to retrieve") {
		t.Errorf("Beer page retrieve did not fail\n%q\n", beerPage)
	}
}

func TestGetVenuePageConvertHttpFail(t *testing.T) {
	var fetcher httpResponseFetcher = stubPassFetcher{}
	var converter = stubFailConverter{}
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	venuePage := u.getVenuePage(fetcher, converter, 2194560)
	if !strings.Contains(venuePage, "Failed to retrieve") {
		t.Errorf("Venue page retrieve did not fail\n%q\n", venuePage)
	}
}

func TestGetUserPage(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	u.l = doLog
	text, _ := u.getUserPage(fileFetcher{}, mainConverter{}, "brotherlogic")

	if !strings.Contains(text, "Simon") {
		t.Errorf("User page retrieve failed: %v", text)
	}

	vals, err := u.convertUserPageToDrinks(text, mainUnmarshaller{})

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(vals) != 25 {
		t.Errorf("Wrong number pulled: %v", len(vals))
	}
}

func TestBadUserPage(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	u.l = doLog
	text, _ := u.getUserPage(fileFetcher{}, mainConverter{}, "blahdyblah")
	_, err := u.convertUserPageToDrinks(text, mainUnmarshaller{})

	if err == nil {
		t.Errorf("Did not error")
	}
}

func TestBadUserPageUnmarshal(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	u.l = doLog
	text, _ := u.getUserPage(fileFetcher{}, mainConverter{}, "brotherlogic")
	_, err := u.convertUserPageToDrinks(text, stubFailUnmarshaller{})

	if err == nil {
		t.Errorf("Did not error")
	}
}

func TestBadUserPageUnmarshalFail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	u.l = doLog
	_, err := u.getUserPage(fileFetcher{fail: true}, mainConverter{}, "brotherlogic")

	if err == nil {
		t.Errorf("Did not error")
	}
}

func TestGetLastBeersFail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	u.l = doLog
	beers, _ := u.getLastBeers(fileFetcher{fail: true}, mainConverter{}, mainUnmarshaller{})

	if len(beers) != 0 {
		t.Errorf("Did not error")
	}
}

func TestGetLastBeers(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	u.l = doLog
	beers, _ := u.getLastBeers(fileFetcher{}, mainConverter{}, mainUnmarshaller{})

	if len(beers) == 0 {
		t.Errorf("Did not ru")
	}
}
