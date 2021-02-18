package gonomics

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestGetMarkets tests markets fetching function.
func TestGetMarkets(t *testing.T) {
	t.Log("Testing /markets API endpoint.")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	// Json format.
	t.Log("Testing for JSON format.")
	mReqJSON := MarketsRequest{
		Exchange: "binance",
		Base:     []string{"BNB", "LTC"},
		Quote:    []string{"BTC", "ETH"},
	}
	mRespJSON, err := c.GetMarkets(mReqJSON)
	if err != nil {
		t.Error(err)
	}
	if len(mRespJSON) < 1 {
		t.Error("Something is wrong here, length of response array is 0. JSON Format.")
	}
	for _, m := range mRespJSON {
		if m.Base == "" {
			t.Error("Something is wrong here, value of response arrays base is empty. JSON Format.")
		}
	}

	// CSV format.
	t.Log("Testing for CSV format.")
	mReqCSV := MarketsRequest{
		Exchange:         "binance",
		Base:             []string{"BNB", "LTC"},
		Quote:            []string{"BTC", "ETH"},
		Format:           "csv",
		FileNameWithPath: "./testdata/markets.csv",
	}
	_, err = c.GetMarkets(mReqCSV)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(mReqCSV.FileNameWithPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	csvData, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Error(err)
	}
	if len(csvData) < 1 {
		t.Error("Something is wrong here, length of response array is 0. CSV Format.")
	}
	for _, data := range csvData {
		if data[1] == "" {
			t.Error("Something is wrong here, value of response arrays base is empty. CSV Format.")
		}
	}
}

// TestGetMarketsCapHistory tests market-cap history fetching function.
func TestGetMarketsCapHistory(t *testing.T) {
	t.Log("Testing /market-cap/history API endpoint. (Partial Paid Plan)")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	// Json format.
	t.Log("Testing for JSON format.")
	startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
	mchReqJSON := MarketsCapHistoryRequest{
		Start:   startTime,
		End:     endTime,
		Convert: "EUR",
	}
	mchRespJSON, err := c.GetMarketsCapHistory(mchReqJSON)
	if err != nil {
		t.Error(err)
	}
	if len(mchRespJSON) < 1 {
		t.Error("Something is wrong here, length of response array is 0. JSON Format.")
	}
	for _, mch := range mchRespJSON {
		if mch.MarketCap == 0 {
			t.Error("Something is wrong here, value of response arrays market-cap is 0. JSON Format.")
		}
	}

	// CSV format.
	t.Log("Testing for CSV format.")
	mchReqCSV := MarketsCapHistoryRequest{
		Start:            startTime,
		End:              endTime,
		Convert:          "EUR",
		Format:           "csv",
		FileNameWithPath: "./testdata/market-cap_history.csv",
	}
	_, err = c.GetMarketsCapHistory(mchReqCSV)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(mchReqCSV.FileNameWithPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	csvData, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Error(err)
	}
	if len(csvData) < 1 {
		t.Error("Something is wrong here, length of response array is 0. CSV Format.")
	}
	for _, data := range csvData {
		mc, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			t.Error("Something is wrong here, value of response arrays market-cap is 0. CSV Format.")
		}
		if mc == 0 {
			t.Error("Something is wrong here, value of response arrays market-cap is 0. CSV Format.")
		}
	}
}

// TestGetExchangeMarketsTicker tests exchange-markets ticker fetching function.
func TestGetExchangeMarketsTicker(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /exchange-markets/ticker API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /exchange-markets/ticker API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		emtReq := ExchangeMarketsTickerRequest{
			Interval: []string{"1d", "30d"},
			Currency: []string{"BTC", "ETH"},
			Exchange: []string{"binance", "gdax"},
			Convert:  "BTC",
			PerPage:  100,
			Page:     1,
		}
		emtResp, err := c.GetExchangeMarketsTicker(emtReq)
		if err != nil {
			t.Error(err)
		}
		if len(emtResp) < 1 {
			t.Error("Something is wrong here, length of response array is 0.")
		}
		for _, emt := range emtResp {
			if emt.Three0D.Volume == 0 {
				t.Error("Something is wrong here, value of response arrays 30d volume is 0.")
			}
		}
	}
}
