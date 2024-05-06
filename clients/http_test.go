package clients

import (
	"errors"
	"net/url"
	"sort"
	"strconv"
	"testing"

	"github.com/Pallinder/go-randomdata"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/fatih/structs"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var TestSortBy = struct {
	Name        string
	Symbol      string
	MarketCap   string
	Price       string
	TotalSupply string
	MaxSupply   string
}{
	Name:        "name",
	Symbol:      "symbol",
	MarketCap:   "market_cap",
	Price:       "price",
	TotalSupply: "total_supply",
	MaxSupply:   "max_supply",
}

var TestSortType = struct {
	ASC  string
	DESC string
}{
	ASC:  "asc",
	DESC: "desc",
}

func TestAddQueryParamsErr(t *testing.T) {
	expectedErr := &url.Error{
		Op:  "parse",
		URL: "url",
		Err: errors.New("error"),
	}
	patches := ApplyFuncReturn(url.Parse, nil, expectedErr)
	defer patches.Reset()

	_, err := addQueryParams("", map[string]string{})

	assert.EqualErrorf(t, err, "parse \"url\": error", "")
}

func TestAddQueryParams(t *testing.T) {
	testURL := "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/listing"
	sortByParam := randomdata.StringSample(getStructValues(TestSortBy)...)
	sortTypeParam := randomdata.StringSample(getStructValues(TestSortType)...)
	startParam := strconv.Itoa(randomdata.Number(1, 9999))
	limitParam := strconv.Itoa(randomdata.Number(1, 9999))
	getParams := map[string]string{
		"start":    startParam,
		"limit":    limitParam,
		"sortBy":   sortByParam,
		"sortType": sortTypeParam,
	}
	expectedURL := addQueryParamsToExpectedURL(testURL, getParams)

	parsedURL, _ := addQueryParams(testURL, getParams)

	if parsedURL.String() != expectedURL {
		t.Errorf(
			"parsedURL %s is not correct, expected url is %s",
			parsedURL.String(),
			expectedURL,
		)
	}
}

func TestMakeGetRequest(t *testing.T) {
	defer gock.Off()
	host := "server.com"
	expectedURL := "http://" + host
	expectedParams := map[string]string{
		"param1": "value1",
	}
	gock.New(expectedURL).Get("/").Reply(200).JSON(map[string]string{})
	patches := ApplyFuncReturn(addQueryParams, &url.URL{Host: host}, nil)
	defer patches.Reset()

	httpClient := getHTTPClient()
	res, err := httpClient.MakeGetRequest(expectedURL, expectedParams)

	if err != nil {
		t.Fatalf("Err has to be nil but err is %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatal("Status code has to be 200")
	}
}

func addQueryParamsToExpectedURL(url string, getParams map[string]string) string {
	keys := make([]string, 0, len(getParams))
	for k := range getParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i == 0 {
			url += "?" + k + "=" + getParams[k]
		} else {
			url += "&" + k + "=" + getParams[k]
		}
	}
	return url
}

func getStructValues(item interface{}) []string {
	values := structs.Values(item)
	convertedValues := make([]string, len(values))
	for i, v := range values {
		convertedValues[i] = v.(string)
	}
	return convertedValues
}
