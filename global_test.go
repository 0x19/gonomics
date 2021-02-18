package gonomics

import (
	"testing"
	"time"
)

// TestGetGlobalTicker tests global-ticker fetching function.
func TestGetGlobalTicker(t *testing.T) {
	if plan == "free" {
		t.Skip("Skipping /global-ticker API endpoint. (Paid Plan) testing.")
	} else {
		t.Log("Testing /global-ticker API endpoint. (Paid Plan)")
		// demoAPIKey is defined in connector.go
		// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
		c := New(demoAPIKey)

		// We can modify Timeout, Transport etc of http if the default is not good.
		c.HTTPClient.Timeout = time.Second * 10

		gtReq := GlobalTickerRequest{
			Convert: "EUR",
		}
		gtResp, err := c.GetGlobalTicker(gtReq)
		if err != nil {
			t.Error(err)
		}
		if len(gtResp) < 1 {
			t.Error("Something is wrong here, length of response array is 0.")
		}
		for _, gt := range gtResp {
			if gt.Three0D.Volume == 0 {
				t.Error("Something is wrong here, value of response arrays 30d volume is 0.")
			}
		}
	}
}
