package main

import (
	"errors"
	"log"
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

type fileFetcher struct{}

func (fetcher fileFetcher) Fetch(url string) (*http.Response, error) {
	strippedURL := strings.Replace(strings.Replace(url[24:], "?", "_", -1), "&", "_", -1)
	data, err := os.Open("testdata/" + strippedURL)
	log.Printf("Loading %v", "testdata/"+strippedURL)
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

func TestGetBeerName(t *testing.T) {
	log.Printf("Running TESTGETBEERNAME maybe")
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	beerName := u.GetBeerName(7936, fileFetcher{}, mainConverter{}, mainUnmarshaller{})
	if beerName != "Firestone Walker Brewing Company - Parabola" {
		t.Errorf("Beer name %q is not firestone, parabola\n", beerName)
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
	beers, err := convertPageToDrinks(beerPage, stubFailUnmarshaller{})
	if err == nil {
		t.Errorf("Bad unmarshal has not failed: %v", beers)
	}
}

func TestGetBeerPage200Fail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 0)
	beers, err := convertPageToDrinks(beerPage, mainUnmarshaller{})
	if err == nil {
		t.Errorf("Bad unmarshal has not failed: %v", beers)
	}
}

func TestGetBeerName200Fail(t *testing.T) {
	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	var fetcher = fileFetcher{}
	var converter = mainConverter{}
	beerPage := u.getBeerPage(fetcher, converter, 0)
	name := convertPageToName(beerPage, mainUnmarshaller{})
	if strings.Contains(name, "Firestone") {
		t.Errorf("Get name worked: %v", name)
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

func TestGetBeerFromCache(t *testing.T) {
	beerMap = make(map[int]string)
	u := &Untappd{untappdID: "blah", untappdSecret: "blah"}
	beerName := u.GetBeerName(7936, fileFetcher{}, mainConverter{}, mainUnmarshaller{})
	if strings.Contains(beerName, "Firestone") {
		t.Errorf("Beer has been retrieved despite no caching / credentials\n")
	}

	//Insert the beer into the cache
	cacheBeer(7936, "Firestone Walker Brewing Company - Parabola")

	beerName = u.GetBeerName(7936, fileFetcher{}, mainConverter{}, mainUnmarshaller{})
	if !strings.Contains(beerName, "Firestone") {
		t.Errorf("Beer name cannot be retrieved despite caching")
	}
}

func TestCacheFail(t *testing.T) {
	cacheBeer(7936, "Firestone Walker Brewing Company - Parabola")
	SaveCache("madeup_folder/makeup_folder2")
	LoadCache("madeup_folder/makeup_folder2")
}

func TestGetBeerFromSavedCache(t *testing.T) {
	cacheBeer(7936, "Firestone Walker Brewing Company - Parabola")
	SaveCache("testing_cache")

	beerMap = make(map[int]string)
	LoadCache("testing_cache")

	u := &Untappd{untappdID: "testid", untappdSecret: "testsecret"}
	beerName := u.GetBeerName(7936, fileFetcher{}, mainConverter{}, mainUnmarshaller{})
	if !strings.Contains(beerName, "Firestone") {
		t.Errorf("Cache reload failed\n")
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
