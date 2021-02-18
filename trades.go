package gonomics

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Trades.

// TradesRequest represents trades request parameters.
type TradesRequest struct {
	Exchange string
	Market   string
	Limit    int
	Order    string
	From     time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/trades_1613046296.csv.
	// Here, new "trades_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// TradesResponse represents trades response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type TradesResponse struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
	Volume    float64   `json:"volume,string"`
}

// GetTrades fetches the trades from the server and returns array of
// TradesResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, TradesRequest.FileNameWithPath is required
// and the []TradesResponse return data is nil.
func (c *Connecter) GetTrades(tReq TradesRequest) ([]TradesResponse, error) {
	req, err := c.newRequest(tradesURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if tReq.Exchange == "" {
		return nil, errors.New("exchange is required")
	}
	q.Add("exchange", tReq.Exchange)
	if tReq.Market == "" {
		return nil, errors.New("market is required")
	}
	q.Add("market", tReq.Market)
	if tReq.Limit != 0 {
		q.Add("limit", strconv.Itoa(tReq.Limit))
	}
	if tReq.Order != "" {
		q.Add("order", tReq.Order)
	}
	if !tReq.From.IsZero() {
		q.Add("from", tReq.From.Format(time.RFC3339))
	}
	if tReq.Format != "" {
		if tReq.Format == "csv" && tReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", tReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if tReq.Format == "" || tReq.Format == "json" {
		var tResp []TradesResponse
		if err := json.NewDecoder(resp.Body).Decode(&tResp); err != nil {
			return nil, err
		}
		return tResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, tReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
