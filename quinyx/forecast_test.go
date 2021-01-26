package quinyx

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestUploadActualData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	want := &DataProviderInputList{
		DataProviderInputs: []DataProviderInput{
			{
				ExternalForecastVariableID: String("a"),
				ExternalSectionID:          String("b"),
				ExternalUnitID:             String("c"),
				DataPayload: []*Payload{
					{
						Data:      Float64(123),
						Timestamp: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
					},
				},
			},
		},
	}
	mux.HandleFunc("/forecasts/actual-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		appendData := r.URL.Query().Get("appendData")
		assert.Equal(t, "false", appendData)

		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, `{"requests":[{"externalForecastVariableId":"a","externalUnitId":"c","externalSectionId":"b","forecastDataPayload":[{"data":123,"timestamp":"2019-10-12T07:20:50.52Z"}]}]}
`, string(body))
		assert.NilError(t, err)
	})

	_, err := client.Forecast.UploadActualData(context.Background(), false, want)
	assert.NilError(t, err)
}

func TestUploadBudgetData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	want := &DataProviderInputList{
		DataProviderInputs: []DataProviderInput{
			{
				ExternalForecastVariableID: String("a"),
				ExternalSectionID:          String("b"),
				ExternalUnitID:             String("c"),
				DataPayload: []*Payload{
					{
						Data:      Float64(123),
						Timestamp: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
					},
				},
			},
		},
	}
	mux.HandleFunc("/forecasts/budget-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		appendData := r.URL.Query().Get("appendData")
		assert.Equal(t, "false", appendData)

		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, `{"requests":[{"externalForecastVariableId":"a","externalUnitId":"c","externalSectionId":"b","forecastDataPayload":[{"data":123,"timestamp":"2019-10-12T07:20:50.52Z"}]}]}
`, string(body))
		assert.NilError(t, err)
	})

	_, err := client.Forecast.UploadBudgetData(context.Background(), false, want)
	assert.NilError(t, err)
}

func TestGetActualData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}

	want := []*DataProvider{
		{
			ExternalForecastVariableID: String("b"),
			ExternalSectionID:          String("c"),
			ExternalUnitID:             String("d"),
			DataPayload: []*Payload{
				{
					Data:      Float64(123),
					Timestamp: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
				},
			},
		},
	}
	mux.HandleFunc("/forecasts/forecast-variables/a/actual-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "GET")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetActualData(context.Background(), "a", options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetActualData(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/actual-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetActualData(context.Background(), "a", options)
}

func TestGetActualDataStream(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}

	want := []*DataProvider{
		{
			ExternalForecastVariableID: String("b"),
			ExternalSectionID:          String("c"),
			ExternalUnitID:             String("d"),
			DataPayload: []*Payload{
				{
					Data:      Float64(123),
					Timestamp: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
				},
			},
		},
	}
	mux.HandleFunc("/forecasts/forecast-variables/a/actual-data-stream", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "GET")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetActualDataStream(context.Background(), "a", options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetActualDataStream(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/actual-data-stream", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetActualData(context.Background(), "a", options)
}

func TestGetAggregatedData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}

	want := []*AggregatedPayload{
		{
			Data:      Float64(1234),
			StartTime: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
			EndTime:   &Timestamp{time.Date(2019, time.October, 12, 07, 21, 50, 520000000, time.UTC)},
		},
	}
	mux.HandleFunc("/forecasts/forecast-variables/a/aggregated-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "GET")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetAggregatedData(context.Background(), "a", options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetAggregatedData(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/aggregated-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetActualData(context.Background(), "a", options)
}

func TestGetCalculatedForecast(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}

	want := []*CalculatedForecast{
		{
			DataPayload: []*CalculatedPayload{
				{
					Data:       Float64(1234),
					EditedData: Float64(4321),
					StartTime:  &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
					EndTime:    &Timestamp{time.Date(2019, time.October, 12, 07, 21, 50, 520000000, time.UTC)},
				},
			},
			ExternalForecastConfigurationID: String("b"),
			ExternalSectionID:               String("c"),
			ExternalUnitID:                  String("d"),
		},
	}
	mux.HandleFunc("/forecasts/forecast-variables/a/calculated-forecast", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "GET")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetCalculatedForecast(context.Background(), "a", options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetCalculatedForecast(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/calculated-forecast", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetActualData(context.Background(), "a", options)
}

func TestEditCalculatedForecast(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := EditCalculatedRequest{
		RepetitionSetup:        true,
		StartTime:              Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		EndTime:                Timestamp{time.Date(2019, time.October, 12, 07, 21, 50, 520000000, time.UTC)},
		PercentageModification: 1234,
		NewValueForPeriod:      4321,
		WeekDays:               []Weekday{Wednesday},
		RepetitionEndDate:      Timestamp{time.Date(2019, time.October, 12, 07, 22, 50, 520000000, time.UTC)},
		WeekPattern:            1,
	}
	mux.HandleFunc("/forecasts/forecast-variables/a/forecast-configurations/b/edit-forecast", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		assert.Equal(t, `{"repetitionSetup":true,"startTime":"2019-10-12T07:20:50.52Z","endTime":"2019-10-12T07:21:50.52Z","percentageModification":1234,"newValueForPeriod":4321,"weekdays":["2"],"repetitionEndDate":"2019-10-12T07:22:50.52Z","weekPattern":1}`, string(json))
	})

	_, err := client.Forecast.EditCalculatedForecast(context.Background(), "a", "b", options, &want)
	assert.NilError(t, err)
}

func TestGetForecastData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}
	want := []*DataProvider{
		{
			ExternalForecastVariableID: String("a"),
			ExternalSectionID:          String("b"),
			ExternalUnitID:             String("c"),
			DataPayload: []*Payload{
				{
					Data:      Float64(1234.4),
					Timestamp: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
				},
			},
		},
	}

	mux.HandleFunc("/forecasts/forecast-variables/a/forecast-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "GET")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetForecastData(context.Background(), "a", options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetForecastData(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/forecast-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetActualData(context.Background(), "a", options)
}

func TestDeleteForecast(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}

	mux.HandleFunc("/forecasts/forecast-variables/a/forecast-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "DELETE")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
	})

	_, err := client.Forecast.DeleteForecastData(context.Background(), "a", options)
	assert.NilError(t, err)
	_, err = client.Forecast.DeleteForecastData(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/forecast-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetActualData(context.Background(), "a", options)
}

func TestDeleteActualData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 13, 07, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}

	mux.HandleFunc("/forecasts/forecast-variables/a/actual-data", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		assert.Equal(t, "2019-10-13T07:20:50Z", r.URL.Query().Get("endTime"))
		assert.Equal(t, "2019-10-12T07:20:50Z", r.URL.Query().Get("startTime"))
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
	})

	_, err := client.Forecast.DeleteActualData(context.Background(), "a", options)
	assert.NilError(t, err)
	_, err = client.Forecast.DeleteActualData(context.Background(), "a", &RequestRangeOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/actual-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.DeleteActualData(context.Background(), "a", options)
}

func TestUploadPredictedData(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	want := &PredictedDataInputList{
		ForecastPredictions: []ForecastPrediction{
			{
				ExternalForecastVariableID:      String("a"),
				ExternalForecastConfigurationID: String("d"),
				ExternalSectionID:               String("b"),
				ExternalUnitID:                  String("c"),
				RunIdentifier:                   String("e"),
				RunTimestamp:                    &Timestamp{time.Date(2020, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
				Payloads: []*Payload{
					{
						Data:      Float64(123),
						Timestamp: &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
					},
				},
			},
		},
	}
	mux.HandleFunc("/forecasts/predicted-data", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, `{"requests":[{"externalForecastVariableId":"a","externalForecastConfigurationId":"d","externalUnitId":"c","externalSectionId":"b","runIdentifier":"e","runTimestamp":"2020-10-12T07:20:50.52Z","forecastDataPayload":[{"data":123,"timestamp":"2019-10-12T07:20:50.52Z"}]}]}
`, string(body))
		assert.NilError(t, err)
	})

	_, err := client.Forecast.UploadPredictedData(context.Background(), want)
	assert.NilError(t, err)
}

func TestDateRangeCheck(t *testing.T) {
	options := &RequestRangeOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 17, 12, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}
	assert.Equal(t, float64(5), options.dayDistance())
	options = &RequestRangeOptions{
		StartTime:      time.Date(2019, time.July, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.December, 17, 12, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}
	assert.Assert(t, options.dayDistance() > 120)
}

func TestGetDynamicRules(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := []*DynamicRule{
		{
			Amount: 123,
			StartTime: LocalTime{
				Hour:   17,
				Minute: 30,
				Nano:   1234,
				Second: 13,
			},
			EndTime: LocalTime{
				Hour:   18,
				Minute: 31,
				Nano:   12345,
				Second: 14,
			},
			ExternalID:                 "e",
			ExternalForecastVariableID: "efvid",
			ShiftTypes: []ShiftType{
				{
					Amount:      321,
					ShiftTypeID: "abc",
				},
			},
			Weekdays: []Weekday{Monday},
		},
	}
	mux.HandleFunc("/forecasts/dynamic-rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetDynamicRules(context.Background(), options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetDynamicRules(context.Background(), &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/actual-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetDynamicRules(context.Background(), options)
}

func TestGetStaticRules(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := []*StaticRule{
		{
			Comment:   "hello there",
			StartDate: time.Date(2019, time.July, 12, 00, 00, 00, 0, time.UTC),
			EndDate:   time.Date(2019, time.July, 13, 10, 00, 00, 0, time.UTC),
			StartTime: LocalTime{
				Hour:   17,
				Minute: 30,
				Nano:   1234,
				Second: 13,
			},
			EndTime: LocalTime{
				Hour:   18,
				Minute: 31,
				Nano:   12345,
				Second: 14,
			},
			ExternalID: "e",
			ShiftType: ShiftType{

				Amount:      321,
				ShiftTypeID: "abc",
			},
			Weekdays: []Weekday{Monday},
		},
	}
	mux.HandleFunc("/forecasts/static-rules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	dps, _, err := client.Forecast.GetStaticRules(context.Background(), options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
	_, _, err = client.Forecast.GetStaticRules(context.Background(), &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	mux.HandleFunc("/forecasts/forecast-variables/b/actual-data", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "foo", r.URL.Query().Get("externalSectionID"))
	})
	options.ExternalSectionID = String("foo")
	client.Forecast.GetStaticRules(context.Background(), options)
}

func TestGetDynamicRulesMatchingExampleJSON(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := []*DynamicRule{
		{
			Amount: 123,
			StartTime: LocalTime{
				Hour:   17,
				Minute: 30,
				Nano:   1234,
				Second: 13,
			},
			EndTime: LocalTime{
				Hour:   18,
				Minute: 31,
				Nano:   12345,
				Second: 14,
			},
			ExternalID:                 "e",
			ExternalForecastVariableID: "efvid",
			ShiftTypes: []ShiftType{
				{
					Amount:      321,
					ShiftTypeID: "abc",
				},
			},
			Weekdays: []Weekday{Monday},
		},
	}

	mux.HandleFunc("/forecasts/dynamic-rules", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(`[
			{
			  "amount": 123,
			  "endTime": {
				"hour": 18,
				"minute": 31,
				"nano": 12345,
				"second": 14
			  },
			  "externalId": "e",
			  "forecastExternalVariableId": "efvid",
			  "shiftTypes": [
				{
				  "amount": 321,
				  "externalShiftTypeId": "abc"
				}
			  ],
			  "startTime": {
				"hour": 17,
				"minute": 30,
				"nano": 1234,
				"second": 13
			  },
			  "weekdays": [
				"0"
			  ]
			}
		  ]`))
	})
	dps, _, err := client.Forecast.GetDynamicRules(context.Background(), options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
}

func TestGetStaticRulesMatchingExampleJSON(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := []*StaticRule{
		{
			Comment:   "hello there",
			StartDate: time.Date(2019, time.July, 12, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2019, time.July, 13, 0, 0, 0, 0, time.UTC),
			StartTime: LocalTime{
				Hour:   17,
				Minute: 30,
				Nano:   1234,
				Second: 13,
			},
			EndTime: LocalTime{
				Hour:   18,
				Minute: 31,
				Nano:   12345,
				Second: 14,
			},
			ExternalID: "e",
			ShiftType: ShiftType{

				Amount:      321,
				ShiftTypeID: "abc",
			},
			Weekdays: []Weekday{Monday},
		},
	}

	mux.HandleFunc("/forecasts/static-rules", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(`[
			{
			  "comment": "hello there",
			  "startDate": "2019-07-12T00:00:00Z",
			  "endDate": "2019-07-13T00:00:00Z",
			  "endTime": {
				"hour": 18,
				"minute": 31,
				"nano": 12345,
				"second": 14
			  },
			  "externalId": "e",
			  "forecastExternalVariableId": "efvid",
			  "shiftType": {
				"amount": 321,
				"externalShiftTypeId": "abc"
			  },
			  "startTime": {
				"hour": 17,
				"minute": 30,
				"nano": 1234,
				"second": 13
			  },
			  "weekdays": [
				"0"
			  ]
			}
		  ]`))
	})
	dps, _, err := client.Forecast.GetStaticRules(context.Background(), options)
	assert.NilError(t, err)
	assert.DeepEqual(t, want, dps)
}

func TestCreateDynamicRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := DynamicRule{

		Amount: 123,
		StartTime: LocalTime{
			Hour:   17,
			Minute: 30,
			Nano:   1234,
			Second: 13,
		},
		EndTime: LocalTime{
			Hour:   18,
			Minute: 31,
			Nano:   12345,
			Second: 14,
		},
		ExternalID:                 "e",
		ExternalForecastVariableID: "efvid",
		ShiftTypes: []ShiftType{
			{
				Amount:      321,
				ShiftTypeID: "abc",
			},
		},
		Weekdays: []Weekday{Monday},
	}
	mux.HandleFunc("/forecasts/dynamic-rules", func(w http.ResponseWriter, r *http.Request) {
		// Make sure it's POST
		testMethod(t, r, "POST")

		// Options are added as Query
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))

		// Body contains the proper JSON
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, fmt.Sprintf("%s\n", `{"amount":123,"endTime":{"hour":18,"minute":31,"nano":12345,"second":14},"startTime":{"hour":17,"minute":30,"nano":1234,"second":13},"externalId":"e","forecastExternalVariableId":"efvid","shiftTypes":[{"amount":321,"externalShiftTypeId":"abc"}],"weekdays":["0"]}`), string(body))
		assert.NilError(t, err)

		// Return the proper JSON
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	// Create a rule
	dps, _, err := client.Forecast.CreateDynamicRule(context.Background(), &want, options)
	assert.NilError(t, err)
	assert.DeepEqual(t, &want, dps)

	// Invalid options
	_, _, err = client.Forecast.CreateDynamicRule(context.Background(), &want, &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	options.ExternalSectionID = String("foo")
	client.Forecast.CreateDynamicRule(context.Background(), &want, options)
}

func TestCreateStaticRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := StaticRule{

		Comment:   "hello there",
		StartDate: time.Date(2019, time.July, 12, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2019, time.July, 13, 0, 0, 0, 0, time.UTC),
		StartTime: LocalTime{
			Hour:   17,
			Minute: 30,
			Nano:   1234,
			Second: 13,
		},
		EndTime: LocalTime{
			Hour:   18,
			Minute: 31,
			Nano:   12345,
			Second: 14,
		},
		ExternalID: "e",
		ShiftType: ShiftType{

			Amount:      321,
			ShiftTypeID: "abc",
		},
		Weekdays: []Weekday{Monday},
	}
	mux.HandleFunc("/forecasts/static-rules", func(w http.ResponseWriter, r *http.Request) {
		// Make sure it's POST
		testMethod(t, r, "POST")

		// Options are added as Query
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))

		// Body contains the proper JSON
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, fmt.Sprintf("%s\n", `{"comment":"hello there","startDate":"2019-07-12T00:00:00Z","endDate":"2019-07-13T00:00:00Z","startTime":{"hour":17,"minute":30,"nano":1234,"second":13},"endTime":{"hour":18,"minute":31,"nano":12345,"second":14},"externalId":"e","repeatPeriod":0,"shiftType":{"amount":321,"externalShiftTypeId":"abc"},"weekdays":["0"]}`), string(body))
		assert.NilError(t, err)

		// Return the proper JSON
		json, err := json.Marshal(&want)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	// Create a rule
	dps, _, err := client.Forecast.CreateStaticRule(context.Background(), &want, options)
	assert.NilError(t, err)
	assert.DeepEqual(t, &want, dps)

	// Invalid options
	_, _, err = client.Forecast.CreateStaticRule(context.Background(), &want, &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	options.ExternalSectionID = String("foo")
	client.Forecast.CreateStaticRule(context.Background(), &want, options)
}

func TestUpdateDynamicRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := DynamicRule{

		Amount: 123,
		StartTime: LocalTime{
			Hour:   17,
			Minute: 30,
			Nano:   1234,
			Second: 13,
		},
		EndTime: LocalTime{
			Hour:   18,
			Minute: 31,
			Nano:   12345,
			Second: 14,
		},
		ExternalID:                 "e",
		ExternalForecastVariableID: "efvid",
		ShiftTypes: []ShiftType{
			{
				Amount:      321,
				ShiftTypeID: "abc",
			},
		},
		Weekdays: []Weekday{Monday},
	}
	mux.HandleFunc("/forecasts/dynamic-rules", func(w http.ResponseWriter, r *http.Request) {
		// Make sure it's PUT
		testMethod(t, r, "PUT")

		// Options are added as Query
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))

		// Body contains the proper JSON
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, fmt.Sprintf("%s\n", `{"amount":123,"endTime":{"hour":18,"minute":31,"nano":12345,"second":14},"startTime":{"hour":17,"minute":30,"nano":1234,"second":13},"externalId":"e","forecastExternalVariableId":"efvid","shiftTypes":[{"amount":321,"externalShiftTypeId":"abc"}],"weekdays":["0"]}`), string(body))
		assert.NilError(t, err)
		fmt.Fprint(w, ``)
	})

	// Update a rule
	_, err := client.Forecast.UpdateDynamicRule(context.Background(), &want, options)
	assert.NilError(t, err)

	// Invalid options
	_, err = client.Forecast.UpdateDynamicRule(context.Background(), &want, &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	options.ExternalSectionID = String("foo")
	client.Forecast.UpdateDynamicRule(context.Background(), &want, options)
}

func TestUpdateStaticRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}

	want := StaticRule{

		Comment:   "hello there",
		StartDate: time.Date(2019, time.July, 12, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2019, time.July, 13, 0, 0, 0, 0, time.UTC),
		StartTime: LocalTime{
			Hour:   17,
			Minute: 30,
			Nano:   1234,
			Second: 13,
		},
		EndTime: LocalTime{
			Hour:   18,
			Minute: 31,
			Nano:   12345,
			Second: 14,
		},
		ExternalID: "e",
		ShiftType: ShiftType{

			Amount:      321,
			ShiftTypeID: "abc",
		},
		Weekdays: []Weekday{Monday},
	}
	mux.HandleFunc("/forecasts/static-rules", func(w http.ResponseWriter, r *http.Request) {
		// Make sure it's PUT
		testMethod(t, r, "PUT")

		// Options are added as Query
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))

		// Body contains the proper JSON
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, fmt.Sprintf("%s\n", `{"comment":"hello there","startDate":"2019-07-12T00:00:00Z","endDate":"2019-07-13T00:00:00Z","startTime":{"hour":17,"minute":30,"nano":1234,"second":13},"endTime":{"hour":18,"minute":31,"nano":12345,"second":14},"externalId":"e","repeatPeriod":0,"shiftType":{"amount":321,"externalShiftTypeId":"abc"},"weekdays":["0"]}`), string(body))
		assert.NilError(t, err)
		fmt.Fprint(w, ``)
	})

	// Update a rule
	_, err := client.Forecast.UpdateStaticRule(context.Background(), &want, options)
	assert.NilError(t, err)

	// Invalid options
	_, err = client.Forecast.UpdateStaticRule(context.Background(), &want, &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	options.ExternalSectionID = String("foo")
	client.Forecast.UpdateStaticRule(context.Background(), &want, options)
}

func TestDeleteDynamicRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}
	wantID := "mycoolID"

	mux.HandleFunc(fmt.Sprintf("/forecasts/dynamic-rules/%s", wantID), func(w http.ResponseWriter, r *http.Request) {
		// Make sure it's DELETE
		testMethod(t, r, "DELETE")

		// Options are added as Query
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))

		fmt.Fprint(w, ``)
	})

	// Update a rule
	_, err := client.Forecast.DeleteDynamicRule(context.Background(), wantID, options)
	assert.NilError(t, err)

	// Invalid options
	_, err = client.Forecast.DeleteDynamicRule(context.Background(), wantID, &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	options.ExternalSectionID = String("foo")
	client.Forecast.DeleteDynamicRule(context.Background(), wantID, options)
}

func TestDeleteStaticRule(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	options := &RequestOptions{
		ExternalUnitID: String("d"),
	}
	wantID := "mycoolID"

	mux.HandleFunc(fmt.Sprintf("/forecasts/static-rules/%s", wantID), func(w http.ResponseWriter, r *http.Request) {
		// Make sure it's DELETE
		testMethod(t, r, "DELETE")

		// Options are added as Query
		assert.Equal(t, "d", r.URL.Query().Get("externalUnitId"))

		fmt.Fprint(w, ``)
	})

	// Update a rule
	_, err := client.Forecast.DeleteStaticRule(context.Background(), wantID, options)
	assert.NilError(t, err)

	// Invalid options
	_, err = client.Forecast.DeleteStaticRule(context.Background(), wantID, &RequestOptions{})
	assert.ErrorContains(t, err, "Required fields in the Options")

	options.ExternalSectionID = String("foo")
	client.Forecast.DeleteStaticRule(context.Background(), wantID, options)
}
