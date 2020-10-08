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

// PredictedDataInputList PredictedDataInputList
type PredictedDataInputList struct {
	ForecastPredictions []ForecastPrediction `json:"requests"`
}

// ForecastPrediction ForecastPrediction
type ForecastPrediction struct {
	ExternalForecastVariableID      string    `json:"externalForecastVariableId,omitempty"`
	ExternalForecastConfigurationID string    `json:"externalForecastConfigurationId,omitempty"`
	ExternalUnitID                  string    `json:"externalUnitId,omitempty"`
	ExternalSectionID               string    `json:"externalSectionId,omitempty"`
	RunIdentifier                   string    `json:"runIdentifier,omitempty"`
	RunTimestamp                    Timestamp `json:"runTimestamp,omitempty"`
	Payloads                        []Payload `json:"forecastDataPayload"`
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

// AggregatedPayload is the aggregated data object
type AggregatedPayload struct {
	Data      float64   `json:"data,omitempty"`
	EndTime   Timestamp `json:"endTime,omitempty"`
	StartTime Timestamp `json:"startTime,omitempty"`
}

// RequestOptions is the options object for querying forecast
type RequestOptions struct {
	StartTime         time.Time `url:"startTime,omitempty"`
	EndTime           time.Time `url:"endTime,omitempty"`
	ExternalSectionID *string   `url:"externalSectionId,omitempty"`
	ExternalUnitID    *string   `url:"externalUnitId,omitempty"`
}

// EditCalculatedOptions is the options object for editing calculated forecast
type EditCalculatedOptions struct {
	ExternalSectionID *string `url:"externalSectionId,omitempty"`
	ExternalUnitID    *string `url:"externalUnitId,omitempty"`
}

// CalculatedForecast is the calculated forecast generated in Quinyx Forecast
type CalculatedForecast struct {
	DataPayload                     []CalculatedPayload `json:"dataPayload,omitempty"`
	ExternalForecastConfigurationID *string             `json:"externalForecastConfigurationId,omitempty"`
	ExternalSectionID               *string             `json:"externalSectionId,omitempty"`
	ExternalUnitID                  *string             `json:"externalUnitId,omitempty"`
}

// CalculatedPayload is the payload data object for calculated data
type CalculatedPayload struct {
	Data       float64   `json:"data,omitempty"`
	EditedData float64   `json:"editedData,omitempty"`
	StartTime  Timestamp `json:"startTime,omitempty"`
	EndTime    Timestamp `json:"endTime,omitempty"`
}

// EditCalculatedRequest is the object used to edit calculated forecast
type EditCalculatedRequest struct {
	RepetitionSetup        bool      `json:"repetitionSetup,omitempty"`
	StartTime              Timestamp `json:"startTime,omitempty"`
	EndTime                Timestamp `json:"endTime,omitempty"`
	PercentageModification float64   `json:"percentageModification,omitempty"`
	NewValueForPeriod      float64   `json:"newValueForPeriod,omitempty"`
	WeekDays               []Weekday `json:"weekdays,omitempty"`
	RepetitionEndDate      Timestamp `json:"repetitionEndDate,omitempty"`
	WeekPattern            int32     `json:"weekPattern,omitempty"`
}

// Weekday defines a weekday
type Weekday string

// TagTypes
const (
	Monday    Weekday = "0"
	Tuesday           = "1"
	Wednesday         = "2"
	Thursday          = "3"
	Friday            = "4"
	Saturday          = "5"
	Sunday            = "6"
)

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

func (g *EditCalculatedOptions) hasRequiredFields() bool {
	if g == nil {
		return false
	}
	if g.ExternalUnitID == nil {
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

// GetActualDataStream gets the actual data previously uploaded for the given forecast variable.
// The range between these two dates can not exceed 120 days.
func (s *ForecastService) GetActualDataStream(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) ([]*DataProvider, *Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/actual-data-stream", externalForecastVariableID)
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

// GetAggregatedData gets the aggregated data for the given forecast variable.
// The range between these two dates can not exceed 120 days.
func (s *ForecastService) GetAggregatedData(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) ([]*AggregatedPayload, *Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/aggregated-data", externalForecastVariableID)
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
	var dp []*AggregatedPayload
	resp, err := s.client.Do(ctx, req, &dp)
	if err != nil {
		return nil, resp, err
	}
	return dp, resp, nil
}

// GetCalculatedForecast gets the calculated forecast for the given forecast variable.
// The range between these two dates can not exceed 120 days.
func (s *ForecastService) GetCalculatedForecast(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) ([]*CalculatedForecast, *Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/calculated-forecast", externalForecastVariableID)
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
	var cf []*CalculatedForecast
	resp, err := s.client.Do(ctx, req, &cf)
	if err != nil {
		return nil, resp, err
	}
	return cf, resp, nil
}

// EditCalculatedForecast changes the calculated forecast.
func (s *ForecastService) EditCalculatedForecast(ctx context.Context, externalForecastVariableID string, externalForecastConfigurationID string, requestoptions *EditCalculatedOptions, modrequest *EditCalculatedRequest) (*Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/forecast-configurations/%v/edit-forecast", externalForecastVariableID, externalForecastConfigurationID)
	if !requestoptions.hasRequiredFields() {
		return nil, fmt.Errorf("Required fields in the Options not provided, see docs")
	}
	req, err := s.client.NewRequest("POST", u, modrequest)
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

// GetForecastData gets the uploaded forecast data for the given forecast variable.
// The range between these two dates can not exceed 120 days.
func (s *ForecastService) GetForecastData(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) ([]*DataProvider, *Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/forecast-data", externalForecastVariableID)
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

// DeleteForecastData deletes the previously uploaded forecast data for the the given forecast variable.
// The startTime and endTime must be at the start of hour and the range between these two dates can not exceed 120 days.
func (s *ForecastService) DeleteForecastData(ctx context.Context, externalForecastVariableID string, requestoptions *RequestOptions) (*Response, error) {
	u := fmt.Sprintf("/forecasts/forecast-variables/%v/forecast-data", externalForecastVariableID)
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

// UploadPredictedData is the Operation used to upload generated prediction data.
// Resolution of datapoints must match expected resolution of variable.
// The total amount of data rows must not exceed 366
func (s *ForecastService) UploadPredictedData(ctx context.Context, inlist *PredictedDataInputList) (*Response, error) {
	u := fmt.Sprintf("/forecasts/predicted-data")
	if inlist != nil {
		if len(inlist.ForecastPredictions) > maxRowsPerCall {
			return nil, fmt.Errorf("The total amount of data rows must not exceed 366 in a single call")
		}
	}
	req, err := s.client.NewRequest("POST", u, inlist)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
