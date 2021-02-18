package gonomics

import (
	"testing"
	"time"
)

// TestGetCurrenciesPredictionsTicker tests currencies predictions ticker fetching function.
func TestGetCurrenciesPredictionsTicker(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /currencies/predictions/ticker API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /currencies/predictions/ticker API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		cptReq := CurrenciesPredictionsTickerRequest{
			Ids: []string{"BTC"},
		}
		cptResp, err := c.GetCurrenciesPredictionsTicker(cptReq)
		if err != nil {
			t.Error(err)
		}
		if len(cptResp) < 1 {
			t.Error("Something is wrong here, length of response array is 0.")
		}
		for _, cpt := range cptResp {
			if cpt.ID == "" {
				t.Error("Something is wrong here, value of response arrays id is empty.")
			}
		}
	}
}

// TestGetCurrenciesPredictionsHistory tests currencies predictions history fetching function.
func TestGetCurrenciesPredictionsHistory(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /currencies/predictions/history API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /currencies/predictions/history API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		cphReq := CurrenciesPredictionsHistoryRequest{
			ID:       "BTC",
			Interval: "7d",
		}
		cphResp, err := c.GetCurrenciesPredictionsHistory(cphReq)
		if err != nil {
			t.Error(err)
		}
		if cphResp.ID == "" {
			t.Error("Something is wrong here, value of response id is empty.")
		}
	}
}
