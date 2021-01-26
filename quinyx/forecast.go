package quinyx

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/go-querystring/query"
)

// ForecastService handles Quinyx Forecast data
//
// Quinyx API docs: https://api.quinyx.com/v2/docs/swagger-ui.html?urls.primaryName=forecast#/
type ForecastService service

const maxRowsPerCall = 366
const maxDaysRange = 120

var (
	// ErrorReqfieldsMissing is the error returned when the request Options does not have all required fields defined
	ErrorReqfieldsMissing = fmt.Errorf("Required fields in the Options not provided, see docs")
	// ErrorDaterangeTooWide is the error returned when the requested daterange in days is to wide
	ErrorDaterangeTooWide = fmt.Errorf("The amount of days between StartTime and EndTime is above the limit, see docs")
)

// DataProviderInputList is the object used to update actual-data in Quinyx Forecast
type DataProviderInputList struct {
	DataProviderInputs []DataProviderInput `json:"requests"`
}

// PredictedDataInputList provides a list of ForecastPredictions
type PredictedDataInputList struct {
	ForecastPredictions []ForecastPrediction `json:"requests"`
}

// ForecastPrediction defines the prediction object
type ForecastPrediction struct {
	ExternalForecastVariableID      *string    `json:"externalForecastVariableId,omitempty"`
	ExternalForecastConfigurationID *string    `json:"externalForecastConfigurationId,omitempty"`
	ExternalUnitID                  *string    `json:"externalUnitId,omitempty"`
	ExternalSectionID               *string    `json:"externalSectionId,omitempty"`
	RunIdentifier                   *string    `json:"runIdentifier,omitempty"`
	RunTimestamp                    *Timestamp `json:"runTimestamp,omitempty"`
	Payloads                        []*Payload `json:"forecastDataPayload"`
}

// DataProviderInput is the input object to feed Quinyx forecast
type DataProviderInput struct {
	ExternalForecastVariableID *string    `json:"externalForecastVariableId,omitempty"`
	ExternalUnitID             *string    `json:"externalUnitId,omitempty"`
	ExternalSectionID          *string    `json:"externalSectionId,omitempty"`
	DataPayload                []*Payload `json:"forecastDataPayload,omitempty"`
}

// DataProvider is the output object from Quinyx forecast
type DataProvider struct {
	ExternalForecastVariableID *string    `json:"externalForecastVariableId,omitempty"`
	ExternalUnitID             *string    `json:"externalUnitId,omitempty"`
	ExternalSectionID          *string    `json:"externalSectionId,omitempty"`
	DataPayload                []*Payload `json:"dataPayload,omitempty"`
}

// Payload is the raw data object
type Payload struct {
	Data      *float64   `json:"data,omitempty"`
	Timestamp *Timestamp `json:"timestamp,omitempty"`
}

// AggregatedPayload is the aggregated data object
type AggregatedPayload struct {
	Data      *float64   `json:"data,omitempty"`
	EndTime   *Timestamp `json:"endTime,omitempty"`
	StartTime *Timestamp `json:"startTime,omitempty"`
}

// RequestRangeOptions is the options object for querying forecast
type RequestRangeOptions struct {
	// StartTime is required
	StartTime time.Time `url:"startTime"`
	// EndTime is required
	EndTime time.Time `url:"endTime"`
	// ExternalSectionID is optional
	ExternalSectionID *string `url:"externalSectionId,omitempty"`
	// ExternalUnitID is required
	ExternalUnitID *string `url:"externalUnitId"`
}

// RequestOptions is the options object for editing calculated forecast
type RequestOptions struct {
	// ExternalSectionID is optional
	ExternalSectionID *string `url:"externalSectionId,omitempty"`
	// ExternalUnitID is required
	ExternalUnitID *string `url:"externalUnitId"`
}

// CalculatedForecast is the calculated forecast generated in Quinyx Forecast
type CalculatedForecast struct {
	DataPayload                     []*CalculatedPayload `json:"dataPayload,omitempty"`
	ExternalForecastConfigurationID *string              `json:"externalForecastConfigurationId,omitempty"`
	ExternalSectionID               *string              `json:"externalSectionId,omitempty"`
	ExternalUnitID                  *string              `json:"externalUnitId,omitempty"`
}

// CalculatedPayload is the payload data object for calculated data
type CalculatedPayload struct {
	Data       *float64   `json:"data,omitempty"`
	EditedData *float64   `json:"editedData,omitempty"`
	StartTime  *Timestamp `json:"startTime,omitempty"`
	EndTime    *Timestamp `json:"endTime,omitempty"`
}

// EditCalculatedRequest is the object used to edit calculated forecast
type EditCalculatedRequest struct {
	RepetitionSetup        bool      `json:"repetitionSetup"`
	StartTime              Timestamp `json:"startTime"`
	EndTime                Timestamp `json:"endTime"`
	PercentageModification float64   `json:"percentageModification"`
	NewValueForPeriod      float64   `json:"newValueForPeriod"`
	WeekDays               []Weekday `json:"weekdays"`
	RepetitionEndDate      Timestamp `json:"repetitionEndDate"`
	WeekPattern            int32     `json:"weekPattern"`
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

// DynamicRule defines a dynamic rule
type DynamicRule struct {
	Amount                     int64       `json:"amount"`
	EndTime                    LocalTime   `json:"endTime"`
	StartTime                  LocalTime   `json:"startTime"`
	ExternalID                 string      `json:"externalId"`
	ExternalForecastVariableID string      `json:"forecastExternalVariableId"`
	ShiftTypes                 []ShiftType `json:"shiftTypes"`
	Weekdays                   []Weekday   `json:"weekdays"`
}

// StaticRule defines a static rule
type StaticRule struct {
	Comment      string    `json:"comment"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	StartTime    LocalTime `json:"startTime"`
	EndTime      LocalTime `json:"endTime"`
	ExternalID   string    `json:"externalId"`
	RepeatPeriod int       `json:"repeatPeriod"`
	ShiftType    ShiftType `json:"shiftType"`
	Weekdays     []Weekday `json:"weekdays"`
}

// LocalTime is a specific time on a day
type LocalTime struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Nano   int `json:"nano"`
	Second int `json:"second"`
}

// ShiftType ShiftType
type ShiftType struct {
	Amount      int    `json:"amount"`
	ShiftTypeID string `json:"externalShiftTypeId"`
}

// GetDynamicRules lists dynamic rules
func (s *ForecastService) GetDynamicRules(ctx context.Context, RequestOptions *RequestOptions) ([]*DynamicRule, *Response, error) {
	var r []*DynamicRule
	u := "forecasts/dynamic-rules"
	if !RequestOptions.hasRequiredFields() {
		return r, nil, ErrorReqfieldsMissing
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return r, nil, err
	}
	v, err := query.Values(RequestOptions)
	if err != nil {
		return r, nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return r, resp, err
	}
	return r, resp, nil
}

// GetStaticRules lists static rules
func (s *ForecastService) GetStaticRules(ctx context.Context, RequestOptions *RequestOptions) ([]*StaticRule, *Response, error) {
	var r []*StaticRule
	u := "forecasts/static-rules"
	if !RequestOptions.hasRequiredFields() {
		return r, nil, ErrorReqfieldsMissing
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return r, nil, err
	}
	v, err := query.Values(RequestOptions)
	if err != nil {
		return r, nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return r, resp, err
	}
	return r, resp, nil
}

// CreateDynamicRule creates a dynamic rule
func (s *ForecastService) CreateDynamicRule(ctx context.Context, rule *DynamicRule, RequestOptions *RequestOptions) (*DynamicRule, *Response, error) {
	var r *DynamicRule
	u := "forecasts/dynamic-rules"
	if !RequestOptions.hasRequiredFields() {
		return r, nil, ErrorReqfieldsMissing
	}
	req, err := s.client.NewRequest("POST", u, rule)
	if err != nil {
		return r, nil, err
	}

	v, err := query.Values(RequestOptions)
	if err != nil {
		return r, nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, &r)
	return r, resp, err
}

// CreateStaticRule creates a static rule
func (s *ForecastService) CreateStaticRule(ctx context.Context, rule *StaticRule, RequestOptions *RequestOptions) (*StaticRule, *Response, error) {
	var r *StaticRule
	u := "forecasts/static-rules"
	if !RequestOptions.hasRequiredFields() {
		return r, nil, ErrorReqfieldsMissing
	}
	req, err := s.client.NewRequest("POST", u, rule)
	if err != nil {
		return r, nil, err
	}

	v, err := query.Values(RequestOptions)
	if err != nil {
		return r, nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, &r)
	return r, resp, err
}

// UpdateDynamicRule updates the existing dynamic rule
func (s *ForecastService) UpdateDynamicRule(ctx context.Context, rule *DynamicRule, RequestOptions *RequestOptions) (*Response, error) {
	u := "forecasts/dynamic-rules"
	if !RequestOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}

	req, err := s.client.NewRequest("PUT", u, rule)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(RequestOptions)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// UpdateStaticRule updates the existing static rule
func (s *ForecastService) UpdateStaticRule(ctx context.Context, rule *StaticRule, RequestOptions *RequestOptions) (*Response, error) {
	u := "forecasts/static-rules"
	if !RequestOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}

	req, err := s.client.NewRequest("PUT", u, rule)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(RequestOptions)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeleteDynamicRule deletes a dynamic rule
func (s *ForecastService) DeleteDynamicRule(ctx context.Context, dynamicRuleID string, RequestOptions *RequestOptions) (*Response, error) {
	u := "forecasts/dynamic-rules"
	if !RequestOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}

	// url parameters for dynamic-rule-controller
	type params struct {
		ExternalDynamicRuleID string  `url:"externalDynamicRuleId"`
		ExternalSectionID     *string `url:"externalSectionId,omitempty"`
		ExternalUnitID        *string `url:"externalUnitId"`
	}
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(params{
		ExternalDynamicRuleID: dynamicRuleID,
		ExternalSectionID:     RequestOptions.ExternalSectionID,
		ExternalUnitID:        RequestOptions.ExternalUnitID,
	})
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeleteStaticRule deletes a static rule
func (s *ForecastService) DeleteStaticRule(ctx context.Context, staticRuleID string, RequestOptions *RequestOptions) (*Response, error) {
	u := "forecasts/static-rules"
	if !RequestOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}

	// url parameters for static-rule-controller
	type params struct {
		ExternalStaticRuleID string  `url:"externalStaticRuleId"`
		ExternalSectionID    *string `url:"externalSectionId,omitempty"`
		ExternalUnitID       *string `url:"externalUnitId"`
	}
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(params{
		ExternalStaticRuleID: staticRuleID,
		ExternalSectionID:    RequestOptions.ExternalSectionID,
		ExternalUnitID:       RequestOptions.ExternalUnitID,
	})
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// UploadActualData sends raw datapoints to Quinyx Forecast API
func (s *ForecastService) UploadActualData(ctx context.Context, appendData bool, dil *DataProviderInputList) (*Response, error) {
	u := fmt.Sprintf("forecasts/actual-data?appendData=%v", appendData)
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
	return resp, err
}

// UploadBudgetData sends budget datapoints to Quinyx Forecast API
func (s *ForecastService) UploadBudgetData(ctx context.Context, appendData bool, dil *DataProviderInputList) (*Response, error) {
	u := fmt.Sprintf("forecasts/budget-data?appendData=%v", appendData)
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
	return resp, err
}

func (g *RequestRangeOptions) hasRequiredFields() bool {
	if g == nil {
		return false
	}
	if g.EndTime.IsZero() || g.ExternalUnitID == nil || g.StartTime.IsZero() {
		return false
	}
	return true
}

func (g *RequestRangeOptions) dayDistance() float64 {
	return math.Floor(g.EndTime.Sub(g.StartTime).Hours() / 24)
}

func (g *RequestOptions) hasRequiredFields() bool {
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
func (s *ForecastService) GetActualData(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) ([]*DataProvider, *Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/actual-data", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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

// DeleteActualData deletes actual data.
// Deleting the actual data previously uploaded for the given forecast variable.
// This operation will also delete the corresponding calculated forecast data.
// The startTime and endTime must be at the start of hour and the range between these two dates can not exceed 120 days.
func (s *ForecastService) DeleteActualData(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) (*Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/actual-data", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
func (s *ForecastService) GetActualDataStream(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) ([]*DataProvider, *Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/actual-data-stream", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
func (s *ForecastService) GetAggregatedData(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) ([]*AggregatedPayload, *Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/aggregated-data", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
func (s *ForecastService) GetCalculatedForecast(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) ([]*CalculatedForecast, *Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/calculated-forecast", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
func (s *ForecastService) EditCalculatedForecast(ctx context.Context, externalForecastVariableID string, externalForecastConfigurationID string, RequestRangeOptions *RequestOptions, modrequest *EditCalculatedRequest) (*Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/forecast-configurations/%v/edit-forecast", externalForecastVariableID, externalForecastConfigurationID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}
	req, err := s.client.NewRequest("POST", u, modrequest)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
func (s *ForecastService) GetForecastData(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) ([]*DataProvider, *Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/forecast-data", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
func (s *ForecastService) DeleteForecastData(ctx context.Context, externalForecastVariableID string, RequestRangeOptions *RequestRangeOptions) (*Response, error) {
	u := fmt.Sprintf("forecasts/forecast-variables/%v/forecast-data", externalForecastVariableID)
	if !RequestRangeOptions.hasRequiredFields() {
		return nil, ErrorReqfieldsMissing
	}
	if RequestRangeOptions.dayDistance() > maxDaysRange {
		return nil, ErrorDaterangeTooWide
	}
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	v, err := query.Values(RequestRangeOptions)
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
	u := fmt.Sprintf("forecasts/predicted-data")
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
