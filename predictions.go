package gonomics

import (
	"encoding/json"
	"strings"
	"time"
)

// Currencies Predictions Ticker.

// CurrenciesPredictionsTickerRequest represents currencies predictions ticker request parameters.
type CurrenciesPredictionsTickerRequest struct {
	Ids []string
}

// CurrenciesPredictionsTickerResponse represents currencies predictions ticker response.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesPredictionsTickerResponse struct {
	ID          string                                           `json:"id"`
	Predictions []CurrenciesPredictionsTickerPredictionsResponse `json:"predictions"`
}

// CurrenciesPredictionsTickerPredictionsResponse represents currencies predictions ticker predictions response,
// Included in CurrenciesPredictionsTickerResponse.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesPredictionsTickerPredictionsResponse struct {
	AvgErrorPct    float64   `json:"avg_error_pct,string"`
	AvgErrorPct30D float64   `json:"avg_error_pct_30d,string"`
	AvgErrorPct7D  float64   `json:"avg_error_pct_7d,string"`
	Interval       string    `json:"interval"`
	PriceChangePct float64   `json:"price_change_pct,string"`
	PriceEnd       float64   `json:"price_end,string"`
	PriceStart     float64   `json:"price_start,string"`
	TimestampEnd   time.Time `json:"timestamp_end"`
	TimestampStart time.Time `json:"timestamp_start"`
}

// GetCurrenciesPredictionsTicker fetches the currencies predictions ticker from the server
// and returns array of CurrenciesPredictionsTickerResponse.
func (c *Connecter) GetCurrenciesPredictionsTicker(cptReq CurrenciesPredictionsTickerRequest) ([]CurrenciesPredictionsTickerResponse, error) {
	req, err := c.newRequest(currenciesPredictionsTickerURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(cptReq.Ids) > 0 {
		q.Add("ids", strings.Join(cptReq.Ids[:], ","))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var cptResp []CurrenciesPredictionsTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&cptResp); err != nil {
		return nil, err
	}
	return cptResp, nil
}

// Currencies Predictions History.

// CurrenciesPredictionsHistoryRequest represents currencies predictions history request parameters.
type CurrenciesPredictionsHistoryRequest struct {
	ID       string
	Interval string
}

// CurrenciesPredictionsHistoryResponse represents currencies predictions history response.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesPredictionsHistoryResponse struct {
	ID          string                                            `json:"id"`
	Interval    string                                            `json:"interval"`
	Predictions []CurrenciesPredictionsHistoryPredictionsResponse `json:"predictions"`
}

// CurrenciesPredictionsHistoryPredictionsResponse represents currencies predictions history predictions response,
// Included in CurrenciesPredictionsHistoryResponse.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesPredictionsHistoryPredictionsResponse struct {
	PriceChangePct float64   `json:"price_change_pct,string"`
	PriceEnd       float64   `json:"price_end,string"`
	ActualPriceEnd float64   `json:"actual_price_end,string"`
	PriceStart     float64   `json:"price_start,string"`
	TimestampEnd   time.Time `json:"timestamp_end"`
	TimestampStart time.Time `json:"timestamp_start"`
}

// GetCurrenciesPredictionsHistory fetches the currencies predictions history from the server
// and returns CurrenciesPredictionsTickerResponse.
func (c *Connecter) GetCurrenciesPredictionsHistory(cphReq CurrenciesPredictionsHistoryRequest) (CurrenciesPredictionsHistoryResponse, error) {
	req, err := c.newRequest(currenciesPredictionsHistoryURL)
	if err != nil {
		return CurrenciesPredictionsHistoryResponse{}, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if cphReq.ID != "" {
		q.Add("id", cphReq.ID)
	}
	if cphReq.Interval != "" {
		q.Add("interval", cphReq.Interval)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return CurrenciesPredictionsHistoryResponse{}, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var cphResp CurrenciesPredictionsHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&cphResp); err != nil {
		return CurrenciesPredictionsHistoryResponse{}, err
	}
	return cphResp, nil
}
