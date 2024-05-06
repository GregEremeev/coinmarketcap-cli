package clients

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
)

func TestGetSymbols(t *testing.T) {
	apiClient := GetCoinmarketcapAPIClient()
	rawContent := []byte(
		`{"data": {"cryptoCurrencyList": [
{"id": 30953, "name": "name", "symbol": "BTC",
"slug": "BTC-Slug", "total_supply": 10, "max_supply": 10,
"quotes": [{"name": "USD", "price": 57348.550077281856, "volume24h": 49666217758.26597,
"marketCap": 1129355648858.249, "percentChange1h": -0.70360978,
"percentChange24h": -3.87692361, "percentChange7d": -10.61822801}]}]}}`,
	)
	patches := ApplyMethodReturn(
		apiClient.httpClient,
		"MakeGetRequest",
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(rawContent))},
		nil,
	)
	defer patches.Reset()

	symbols, err := apiClient.GetSymbols(
		1,
		10,
		"sortBy",
		"sortTyp",
	)

	if err != nil {
		t.Fatalf("err has to be nil but it is %v", err)
	}
	if len(symbols) == 0 {
		t.Fatal("there are no symbols")
	}
}

func TestGetCoinIdBySymbol(t *testing.T) {
	apiClient := GetCoinmarketcapAPIClient()
	rawContent := []byte(
		`{"data":[{"id":1,"name":"Bitcoin",
"symbol":"BTC","slug":"bitcoin","is_active":1,
"first_historical_data":"2010-07-13T00:05:00.000Z",
"last_historical_data":"2024-05-04T08:05:00.000Z","rank":1}],
"status":{"timestamp":"2024-05-04T08:16:35.616Z",
"error_code":"0","error_message":"SUCCESS","elapsed":"1","credit_count":0}}`)

	patches := ApplyMethodReturn(
		apiClient.httpClient,
		"MakeGetRequest",
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(rawContent))},
		nil,
	)
	defer patches.Reset()

	id, err := apiClient.GetCoinIDBySymbol("btc")

	if err != nil {
		t.Fatalf("err has to be nil but it is %v", err)
	}
	if id != 1 {
		t.Fatalf("id is not correct")
	}
}

func TestGetSymbolChart(t *testing.T) {
	apiClient := GetCoinmarketcapAPIClient()
	rawContent := []byte(
		`{"data":{"points":{"1279065600":{"v":[0.05815725,261.54,196180.04855025,1,3373269],
"c":[0.05815725,261.54,196180.04855025]},
"1279756800":{"v":[0.07417871,2167.06,258429.35330125,1,3483875],
"c":[0.07417871,2167.06,258429.35330125]},
"1280361600":{"v":[0.06785855,8091.7,241028.48020875,1,3551925]}}},
"status":{"timestamp":"2024-05-04T09:02:45.581Z","error_code":"0",
"error_message":"SUCCESS","elapsed":"21","credit_count":0}}
`)

	patches := ApplyMethodReturn(
		apiClient.httpClient,
		"MakeGetRequest",
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(rawContent))},
		nil,
	)
	defer patches.Reset()

	points, err := apiClient.GetSymbolChart(1, "1d")

	if err != nil {
		t.Fatalf("err has to be nil but it is %v", err)
	}
	if len(points) == 0 {
		t.Fatalf("ponts are not correct")
	}
	if points["1279065600"].V[0] != 0.05815725 {
		t.Fatalf("ponts are not correct")
	}
	if points["1279065600"].C[2] != 196180.04855025 {
		t.Fatalf("ponts are not correct")
	}
}
