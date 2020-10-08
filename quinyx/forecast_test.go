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
		DataProviderInputs: []DataProvider{
			{
				ExternalForecastVariableID: "a",
				ExternalSectionID:          "b",
				ExternalUnitID:             "c",
				DataPayload: []Payload{
					{
						Data:      123,
						Timestamp: Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
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
		DataProviderInputs: []DataProvider{
			{
				ExternalForecastVariableID: "a",
				ExternalSectionID:          "b",
				ExternalUnitID:             "c",
				DataPayload: []Payload{
					{
						Data:      123,
						Timestamp: Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
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
			ExternalForecastVariableID: "b",
			ExternalSectionID:          "c",
			ExternalUnitID:             "d",
			DataPayload: []Payload{
				{
					Data:      123,
					Timestamp: Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
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
