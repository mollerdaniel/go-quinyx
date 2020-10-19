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
	options := &RequestOptions{
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
	_, _, err = client.Forecast.GetActualData(context.Background(), "a", &RequestOptions{})
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
	options := &RequestOptions{
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
	_, _, err = client.Forecast.GetActualDataStream(context.Background(), "a", &RequestOptions{})
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
	options := &RequestOptions{
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
	_, _, err = client.Forecast.GetAggregatedData(context.Background(), "a", &RequestOptions{})
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
	options := &RequestOptions{
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
	_, _, err = client.Forecast.GetCalculatedForecast(context.Background(), "a", &RequestOptions{})
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
	options := &EditCalculatedOptions{
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
	options := &RequestOptions{
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
	_, _, err = client.Forecast.GetForecastData(context.Background(), "a", &RequestOptions{})
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
	options := &RequestOptions{
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
	_, err = client.Forecast.DeleteForecastData(context.Background(), "a", &RequestOptions{})
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
	options := &RequestOptions{
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
	_, err = client.Forecast.DeleteActualData(context.Background(), "a", &RequestOptions{})
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
	options := &RequestOptions{
		StartTime:      time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.October, 17, 12, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}
	assert.Equal(t, float64(5), options.dayDistance())
	options = &RequestOptions{
		StartTime:      time.Date(2019, time.July, 12, 07, 20, 50, 520000000, time.UTC),
		EndTime:        time.Date(2019, time.December, 17, 12, 20, 50, 520000000, time.UTC),
		ExternalUnitID: String("d"),
	}
	assert.Assert(t, options.dayDistance() > 120)
}
