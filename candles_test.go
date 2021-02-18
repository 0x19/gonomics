package gonomics

import (
	"encoding/csv"
	"os"
	"testing"
	"time"
)

// TestGetCandles tests candles fetching function.
func TestGetCandles(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /candles API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /candles API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
		cReqJSON := CandlesRequest{
			Interval: "1d",
			Currency: "BTC",
			Start:    startTime,
			End:      endTime,
		}
		cRespJSON, err := c.GetCandles(cReqJSON)
		if err != nil {
			t.Error(err)
		}
		if len(cRespJSON) < 1 {
			t.Error("Something is wrong here, length of response array is 0. JSON Format.")
		}
		for _, c := range cRespJSON {
			if c.Timestamp.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. JSON Format.")
			}
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		cReqCSV := CandlesRequest{
			Interval:         "1d",
			Currency:         "BTC",
			Start:            startTime,
			End:              endTime,
			Format:           "csv",
			FileNameWithPath: "./testdata/candles.csv",
		}
		_, err = c.GetCandles(cReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(cReqCSV.FileNameWithPath)
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
		var ts time.Time
		for _, data := range csvData {
			ts, err = time.Parse(time.RFC3339, data[0])
			if err != nil {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. CSV Format.")
			}
			if ts.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. CSV Format.")
			}
		}
	}
}

// TestGetExchangeCandles tests exchange candles fetching function.
func TestGetExchangeCandles(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /exchange_candles API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /exchange_candles API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
		ecReqJSON := ExchangeCandlesRequest{
			Interval: "1d",
			Exchange: "binance",
			Market:   "BTCUSDT",
			Start:    startTime,
			End:      endTime,
		}
		ecRespJSON, err := c.GetExchangeCandles(ecReqJSON)
		if err != nil {
			t.Error(err)
		}
		if len(ecRespJSON) < 1 {
			t.Error("Something is wrong here, length of response array is 0. JSON Format.")
		}
		for _, ec := range ecRespJSON {
			if ec.Timestamp.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. JSON Format.")
			}
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		ecReqCSV := ExchangeCandlesRequest{
			Interval:         "1d",
			Exchange:         "binance",
			Market:           "BTCUSDT",
			Start:            startTime,
			End:              endTime,
			Format:           "csv",
			FileNameWithPath: "./testdata/exchange_candles.csv",
		}
		_, err = c.GetExchangeCandles(ecReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(ecReqCSV.FileNameWithPath)
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
		var ts time.Time
		for _, data := range csvData {
			ts, err = time.Parse(time.RFC3339, data[0])
			if err != nil {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. CSV Format.")
			}
			if ts.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. CSV Format.")
			}
		}
	}
}

// TestGetMarketsCandles tests markets candles fetching function.
func TestGetMarketsCandles(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /markets/candles API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /markets/candles API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
		mcReqJSON := MarketsCandlesRequest{
			Interval: "1d",
			Base:     "ETH",
			Quote:    "BTC",
			Start:    startTime,
			End:      endTime,
		}
		mcRespJSON, err := c.GetMarketsCandles(mcReqJSON)
		if err != nil {
			t.Error(err)
		}
		if len(mcRespJSON) < 1 {
			t.Error("Something is wrong here, length of response array is 0. JSON Format.")
		}
		for _, mc := range mcRespJSON {
			if mc.Timestamp.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. JSON Format.")
			}
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		mcReqCSV := MarketsCandlesRequest{
			Interval:         "1d",
			Base:             "ETH",
			Quote:            "BTC",
			Start:            startTime,
			End:              endTime,
			Format:           "csv",
			FileNameWithPath: "./testdata/markets_candles.csv",
		}
		_, err = c.GetMarketsCandles(mcReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(mcReqCSV.FileNameWithPath)
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
		var ts time.Time
		for _, data := range csvData {
			ts, err = time.Parse(time.RFC3339, data[0])
			if err != nil {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. CSV Format.")
			}
			if ts.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. CSV Format.")
			}
		}
	}
}
