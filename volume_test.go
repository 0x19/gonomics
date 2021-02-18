package gonomics

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestGetVolumeHistory tests volume history fetching function.
func TestGetVolumeHistory(t *testing.T) {
	t.Log("Testing /volume/history API endpoint. (Partial Paid Plan)")
	// demoAPIKey is defined in connector.go
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	c := New(demoAPIKey)

	// We can modify Timeout, Transport etc of http if the default is not good.
	c.HTTPClient.Timeout = time.Second * 10

	// Json format.
	t.Log("Testing for JSON format.")
	startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
	vhReqJSON := VolumeHistoryRequest{
		Start:   startTime,
		End:     endTime,
		Convert: "EUR",
	}
	vhRespJSON, err := c.GetVolumeHistory(vhReqJSON)
	if err != nil {
		t.Error(err)
	}
	if len(vhRespJSON) < 1 {
		t.Error("Something is wrong here, length of response array is 0. JSON Format.")
	}
	for _, vh := range vhRespJSON {
		if vh.Volume == 0 {
			t.Error("Something is wrong here, value of response arrays volume is 0. JSON Format.")
		}
	}

	// CSV format.
	t.Log("Testing for CSV format.")
	vhReqCSV := VolumeHistoryRequest{
		Start:            startTime,
		End:              endTime,
		Convert:          "EUR",
		Format:           "csv",
		FileNameWithPath: "./testdata/volume_history.csv",
	}
	_, err = c.GetVolumeHistory(vhReqCSV)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(vhReqCSV.FileNameWithPath)
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
