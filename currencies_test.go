package gonomics

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
	"time"
)

// TestGetCurrenciesTicker tests currencies ticker fetching function.
func TestGetCurrenciesTicker(t *testing.T) {
	t.Log("Testing /currencies/ticker API endpoint. (Partial Paid Plan)")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	ctReq := CurrenciesTickerRequest{
		Ids:      []string{"BTC", "ETH"},
		Interval: []string{"1d", "30d"},
		PerPage:  100,
		Page:     1,
	}
	ctResp, err := c.GetCurrenciesTicker(ctReq)
	if err != nil {
		t.Error(err)
	}
	if len(ctResp) < 1 {
		t.Error("Something is wrong here, length of response array is 0.")
	}
	for _, ct := range ctResp {
		if ct.Three0D.Volume == 0 {
			t.Error("Something is wrong here, value of response arrays 30d volume is 0.")
		}
	}
}

// TestGetCurrenciesMetadata tests currencies metadata fetching function.
func TestGetCurrenciesMetadata(t *testing.T) {
	t.Log("Testing /currencies API endpoint.")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	expectedOutput := []CurrenciesMetadataResponse{
		{ID: "BTC", Name: "Bitcoin"},
		{ID: "ETH", Name: "Ethereum"},
	}

	// Json format.
	t.Log("Testing for JSON format.")
	cmReqJSON := CurrenciesMetadataRequest{
		Ids:        []string{"BTC", "ETH"},
		Attributes: []string{"id", "name"},
	}
	cmRespJSON, err := c.GetCurrenciesMetadata(cmReqJSON)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(cmRespJSON, expectedOutput) {
		t.Error("Something is wrong here, actual output is not matching with expected output. JSON Format.")
	}

	// CSV format.
	t.Log("Testing for CSV format.")
	cmReqCSV := CurrenciesMetadataRequest{
		Ids:              []string{"BTC", "ETH"},
		Attributes:       []string{"id", "name"},
		Format:           "csv",
		FileNameWithPath: "./testdata/currencies_metadata.csv",
	}
	_, err = c.GetCurrenciesMetadata(cmReqCSV)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(cmReqCSV.FileNameWithPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	csvData, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Error(err)
	}
	var cmRespCSV []CurrenciesMetadataResponse
	for _, data := range csvData {
		cmRespCSV = append(cmRespCSV, CurrenciesMetadataResponse{
			ID:   data[0],
			Name: data[1],
		})
	}
	if !reflect.DeepEqual(cmRespCSV, expectedOutput) {
		t.Error("Something is wrong here, actual output is not matching with expected output. CSV Format.")
	}
}

// TestGetCurrenciesSparkline tests currencies sparkline fetching function.
func TestGetCurrenciesSparkline(t *testing.T) {
	t.Log("Testing /currencies/sparkline API endpoint.")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	startTime, _ := time.Parse(time.RFC3339, "2021-02-14T00:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2021-02-16T00:00:00Z")
	csReq := CurrenciesSparklineRequest{
		Ids:     []string{"BTC", "ETH"},
		Start:   startTime,
		End:     endTime,
		Convert: "USD",
	}
	csResp, err := c.GetCurrenciesSparkline(csReq)
	if err != nil {
		t.Error(err)
	}
	if len(csResp) < 1 {
		t.Error("Something is wrong here, length of response array is 0.")
	}
	for _, cs := range csResp {
		if len(cs.Prices) < 1 {
			t.Error("Something is wrong here, length of response arrays prices is 0.")
		}
	}
}

// TestGetCurrenciesSupplyHistory tests currencies supply history fetching function.
func TestGetCurrenciesSupplyHistory(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /supplies/history API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /supplies/history API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		startTime, _ := time.Parse(time.RFC3339, "2021-02-14T00:00:00Z")
		endTime, _ := time.Parse(time.RFC3339, "2021-02-16T00:00:00Z")
		cshReqJSON := CurrenciesSupplyHistoryRequest{
			Currency: "BTC,ETH",
			Start:    startTime,
			End:      endTime,
		}
		cshRespJSON, err := c.GetCurrenciesSupplyHistory(cshReqJSON)
		if err != nil {
			t.Error(err)
		}
		if len(cshRespJSON) < 1 {
			t.Error("Something is wrong here, length of response array is 0. JSON Format.")
		}
		for _, csh := range cshRespJSON {
			if csh.Timestamp.IsZero() {
				t.Error("Something is wrong here, value of response arrays timestamp is 0. JSON Format.")
			}
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		cshReqCSV := CurrenciesSupplyHistoryRequest{
			Currency:         "BTC,ETH",
			Start:            startTime,
			End:              endTime,
			Format:           "csv",
			FileNameWithPath: "./testdata/currencies_supply_history.csv",
		}
		_, err = c.GetCurrenciesSupplyHistory(cshReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(cshReqCSV.FileNameWithPath)
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
