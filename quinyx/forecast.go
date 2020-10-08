package quinyx

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
)

// ForecastService handles Quinyx Forecast data
//
// Quinyx API docs: https://api.quinyx.com/v2/docs/swagger-ui.html?urls.primaryName=forecast#/
type ForecastService service

const maxRowsPerCall = 366

// DataProviderInputList is the object used to update actual-data in Quinyx Forecast
type DataProviderInputList struct {
	DataProviderInputs []DataProvider `json:"requests"`
}

// DataProvider is the input object to feed Quinyx forecast
type DataProvider struct {
	ExternalForecastVariableID string    `json:"externalForecastVariableId,omitempty"`
	ExternalUnitID             string    `json:"externalUnitId,omitempty"`
	ExternalSectionID          string    `json:"externalSectionId,omitempty"`
	DataPayload                []Payload `json:"forecastDataPayload,omitempty"`
}

// Payload is the raw data object
type Payload struct {
	Data      float64   `json:"data,omitempty"`
	Timestamp Timestamp `json:"timestamp,omitempty"`
}

// RequestOptions is the options object for querying forecast
type RequestOptions struct {
	StartTime         time.Time `url:"startTime,omitempty"`
	EndTime           time.Time `url:"endTime,omitempty"`
	ExternalSectionID *string   `url:"externalSectionId,omitempty"`
	ExternalUnitID    *string   `url:"externalUnitId,omitempty"`
}

// UploadActualData sends raw datapoints to Quinyx Forecast API
func (s *ForecastService) UploadActualData(ctx context.Context, appendData bool, dil *DataProviderInputList) (*Response, error) {
	u := fmt.Sprintf("/forecasts/actual-data?appendData=%v", appendData)
	if dil != nil {
		if len(dil.DataProviderInputs) > maxRowsPerCall {
			return nil, fmt.Errorf("The total amount of data rows must not exceed 366 in a single call")
		}
	}
	req, err := s.client.NewRequest("POST", u, dil)
	if err != nil {
		return nil, err
	}
	var tagres *Tag
	resp, err := s.client.Do(ctx, req, &tagres)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// UploadBudgetData sends budget datapoints to Quinyx Forecast API
func (s *ForecastService) UploadBudgetData(ctx context.Context, appendData bool, dil *DataProviderInputList) (*Response, error) {
	u := fmt.Sprintf("/forecasts/budget-data?appendData=%v", appendData)
	if dil != nil {
		if len(dil.DataProviderInputs) > maxRowsPerCall {
			return nil, fmt.Errorf("The total amount of data rows must not exceed 366 in a single call")
		}
	}
	req, err := s.client.NewRequest("POST", u, dil)
	if err != nil {
		return nil, err
	}
	var tagres *Tag
	resp, err := s.client.Do(ctx, req, &tagres)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (g *RequestOptions) hasRequiredFields() bool {
	if g == nil {
		return false
	}
	if g.EndTime.IsZero() || g.ExternalUnitID == nil || g.StartTime.IsZero() {
		return false
	}
	return true
}

// GetActualData gets the actual data previously uploaded for the given forecast variable.
// The range between these two dates can not exceed 120 days.
func (s *ForecastService) GetActualData(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) ([]*DataProvider, *Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/actual-data", externalForecastVariableID)
	if !requestoptions.hasRequiredFields() {
		return nil, nil, fmt.Errorf("Required fields in the Options not provided, see docs")
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	v, err := query.Values(requestoptions)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = v.Encode()
	var dp []*DataProvider
	resp, err := s.client.Do(ctx, req, &dp)
	if err != nil {
		return nil, resp, err
	}
	return dp, resp, nil
}

// DeleteActualData deletes actual data
// Deleting the actual data previously uploaded for the given forecast variable.
// This operation will also delete the corresponding calculated forecast data.
// The startTime and endTime must be at the start of hour and the range between these two dates can not exceed 120 days.
func (s *ForecastService) DeleteActualData(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) (*Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/actual-data", externalForecastVariableID)
	if !requestoptions.hasRequiredFields() {
		return nil, fmt.Errorf("Required fields in the Options not provided, see docs")
	}
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(requestoptions)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
