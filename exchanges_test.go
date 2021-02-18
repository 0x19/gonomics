package gonomics

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// TestGetExchangesTicker tests exchanges ticker fetching function.
func TestGetExchangesTicker(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /exchanges/ticker API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /exchanges/ticker API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		etReq := ExchangesTickerRequest{
			Ids:      []string{"BTC", "ETH"},
			Interval: []string{"1d", "30d"},
			PerPage:  100,
			Page:     1,
		}
		etResp, err := c.GetExchangesTicker(etReq)
		if err != nil {
			t.Error(err)
		}
		if len(etResp) < 1 {
			t.Error("Something is wrong here, length of response array is 0.")
		}
		for _, et := range etResp {
			if et.Three0D.Volume == 0 {
				t.Error("Something is wrong here, value of response arrays 30d volume is 0.")
			}
		}
	}
}

// TestGetExchangesVolumeHistory tests volume history fetching function.
func TestGetExchangesVolumeHistory(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /exchanges/volume/history API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /exchanges/volume/history API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
		evhReqJSON := ExchangesVolumeHistoryRequest{
			Exchange: "binance",
			Start:    startTime,
			End:      endTime,
			Convert:  "EUR",
		}
		evhRespJSON, err := c.GetExchangesVolumeHistory(evhReqJSON)
		if err != nil {
			t.Error(err)
		}
		if len(evhRespJSON) < 1 {
			t.Error("Something is wrong here, length of response array is 0. JSON Format.")
		}
		for _, evh := range evhRespJSON {
			if evh.Volume == 0 {
				t.Error("Something is wrong here, value of response arrays volume is 0. JSON Format.")
			}
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		evhReqCSV := ExchangesVolumeHistoryRequest{
			Exchange:         "binance",
			Start:            startTime,
			End:              endTime,
			Convert:          "EUR",
			Format:           "csv",
			FileNameWithPath: "./testdata/exchanges_volume_history.csv",
		}
		_, err = c.GetExchangesVolumeHistory(evhReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(evhReqCSV.FileNameWithPath)
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
			v, err := strconv.ParseFloat(data[1], 64)
			if err != nil {
				t.Error("Something is wrong here, value of response arrays volume is 0. CSV Format.")
			}
			if v == 0 {
				t.Error("Something is wrong here, value of response arrays volume is 0. CSV Format.")
			}
		}
	}
}

// TestGetExchangesMetadata tests exchanges metadata fetching function.
func TestGetExchangesMetadata(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /exchanges API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /exchanges API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		expectedOutput := []ExchangesMetadataResponse{
			{ID: "binance", Name: "Binance"},
		}

		// Json format.
		t.Log("Testing for JSON format.")
		emReqJSON := ExchangesMetadataRequest{
			Ids:        []string{"binance"},
			Attributes: []string{"id", "name"},
		}
		emRespJSON, err := c.GetExchangesMetadata(emReqJSON)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(emRespJSON, expectedOutput) {
			t.Error("Something is wrong here, actual output is not matching with expected output. JSON Format.")
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		emReqCSV := ExchangesMetadataRequest{
			Ids:              []string{"binance"},
			Attributes:       []string{"id", "name"},
			Format:           "csv",
			FileNameWithPath: "./testdata/exchanges_metadata.csv",
		}
		_, err = c.GetExchangesMetadata(emReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(emReqCSV.FileNameWithPath)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()
		csvData, err := csv.NewReader(f).ReadAll()
		if err != nil {
			t.Error(err)
		}
		var emRespCSV []ExchangesMetadataResponse
		for _, data := range csvData {
			emRespCSV = append(emRespCSV, ExchangesMetadataResponse{
				ID:   data[0],
				Name: data[1],
			})
		}
		if !reflect.DeepEqual(emRespCSV, expectedOutput) {
			t.Error("Something is wrong here, actual output is not matching with expected output. CSV Format.")
		}
	}
}
