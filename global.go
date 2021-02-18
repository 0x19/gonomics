package gonomics

import (
	"encoding/json"
)

// Global Ticker.

// GlobalTickerRequest represents global-ticker requests parameters.
type GlobalTickerRequest struct {
	Convert string
}

// GlobalTickerResponse represents global-ticker response.
// Fields will contain default go lang values if there is no value received from the server.
type GlobalTickerResponse struct {
	NumCurrencies        int                          `json:"num_currencies,string"`
	NumCurrenciesActive  int                          `json:"num_currencies_active,string"`
	NumCurrenciesInative int                          `json:"num_currencies_inative,string"`
	NumCurrenciesDead    int                          `json:"num_currencies_dead,string"`
	NumCurrenciesNew     int                          `json:"num_currencies_new,string"`
	MarketCap            float64                      `json:"market_cap,string"`
	TransparentMarketCap float64                      `json:"transparent_market_cap,string"`
	OneD                 GlobalTickerIntervalResponse `json:"1d"`
	SevenD               GlobalTickerIntervalResponse `json:"7d"`
	Three0D              GlobalTickerIntervalResponse `json:"30d"`
	Three65D             GlobalTickerIntervalResponse `json:"365d"`
	Ytd                  GlobalTickerIntervalResponse `json:"ytd"`
}

// GlobalTickerIntervalResponse represents global-ticker interval response,
// Included in GlobalTickerResponse.
// Fields will contain default go lang values if there is no value received from the server.
type GlobalTickerIntervalResponse struct {
	MarketCapChange               float64                                          `json:"market_cap_change,string"`
	MarketCapChangePct            float64                                          `json:"market_cap_change_pct,string"`
	TransparentMarketCapChange    float64                                          `json:"transparent_market_cap_change,string"`
	TransparentMarketCapChangePct float64                                          `json:"transparent_market_cap_change_pct,string"`
	Volume                        float64                                          `json:"volume,string"`
	VolumeChange                  float64                                          `json:"volume_change,string"`
	VolumeChangePct               float64                                          `json:"volume_change_pct,string"`
	TransparentVolume             float64                                          `json:"transparent_volume,string"`
	TransparentVolumeChange       float64                                          `json:"transparent_volume_change,string"`
	TransparentVolumeChangePct    float64                                          `json:"transparent_volume_change_pct,string"`
	VolumeTransparency            []GlobalTickerIntervalVolumeTransparencyResponse `json:"volume_transparency"`
}

// GlobalTickerIntervalVolumeTransparencyResponse represents global-ticker interval volume transaparency response,
// Included in GlobalTickerIntervalResponse.
// Fields will contain default go lang values if there is no value received from the server.
type GlobalTickerIntervalVolumeTransparencyResponse struct {
	Grade           string  `json:"grade"`
	Volume          float64 `json:"volume,string"`
	VolumeChange    float64 `json:"volume_change,string"`
	VolumeChangePct float64 `json:"volume_change_pct,string"`
}

// GetGlobalTicker fetches the global-ticker from the server
// and returns array of GlobalTickerResponse.
func (c *Connecter) GetGlobalTicker(gtReq GlobalTickerRequest) ([]GlobalTickerResponse, error) {
	req, err := c.newRequest(globalTickerURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if gtReq.Convert != "" {
		q.Add("convert", gtReq.Convert)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var gtResp []GlobalTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&gtResp); err != nil {
		return nil, err
	}
	return gtResp, nil
}
