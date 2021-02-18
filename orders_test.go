package gonomics

import (
	"encoding/csv"
	"os"
	"testing"
	"time"
)

// TestGetOrdersSnapshot tests orders snapshot fetching function.
func TestGetOrdersSnapshot(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /orders/snapshot API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /orders/snapshot API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		// Json format.
		t.Log("Testing for JSON format.")
		atTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		osReqJSON := OrdersSnapshotRequest{
			Exchange: "binance",
			Market:   "BTCUSDT",
			At:       atTime,
		}
		osRespJSON, err := c.GetOrdersSnapshot(osReqJSON)
		if err != nil {
			t.Error(err)
		}
		if osRespJSON.Timestamp.IsZero() {
			t.Error("Something is wrong here, value of response timestamp is 0. JSON Format.")
		}

		// CSV format.
		t.Log("Testing for CSV format.")
		osReqCSV := OrdersSnapshotRequest{
			Exchange:         "binance",
			Market:           "BTCUSDT",
			At:               atTime,
			Format:           "csv",
			FileNameWithPath: "./testdata/orders_snapshot.csv",
		}
		_, err = c.GetOrdersSnapshot(osReqCSV)
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open(osReqCSV.FileNameWithPath)
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
