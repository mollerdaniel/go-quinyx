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

func TestGetAllCategories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tags/categories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
			  "color": "blue",
			  "externalId": "eid",
			  "id": 123,
			  "name": "n",
			  "tagType": "COST_CENTER"
			}
		  ]`)
	})
	categoriesResponse, _, err := client.Tags.GetAllCategories(context.Background())
	assert.NilError(t, err)

	want := []*TagCategory{
		{
			Color:      String("blue"),
			ExternalID: String("eid"),
			TagID:      Int32(123),
			Name:       String("n"),
			TagType:    CostCenter,
		},
	}
	assert.DeepEqual(t, categoriesResponse, want)
}

func TestGetCategory(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tags/categories/example", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
			{
			  "color": "blue",
			  "externalId": "example",
			  "id": 123,
			  "name": "string",
			  "tagType": "COST_CENTER"
			}
		  `)
	})
	categoriesResponse, _, err := client.Tags.GetCategory(context.Background(), "example")
	assert.NilError(t, err)

	want := &TagCategory{
		Color:      String("blue"),
		ExternalID: String("example"),
		TagID:      Int32(123),
		Name:       String("string"),
		TagType:    CostCenter,
	}
	assert.DeepEqual(t, categoriesResponse, want)
}

func TestGetAllTags(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tags/categories/example/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		{
			"categoryExternalId": "example",
			"code": "c",
			"coordinates": [
			  {
				"latitude": 12.1234,
				"longitude": 24.1234,
				"radius": 15
			  }
			],
			"customFields": [
			  {
				"label": "l",
				"value": "v"
			  }
			],
			"endDate": "2019-10-12T07:20:50.52Z",
			"externalId": "eid",
			"information": "inf",
			"name": "n",
			"periods": [
			  {
				"from": "2019-10-12T07:20:50.52Z",
				"to": "2019-10-12T12:20:50.52Z",
				"hours": 5,
				"type": "PERIOD",
				"count": 9
			  }
			],
			"startDate": "2020-10-12T07:20:50.52Z",
			"uniqueScheduling": true,
			"unitExternalId": "uid"
		  }
		  `)
	})
	tag, _, err := client.Tags.GetAllTags(context.Background(), "example")
	assert.NilError(t, err)

	want := &Tag{
		CategoryExternalID: String("example"),
		Code:               String("c"),
		Coordinates: []*Coordinate{
			{
				Latitude:  Float64(12.1234),
				Longitude: Float64(24.1234),
				Radius:    Int32(15),
			},
		},
		CustomFields: []*CustomField{
			{
				Label: String("l"),
				Value: String("v"),
			},
		},
		EndDate:     &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		ExternalID:  String("eid"),
		Information: String("inf"),
		Name:        String("n"),
		Periods: []*Period{{
			From:  &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
			To:    &Timestamp{time.Date(2019, time.October, 12, 12, 20, 50, 520000000, time.UTC)},
			Hours: Float64(5),
			Type:  PeriodTypePeriod,
			Count: Float64(9),
		}},
		StartDate:        &Timestamp{time.Date(2020, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		UniqueScheduling: Bool(true),
		UnitExternalID:   String("uid"),
	}
	assert.DeepEqual(t, tag, want)
}

func TestGetTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tags/categories/example/tags/eid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		{
			"categoryExternalId": "example",
			"code": "c",
			"coordinates": [
			  {
				"latitude": 12.1234,
				"longitude": 24.1234,
				"radius": 15
			  }
			],
			"customFields": [
			  {
				"label": "l",
				"value": "v"
			  }
			],
			"endDate": "2019-10-12T07:20:50.52Z",
			"externalId": "eid",
			"information": "inf",
			"name": "n",
			"periods": [
			  {
				"from": "2019-10-12T07:20:50.52Z",
				"to": "2019-10-12T12:20:50.52Z",
				"hours": 5,
				"type": "PERIOD",
				"count": 9
			  }
			],
			"startDate": "2020-10-12T07:20:50.52Z",
			"uniqueScheduling": true,
			"unitExternalId": "uid"
		  }
		  `)
	})
	tag, _, err := client.Tags.GetTag(context.Background(), "example", "eid")
	assert.NilError(t, err)

	want := &Tag{
		CategoryExternalID: String("example"),
		Code:               String("c"),
		Coordinates: []*Coordinate{
			{
				Latitude:  Float64(12.1234),
				Longitude: Float64(24.1234),
				Radius:    Int32(15),
			},
		},
		CustomFields: []*CustomField{
			{
				Label: String("l"),
				Value: String("v"),
			},
		},
		EndDate:     &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		ExternalID:  String("eid"),
		Information: String("inf"),
		Name:        String("n"),
		Periods: []*Period{{
			From:  &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
			To:    &Timestamp{time.Date(2019, time.October, 12, 12, 20, 50, 520000000, time.UTC)},
			Hours: Float64(5),
			Type:  PeriodTypePeriod,
			Count: Float64(9),
		}},
		StartDate:        &Timestamp{time.Date(2020, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		UniqueScheduling: Bool(true),
		UnitExternalID:   String("uid"),
	}
	assert.DeepEqual(t, tag, want)
}

func TestCreateTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	want := &Tag{
		CategoryExternalID: String("example"),
		Code:               String("c"),
		Coordinates: []*Coordinate{
			{
				Latitude:  Float64(12.1234),
				Longitude: Float64(24.1234),
				Radius:    Int32(15),
			},
		},
		CustomFields: []*CustomField{
			{
				Label: String("l"),
				Value: String("v"),
			},
		},
		EndDate:     &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		ExternalID:  String("eid"),
		Information: String("inf"),
		Name:        String("n"),
		Periods: []*Period{{
			From:  &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
			To:    &Timestamp{time.Date(2019, time.October, 12, 12, 20, 50, 520000000, time.UTC)},
			Hours: Float64(5),
			Type:  PeriodTypePeriod,
			Count: Float64(9),
		}},
		StartDate:        &Timestamp{time.Date(2020, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		UniqueScheduling: Bool(true),
		UnitExternalID:   String("uid"),
	}
	mux.HandleFunc("/tags/categories/example/tags", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		tagbody := &Tag{}
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, `{"categoryExternalId":"example","code":"c","coordinates":[{"latitude":12.1234,"longitude":24.1234,"radius":15}],"customFields":[{"label":"l","value":"v"}],"endDate":"2019-10-12T07:20:50.52Z","externalId":"eid","information":"inf","name":"n","periods":[{"from":"2019-10-12T07:20:50.52Z","to":"2019-10-12T12:20:50.52Z","hours":5,"type":"PERIOD","count":9}],"startDate":"2020-10-12T07:20:50.52Z","uniqueScheduling":true,"unitExternalId":"uid"}
`, string(body))
		assert.NilError(t, err)
		err = json.Unmarshal(body, tagbody)
		assert.NilError(t, err)
		json, err := json.Marshal(tagbody)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	tag, _, err := client.Tags.CreateTag(context.Background(), "example", want)
	assert.NilError(t, err)
	assert.DeepEqual(t, tag, want)
}

func TestUpdateTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	want := &Tag{
		CategoryExternalID: String("example"),
		Code:               String("c"),
		Coordinates: []*Coordinate{
			{
				Latitude:  Float64(12.1234),
				Longitude: Float64(24.1234),
				Radius:    Int32(15),
			},
		},
		CustomFields: []*CustomField{
			{
				Label: String("l"),
				Value: String("v"),
			},
		},
		EndDate:     &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		ExternalID:  String("eid"),
		Information: String("inf"),
		Name:        String("n"),
		Periods: []*Period{{
			From:  &Timestamp{time.Date(2019, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
			To:    &Timestamp{time.Date(2019, time.October, 12, 12, 20, 50, 520000000, time.UTC)},
			Hours: Float64(5),
			Type:  PeriodTypePeriod,
			Count: Float64(9),
		}},
		StartDate:        &Timestamp{time.Date(2020, time.October, 12, 07, 20, 50, 520000000, time.UTC)},
		UniqueScheduling: Bool(true),
		UnitExternalID:   String("uid"),
	}
	mux.HandleFunc("/tags/categories/example/tags/eid", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "PUT")
		tagbody := &Tag{}
		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, `{"categoryExternalId":"example","code":"c","coordinates":[{"latitude":12.1234,"longitude":24.1234,"radius":15}],"customFields":[{"label":"l","value":"v"}],"endDate":"2019-10-12T07:20:50.52Z","externalId":"eid","information":"inf","name":"n","periods":[{"from":"2019-10-12T07:20:50.52Z","to":"2019-10-12T12:20:50.52Z","hours":5,"type":"PERIOD","count":9}],"startDate":"2020-10-12T07:20:50.52Z","uniqueScheduling":true,"unitExternalId":"uid"}
`, string(body))
		assert.NilError(t, err)
		err = json.Unmarshal(body, tagbody)
		assert.NilError(t, err)
		json, err := json.Marshal(tagbody)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})
	mux.HandleFunc("/tags/categories/nochange/tags/eid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		tagbody := &Tag{}
		body, err := ioutil.ReadAll(r.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(body, tagbody)
		assert.NilError(t, err)
		json, err := json.Marshal(tagbody)
		assert.NilError(t, err)
		fmt.Fprint(w, string(json))
	})

	tag, _, err := client.Tags.UpdateTag(context.Background(), "example", "eid", want)
	assert.NilError(t, err)
	assert.DeepEqual(t, tag, want)
	want.CategoryExternalID = String("blah")
	_, _, err = client.Tags.UpdateTag(context.Background(), "nochange", "eid", want)
	assert.ErrorContains(t, err, "categoryExternalID cannot be changed")
}

func TestDeleteTag(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tags/categories/example/tags/eid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	_, err := client.Tags.DeleteTag(context.Background(), "example", "eid")
	assert.NilError(t, err)
}
