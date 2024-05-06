package clients

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// coinmarketcap API urls
const (
	BaseAPIURLV3        = "https://api.coinmarketcap.com/data-api/v3"
	BaseAPIURLV1        = "https://api.coinmarketcap.com/data-api/v1"
	CryptoListingAPIURL = BaseAPIURLV3 + "/cryptocurrency/listing"
	CryptoMapAPIURL     = BaseAPIURLV1 + "/cryptocurrency/map"
	CryptoChartAPIURL   = BaseAPIURLV3 + "/cryptocurrency/detail/chart"
)

// SortBy CryptoListingAPIURL params
var SortBy = struct {
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

// SortType CryptoListingAPIURL params
var SortType = struct {
	ASC  string
	DESC string
}{
	ASC:  "asc",
	DESC: "desc",
}

// GetCoinmarketcapAPIClient is a function to obtain a client for public coinmarketcap API
func GetCoinmarketcapAPIClient() CoinmarketcapAPIClient {
	return CoinmarketcapAPIClient{httpClient: getHTTPClient()}
}

// CoinmarketcapAPIClient is a client for public coinmarketcap API
type CoinmarketcapAPIClient struct {
	httpClient customHTTPClient
}

// GetSymbols is a method to obtain crypto symbols with basic info
func (cac CoinmarketcapAPIClient) GetSymbols(
	start, limit int,
	sortBy, sortType string,
) ([]CoinmarketcapSymbol, error) {
	var responseData CoinmarketcapListResponse
	queryParams := map[string]string{
		"start":    strconv.Itoa(start),
		"limit":    strconv.Itoa(limit),
		"sortBy":   sortBy,
		"sortType": sortType,
	}
	res, err := cac.httpClient.MakeGetRequest(CryptoListingAPIURL, queryParams)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}
	return responseData.Data.CryptoCurrencyList, nil
}

// GetCoinIDBySymbol is a method to obtain coin ID by symbol name
func (cac CoinmarketcapAPIClient) GetCoinIDBySymbol(symbol string) (int, error) {
	var responseData CryptoMapListResponse
	queryParams := map[string]string{
		"symbol": symbol,
	}
	res, err := cac.httpClient.MakeGetRequest(CryptoMapAPIURL, queryParams)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return -1, err
	}
	return responseData.Data[0].ID, nil
}

// GetSymbolChart is a method to get points to draw a chart
func (cac CoinmarketcapAPIClient) GetSymbolChart(
	symbolID int,
	rangeType string,
) (map[string]valuesChart, error) {
	var responseData CryptoChartResponse
	queryParams := map[string]string{
		"id":    fmt.Sprintf("%d", symbolID),
		"range": rangeType,
	}
	res, err := cac.httpClient.MakeGetRequest(CryptoChartAPIURL, queryParams)
	if err != nil {
		return map[string]valuesChart{}, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return map[string]valuesChart{}, err
	}
	return responseData.Data.Points, nil
}

// CryptoChartResponse is a data structure with points to draw a chart
type CryptoChartResponse struct {
	Data timestampsChart `json:"data"`
}

type timestampsChart struct {
	Points map[string]valuesChart `json:"points"`
}

type valuesChart struct {
	V []float64 `json:"v"`
	C []float64 `json:"c"`
}

// CoinmarketcapListResponse is just a wrapper for CoinmarketcapListDataResponse
type CoinmarketcapListResponse struct {
	Data CoinmarketcapListDataResponse `json:"data"`
}

// CoinmarketcapListDataResponse is just a wrapper for CryptoCurrencyList
type CoinmarketcapListDataResponse struct {
	CryptoCurrencyList []CoinmarketcapSymbol `json:"cryptoCurrencyList"`
	TotalCount         string                `json:"totalCount"`
}

type quotes struct {
	Name             string  `json:"name"`
	Price            float64 `json:"price"`
	Volume24H        float64 `json:"volume24h"`
	MarketCap        float64 `json:"marketCap"`
	PercentChange1H  float64 `json:"percentChange1h"`
	PercentChange24H float64 `json:"percentChange24h"`
	PercentChange7D  float64 `json:"percentChange7d"`
}

// CoinmarketcapSymbol is a data structure that contains basic data about crypto symbol
type CoinmarketcapSymbol struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Symbol            string   `json:"symbol"`
	Slug              string   `json:"slug"`
	CirculatingSupply float64  `json:"circulatingSupply"`
	Quotes            []quotes `json:"quotes"`
}

// CryptoMapListResponse is just a wrapper for cryptoMapSymbol
type CryptoMapListResponse struct {
	Data []cryptoMapSymbol `json:"data"`
}

type cryptoMapSymbol struct {
	ID int `json:"id"`
}
