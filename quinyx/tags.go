package quinyx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// TagsService handles Quinyx Tags
//
// Quinyx API docs: https://api.quinyx.com/v2/docs/swagger-ui.html?urls.primaryName=tags#/
type TagsService service

// Tag defines a Quinyx TagIntegration object
type Tag struct {
	CategoryExternalID *string        `json:"categoryExternalId,omitempty"`
	Code               *string        `json:"code,omitempty"`
	Coordinates        []*Coordinate  `json:"coordinates,omitempty"`
	CustomFields       []*CustomField `json:"customFields,omitempty"`
	EndDate            *Timestamp     `json:"endDate,omitempty"`
	ExternalID         *string        `json:"externalId,omitempty"`
	Information        *string        `json:"information,omitempty"`
	Name               *string        `json:"name,omitempty"`
	Periods            []*Period      `json:"periods,omitempty"`
	StartDate          *Timestamp     `json:"startDate,omitempty"`
	UniqueScheduling   *bool          `json:"uniqueScheduling,omitempty"`
	UnitExternalID     *string        `json:"unitExternalId,omitempty"`
}

// Coordinate defines a Geofence using Long Lat and a Radius
type Coordinate struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Radius    *int32   `json:"radius,omitempty"`
}

// CustomField is a Tag custom field
type CustomField struct {
	Label *string `json:"label,omitempty"`
	Value *string `json:"value,omitempty"`
}

// Period defines a point in time for the tag
type Period struct {
	From  *Timestamp `json:"from,omitempty"`
	To    *Timestamp `json:"to,omitempty"`
	Hours *float64   `json:"hours,omitempty"`
	Type  PeriodType `json:"type,omitempty"`
	Count *float64   `json:"count,omitempty"`
}

// TagCategory is a Category of tags
type TagCategory struct {
	Color      *string `json:"color,omitempty"`
	ExternalID *string `json:"externalId,omitempty"`
	TagID      *int32  `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	TagType    TagType `json:"tagType,omitempty"`
}

// TagType defines a type of Tag
type TagType string

// TagTypes
const (
	CostCenter TagType = "COST_CENTER"
	Project            = "PROJECT"
	Account            = "ACCOUNT"
	Extended           = "EXTENDED"
)

// PeriodType is the type of the Period
type PeriodType string

// PeriodType
const (
	PeriodTypePeriod PeriodType = "PERIOD"
	PeriodTypeDays              = "DAYS"
	PeriodTypeWeeks             = "WEEKS"
)

// UnmarshalJSON TagType enum
func (tt *TagType) UnmarshalJSON(b []byte) error {
	type TT TagType
	var r *TT = (*TT)(tt)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *tt {
	case CostCenter, Project, Account, Extended:
		return nil
	}
	return errors.New("Invalid leave type")
}

// GetAllCategories gets all categories
func (s *TagsService) GetAllCategories(ctx context.Context) ([]*TagCategory, *Response, error) {
	u := "tags/categories"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var categories []*TagCategory
	resp, err := s.client.Do(ctx, req, &categories)
	if err != nil {
		return nil, resp, err
	}
	return categories, resp, nil
}

// GetCategory from a categoryExternalID
func (s *TagsService) GetCategory(ctx context.Context, categoryExternalID string) (*TagCategory, *Response, error) {
	u := fmt.Sprintf("tags/categories/%v", categoryExternalID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var category *TagCategory
	resp, err := s.client.Do(ctx, req, &category)
	if err != nil {
		return nil, resp, err
	}
	return category, resp, nil
}

// GetAllTags based on categoryExternalID
// While the documentation says it should return all, it actually only returns one tag.
func (s *TagsService) GetAllTags(ctx context.Context, categoryExternalID string) (*Tag, *Response, error) {
	u := fmt.Sprintf("tags/categories/%v/tags", categoryExternalID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var tag *Tag
	resp, err := s.client.Do(ctx, req, &tag)
	if err != nil {
		return nil, resp, err
	}
	return tag, resp, nil
}

// GetTag returns the specified tag by external tag category id and external tag id
func (s *TagsService) GetTag(ctx context.Context, categoryExternalID string, tagExternalID string) (*Tag, *Response, error) {
	u := fmt.Sprintf("tags/categories/%v/tags/%v", categoryExternalID, tagExternalID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var tag *Tag
	resp, err := s.client.Do(ctx, req, &tag)
	if err != nil {
		return nil, resp, err
	}
	return tag, resp, nil
}

// CreateTag creates and then returns the tag
func (s *TagsService) CreateTag(ctx context.Context, categoryExternalID string, tag *Tag) (*Tag, *Response, error) {
	u := fmt.Sprintf("tags/categories/%v/tags", categoryExternalID)
	req, err := s.client.NewRequest("POST", u, tag)
	if err != nil {
		return nil, nil, err
	}
	var tagres *Tag
	resp, err := s.client.Do(ctx, req, &tagres)
	if err != nil {
		return nil, resp, err
	}
	return tagres, resp, nil
}

// UpdateTag using a tagdelta object where set values will be changed
func (s *TagsService) UpdateTag(ctx context.Context, categoryExternalID string, tagExternalID string, tag *Tag) (*Tag, *Response, error) {
	u := fmt.Sprintf("tags/categories/%v/tags/%v", categoryExternalID, tagExternalID)

	// See documentation for limitations on changing categoryExternalID
	// https://api.quinyx.com/v2/docs/swagger-ui.html?urls.primaryName=tags#/tag-integration-api-controller/updateTagByExternalIdUsingPUT
	if *tag.CategoryExternalID != categoryExternalID {
		return nil, nil, fmt.Errorf("categoryExternalID cannot be changed")
	}

	req, err := s.client.NewRequest("PUT", u, tag)

	if err != nil {
		return nil, nil, err
	}
	var tagres *Tag
	resp, err := s.client.Do(ctx, req, &tagres)
	if err != nil {
		return nil, resp, err
	}
	return tagres, resp, nil
}

// DeleteTag removes the tag
func (s *TagsService) DeleteTag(ctx context.Context, categoryExternalID string, tagExternalID string) (*Response, error) {
	u := fmt.Sprintf("tags/categories/%v/tags/%v", categoryExternalID, tagExternalID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
