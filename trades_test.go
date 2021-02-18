package gonomics

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestGetTrades tests trades fetching function.
func TestGetTrades(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /trades API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /trades API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		fromTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		tReqJSON := TradesRequest{
			Exchange: "binance",
			Market:   "BTCUSDT",
			Limit:    100,
			Order:    "asc",
			From:     fromTime,
		}
		tRespJSON, err := c.GetTrades(tReqJSON)
		if err != nil {
			t.Error(err)
		}
		if len(tRespJSON) < 1 {
			t.Error("Something is wrong here, length of response array is 0. JSON Format.")
		}
		for _, tr := range tRespJSON {
			if tr.Volume == 0 {
				t.Error("Something is wrong here, value of response arrays volume is 0. JSON Format.")
			}
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		tReqCSV := TradesRequest{
			Exchange:         "binance",
			Market:           "BTCUSDT",
			Limit:            100,
			Order:            "asc",
			From:             fromTime,
			Format:           "csv",
			FileNameWithPath: "./testdata/trades.csv",
		}
		_, err = c.GetTrades(tReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(tReqCSV.FileNameWithPath)
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
			tr, err := strconv.ParseFloat(data[3], 64)
			if err != nil {
				t.Error("Something is wrong here, value of response arrays volume is 0. CSV Format.")
			}
			if tr == 0 {
				t.Error("Something is wrong here, value of response arrays volume is 0. CSV Format.")
			}
		}
	}
}
