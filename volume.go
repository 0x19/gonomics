package gonomics

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Volume History.

// VolumeHistoryRequest represents volume history request parameters.
type VolumeHistoryRequest struct {
	Start   time.Time
	End     time.Time
	Convert string

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/volume_history_1613046296.csv.
	// Here, new "volume_history_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string

	IncludeTransparency bool
}

// VolumeHistoryResponse represents volume history response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type VolumeHistoryResponse struct {
	Timestamp         time.Time `json:"timestamp"`
	Volume            float64   `json:"volume,string"`
	TransparentVolume float64   `json:"transparent_volume,string"`
}

// GetVolumeHistory fetches the volume history from the server and returns array of
// VolumeHistoryResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, VolumeHistoryRequest.FileNameWithPath is required
// and the []VolumeHistoryResponse return data is nil.
func (c *Connecter) GetVolumeHistory(vhReq VolumeHistoryRequest) ([]VolumeHistoryResponse, error) {
	req, err := c.newRequest(volumeHistoryURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if !vhReq.Start.IsZero() {
		q.Add("start", vhReq.Start.Format(time.RFC3339))
	}
	if !vhReq.End.IsZero() {
		q.Add("end", vhReq.End.Format(time.RFC3339))
	}
	if vhReq.Convert != "" {
		q.Add("convert", vhReq.Convert)
	}
	if vhReq.Format != "" {
		if vhReq.Format == "csv" && vhReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", vhReq.Format)
	}
	if vhReq.IncludeTransparency {
		q.Add("include-transparency", strconv.FormatBool(vhReq.IncludeTransparency))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if vhReq.Format == "" || vhReq.Format == "json" {
		var vhResp []VolumeHistoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&vhResp); err != nil {
			return nil, err
		}
		return vhResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, vhReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
