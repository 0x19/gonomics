package gonomics

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestGetExchangeRates tests exchange rates fetching function.
func TestGetExchangeRates(t *testing.T) {
	t.Log("Testing /exchange-rates API endpoint.")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	// Json format.
	t.Log("Testing for JSON format.")
	erReqJSON := ExchangeRatesRequest{}
	erRespJSON, err := c.GetExchangeRates(erReqJSON)
	if err != nil {
		t.Error(err)
	}
	if len(erRespJSON) < 1 {
		t.Error("Something is wrong here, length of response array is 0. JSON Format.")
	}
	for _, er := range erRespJSON {
		if er.Rate == 0 {
			t.Error("Something is wrong here, value of response arrays rate is 0. JSON Format.")
		}
	}

	// CSV format.
	t.Log("Testing for CSV format.")
	erReqCSV := ExchangeRatesRequest{
		Format:           "csv",
		FileNameWithPath: "./testdata/exchange_rates.csv",
	}
	_, err = c.GetExchangeRates(erReqCSV)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(erReqCSV.FileNameWithPath)
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
		r, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			t.Error("Something is wrong here, value of response arrays rate is 0. CSV Format.")
		}
		if r == 0 {
			t.Error("Something is wrong here, value of response arrays volume is 0. CSV Format.")
		}
	}
}

// TestGetExchangeRatesHistory tests exchange-rates history fetching function.
func TestGetExchangeRatesHistory(t *testing.T) {
	t.Log("Testing /exchange-rates/history API endpoint.")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	// Json format.
	t.Log("Testing for JSON format.")
	startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
	erhReqJSON := ExchangeRatesHistoryRequest{
		Currency: "BTC",
		Start:    startTime,
		End:      endTime,
	}
	erhRespJSON, err := c.GetExchangeRatesHistory(erhReqJSON)
	if err != nil {
		t.Error(err)
	}
	if len(erhRespJSON) < 1 {
		t.Error("Something is wrong here, length of response array is 0. JSON Format.")
	}
	for _, erh := range erhRespJSON {
		if erh.Rate == 0 {
			t.Error("Something is wrong here, value of response arrays rate is 0. JSON Format.")
		}
	}

	// CSV format.
	t.Log("Testing for CSV format.")
	erhReqCSV := ExchangeRatesHistoryRequest{
		Currency:         "BTC",
		Start:            startTime,
		End:              endTime,
		Format:           "csv",
		FileNameWithPath: "./testdata/exchange-rates_history.csv",
	}
	_, err = c.GetExchangeRatesHistory(erhReqCSV)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(erhReqCSV.FileNameWithPath)
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
		r, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			t.Error("Something is wrong here, value of response arrays rate is 0. CSV Format.")
		}
		if r == 0 {
			t.Error("Something is wrong here, value of response arrays volume is 0. CSV Format.")
		}
	}
}
