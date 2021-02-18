package gonomics

import (
	"encoding/json"
	"errors"
	"time"
)

// Orders Snapshot.

// OrdersSnapshotRequest represents orders snapshot request parameters.
type OrdersSnapshotRequest struct {
	Exchange string
	Market   string
	At       time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/orders_snapshot_1613046296.csv.
	// Here, new "orders_snapshot_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// OrdersSnapshotResponse represents orders snapshot response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type OrdersSnapshotResponse struct {
	Timestamp time.Time   `json:"timestamp"`
	Bids      [][]float64 `json:"bids"`
	Asks      [][]float64 `json:"asks"`
}

// GetOrdersSnapshot fetches the orders snapshot from the server and returns
// OrdersSnapshotResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, OrdersSnapshotRequest.FileNameWithPath is required
// and the []OrdersSnapshotResponse return data is nil.
func (c *Connecter) GetOrdersSnapshot(osReq OrdersSnapshotRequest) (OrdersSnapshotResponse, error) {
	req, err := c.newRequest(ordersSnapshotURL)
	if err != nil {
		return OrdersSnapshotResponse{}, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if osReq.Exchange == "" {
		return OrdersSnapshotResponse{}, errors.New("exchange is required")
	}
	q.Add("exchange", osReq.Exchange)
	if osReq.Market == "" {
		return OrdersSnapshotResponse{}, errors.New("market is required")
	}
	q.Add("market", osReq.Market)
	if !osReq.At.IsZero() {
		q.Add("at", osReq.At.Format(time.RFC3339))
	}
	if osReq.Format != "" {
		if osReq.Format == "csv" && osReq.FileNameWithPath == "" {
			return OrdersSnapshotResponse{}, errors.New("csv file path is required")
		}
		q.Add("format", osReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return OrdersSnapshotResponse{}, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if osReq.Format == "" || osReq.Format == "json" {
		var osResp OrdersSnapshotResponse
		if err := json.NewDecoder(resp.Body).Decode(&osResp); err != nil {
			return OrdersSnapshotResponse{}, err
		}
		return osResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, osReq.FileNameWithPath)
	if err != nil {
		return OrdersSnapshotResponse{}, err
	}
	return OrdersSnapshotResponse{}, nil
}
