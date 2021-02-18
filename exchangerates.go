package gonomics

import (
	"encoding/json"
	"errors"
	"time"
)

// Exchange Rates.

// ExchangeRatesRequest represents exchange rates request parameters.
type ExchangeRatesRequest struct {
	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/exchange_rates_1613046296.csv.
	// Here, new "exchange_rates_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// ExchangeRatesResponse represents exchange rates response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangeRatesResponse struct {
	Currency  string    `json:"currency"`
	Rate      float64   `json:"rate,string"`
	Timestamp time.Time `json:"timestamp"`
}

// GetExchangeRates fetches the exchange rates from the server and returns array of
// ExchangeRatesResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, ExchangeRatesRequest.FileNameWithPath is required
// and the []ExchangeRatesResponse return data is nil.
func (c *Connecter) GetExchangeRates(erReq ExchangeRatesRequest) ([]ExchangeRatesResponse, error) {
	req, err := c.newRequest(exchangeRatesURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if erReq.Format != "" {
		if erReq.Format == "csv" && erReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", erReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if erReq.Format == "" || erReq.Format == "json" {
		var erResp []ExchangeRatesResponse
		if err := json.NewDecoder(resp.Body).Decode(&erResp); err != nil {
			return nil, err
		}
		return erResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, erReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Exchange Rates History.

// ExchangeRatesHistoryRequest represents exchange-rates history request parameters.
type ExchangeRatesHistoryRequest struct {
	Currency string
	Start    time.Time
	End      time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/exchange-rates_history_1613046296.csv.
	// Here, new "exchange-rates_history_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// ExchangeRatesHistoryResponse represents exchange-rates history response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangeRatesHistoryResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Rate      float64   `json:"rate,string"`
}

// GetExchangeRatesHistory fetches the exchange-rates history from the server and returns array of
// ExchangeRatesHistoryResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, ExchangeRatesHistoryRequest.FileNameWithPath is required
// and the []ExchangeRatesHistoryResponse return data is nil.
func (c *Connecter) GetExchangeRatesHistory(erhReq ExchangeRatesHistoryRequest) ([]ExchangeRatesHistoryResponse, error) {
	req, err := c.newRequest(exchangeRatesHistoryURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if erhReq.Currency == "" {
		return nil, errors.New("currency is required")
	}
	q.Add("currency", erhReq.Currency)
	if erhReq.Start.IsZero() {
		return nil, errors.New("start is required")
	}
	q.Add("start", erhReq.Start.Format(time.RFC3339))
	if !erhReq.End.IsZero() {
		q.Add("end", erhReq.End.Format(time.RFC3339))
	}
	if erhReq.Format != "" {
		if erhReq.Format == "csv" && erhReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", erhReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if erhReq.Format == "" || erhReq.Format == "json" {
		var erhResp []ExchangeRatesHistoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&erhResp); err != nil {
			return nil, err
		}
		return erhResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, erhReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
